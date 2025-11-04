import os
import pyodbc
import redis
import pika
import json
import logging
from flask import Flask, redirect, jsonify, abort

# --- Configuration ---
app = Flask(__name__)
logging.basicConfig(level=logging.INFO)

# Get all connection info from Environment Variables
DB_HOST = os.environ.get('DB_HOST')
DB_NAME = "UrlShortenerDb"
DB_USER = os.environ.get('DB_USER')
DB_PASSWORD = os.environ.get('DB_PASSWORD')

REDIS_HOST = os.environ.get('REDIS_HOST')

RABBITMQ_HOST = os.environ.get('RABBITMQ_HOST')
RABBITMQ_USER = os.environ.get('RABBITMQ_USER')
RABBITMQ_PASS = os.environ.get('RABBITMQ_PASS')

# --- Redis Connection ---
try:
    # decode_responses=True makes it return strings, not bytes
    redis_client = redis.Redis(host=REDIS_HOST, port=6379, db=0, decode_responses=True)
    redis_client.ping()
    logging.info("Connected to Redis successfully!")
except Exception as e:
    logging.error(f"Could not connect to Redis: {e}")
    redis_client = None

# --- RabbitMQ Connection (Fire-and-Forget) ---
def send_analytics_event(short_code):
    """
    Sends a 'click' event to the analytics queue.
    This is "fire-and-forget". If it fails, we log the error
    but DO NOT stop the user's redirect.
    """
    try:
        creds = pika.PlainCredentials(RABBITMQ_USER, RABBITMQ_PASS)
        params = pika.ConnectionParameters(host=RABBITMQ_HOST, credentials=creds)
        
        with pika.BlockingConnection(params) as connection:
            channel = connection.channel()
            # Ensure the queue exists and is durable (survives broker restart)
            channel.queue_declare(queue='analytics-queue', durable=True)
            
            message = json.dumps({'shortCode': short_code})
            
            channel.basic_publish(
                exchange='',
                routing_key='analytics-queue',
                body=message,
                properties=pika.BasicProperties(
                    delivery_mode = 2, # Make message persistent
                ))
        logging.info(f"Sent analytics event for {short_code}")
    except Exception as e:
        # If this fails, the user redirect is more important.
        logging.error(f"Could not send analytics event: {e}")

# --- SQL Database Connection ---
db_conn = None

def get_db_connection():
    """Establishes or re-establishes the database connection."""
    global db_conn
    try:
        if db_conn is None or db_conn.closed:
            conn_string = (
                f"DRIVER={{ODBC Driver 18 for SQL Server}};"
                f"SERVER={DB_HOST},1433;"
                f"DATABASE={DB_NAME};"
                f"UID={DB_USER};"
                f"PWD={DB_PASSWORD};"
                f"Encrypt=no;"
            )
            db_conn = pyodbc.connect(conn_string, timeout=5)
        return db_conn
    except Exception as e:
        logging.error(f"Error connecting to database: {e}")
        return None

# --- API Routes ---

@app.route('/<string:short_code>')
def handle_redirect(short_code):
    """
    This is the main redirect logic.
    1. Check Redis (Cache)
    2. If miss, check SQL (Database)
    3. If found, cache it and redirect.
    4. If not found, 404.
    """
    
    # 1. --- Check Cache (Redis) ---
    if redis_client:
        try:
            long_url = redis_client.get(short_code)
            if long_url:
                logging.info(f"Cache HIT for: {short_code}")
                # Asynchronously send analytics event
                send_analytics_event(short_code)
                # 302 Found (Temporary Redirect)
                return redirect(long_url, code=302)
        except Exception as e:
            logging.error(f"Redis error: {e}")
    
    # 2. --- Cache MISS, Check Database (SQL) ---
    logging.info(f"Cache MISS for: {short_code}")
    
    conn = get_db_connection()
    if not conn:
        return jsonify({"error": "Database connection failed"}), 503

    try:
        cursor = conn.cursor()
        cursor.execute("SELECT LongUrl FROM Links WHERE ShortCode = ?", (short_code))
        row = cursor.fetchone()
        
        # 3. --- Found in DB ---
        if row:
            long_url = row.LongUrl
            
            # 3a. Update cache for next time (e.g., 1 hour expiry)
            if redis_client:
                try:
                    redis_client.set(short_code, long_url, ex=3600)
                except Exception as e:
                    logging.error(f"Redis set error: {e}")
            
            # 3b. Send analytics event
            send_analytics_event(short_code)
            
            # 3c. Redirect user
            return redirect(long_url, code=302)
        
        # 4. --- Not Found Anywhere ---
        else:
            return jsonify({"error": "Short URL not found"}), 404
            
    except Exception as e:
        logging.error(f"Database query error: {e}")
        return jsonify({"error": "An internal error occurred"}), 500

@app.route('/health', methods=['GET'])
def health_check():
    """Health check for Kubernetes probes."""
    try:
        # Check DB
        conn = get_db_connection()
        if not conn:
            return "DB Connection Failed", 503
        conn.cursor().execute("SELECT 1")
        
        # Check Redis
        if not redis_client or not redis_client.ping():
            return "Redis Connection Failed", 503
            
        return "OK", 200
    except Exception as e:
        return str(e), 503

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001) # Note: Using port 5001 to avoid conflicts

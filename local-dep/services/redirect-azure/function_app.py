import azure.functions as func
import logging
import os
import redis
import pyodbc
import pika
import json

# --- Connection Clients (Initialized globally for reuse) ---
# We will set these connection strings in 'local.settings.json'

# Initialize Redis client
# Note: These connections will FAIL until we start port-forwarding.
try:
    redis_client = redis.Redis(
        host=os.environ.get("REDIS_HOST", "localhost"), 
        port=os.environ.get("REDIS_PORT", 6379), 
        db=0, 
        decode_responses=True,
        socket_timeout=5
    )
    redis_client.ping()
    logging.info("Connected to Redis successfully!")
except Exception as e:
    logging.warning(f"Could not connect to Redis on startup: {e}")
    redis_client = None

# Database connection string (we'll get this from local.settings.json)
DB_CONN_STRING = os.environ.get("SQL_CONNECTION_STRING")

# RabbitMQ Connection Info
RABBITMQ_HOST = os.environ.get("RABBITMQ_HOST", "localhost")
RABBITMQ_USER = os.environ.get("RABBITMQ_USER", "user")
RABBITMQ_PASS = os.environ.get("RABBITMQ_PASS", "password")


# --- RabbitMQ Helper Function ---
def send_analytics_event(short_code):
    """
    Sends a 'click' event to the analytics queue.
    This is "fire-and-forget". If it fails, we log the error
    but DO NOT stop the user's redirect.
    """
    try:
        creds = pika.PlainCredentials(RABBITMQ_USER, RABBITMQ_PASS)
        params = pika.ConnectionParameters(host=RABBITMQ_HOST, credentials=creds, port=5672)
        
        with pika.BlockingConnection(params) as connection:
            channel = connection.channel()
            channel.queue_declare(queue='analytics-queue', durable=True)
            
            message = json.dumps({'shortCode': short_code})
            
            channel.basic_publish(
                exchange='',
                routing_key='analytics-queue',
                body=message,
                properties=pika.BasicProperties(delivery_mode=2) # Make persistent
            )
        logging.info(f"Sent analytics event for {short_code}")
    except Exception as e:
        # If this fails, the user redirect is more important.
        logging.error(f"Could not send analytics event: {e}")


# --- The Main Azure Function ---
app = func.FunctionApp()

@app.route(route="{shortCode}", auth_level=func.AuthLevel.ANONYMOUS)
def RedirectHandler(req: func.HttpRequest) -> func.HttpResponse:
    """
    This is the main redirect logic.
    1. Get 'shortCode' from the route.
    2. Check Redis (Cache)
    3. If miss, check SQL (Database)
    4. If found, cache it, send analytics, and redirect.
    5. If not found, 404.
    """
    
    # 1. Get 'shortCode' from the route
    short_code = req.route_params.get('shortCode')
    if not short_code:
        return func.HttpResponse("No short code provided.", status_code=400)

    logging.info(f"Processing redirect request for: {short_code}")

    # 2. --- Check Cache (Redis) ---
    long_url = None
    if redis_client:
        try:
            long_url = redis_client.get(short_code)
            if long_url:
                logging.info(f"Cache HIT for: {short_code}")
        except Exception as e:
            logging.warning(f"Redis 'get' error (will try DB): {e}")
            
    # 3. --- Cache MISS, Check Database (SQL) ---
    if not long_url:
        logging.info(f"Cache MISS for: {short_code}")
        if not DB_CONN_STRING:
            logging.error("SQL_CONNECTION_STRING is not set!")
            return func.HttpResponse("Server configuration error.", status_code=500)
            
        try:
            with pyodbc.connect(DB_CONN_STRING, timeout=5) as conn:
                cursor = conn.cursor()
                cursor.execute("SELECT LongUrl FROM Links WHERE ShortCode = ?", (short_code))
                row = cursor.fetchone()
                
                if row:
                    # 4. --- Found in DB ---
                    logging.info(f"Database HIT for: {short_code}")
                    long_url = row.LongUrl
                    
                    # 4a. Update cache for next time
                    if redis_client:
                        try:
                            redis_client.set(short_code, long_url, ex=3600) # 1 hour expiry
                        except Exception as e:
                            logging.warning(f"Redis 'set' error: {e}")
                else:
                    # 5. --- Not Found Anywhere ---
                    logging.warning(f"Code NOT FOUND: {short_code}")
                    return func.HttpResponse("Short URL not found.", status_code=404)

        except Exception as e:
            logging.error(f"Database query error: {e}")
            return func.HttpResponse("Server error during database query.", status_code=500)

    # --- If we are here, 'long_url' was found (from cache or DB) ---
    
    # 4b. Send analytics event
    send_analytics_event(short_code)
    
    # 4c. Redirect user
    return func.HttpResponse(
        status_code=302,  # 302 = Temporary Redirect
        headers={
            "Location": long_url
        }
    )

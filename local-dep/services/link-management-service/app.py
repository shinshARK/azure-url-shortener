import os
import pyodbc
import shortuuid
from flask import Flask, request, jsonify

# --- Configuration ---

app = Flask(__name__)

# Get database config from Environment Variables
# We will set these in our Kubernetes Deployment
DB_HOST = os.environ.get('DB_HOST')
DB_NAME = "UrlShortenerDb" # We can hardcode this or use an env var
DB_USER = os.environ.get('DB_USER')
DB_PASSWORD = os.environ.get('DB_PASSWORD')

# Build the connection string
connection_string = (
    f"DRIVER={{ODBC Driver 18 for SQL Server}};"
    f"SERVER={DB_HOST},1433;"
    f"DATABASE={DB_NAME};"
    f"UID={DB_USER};"
    f"PWD={DB_PASSWORD};"
    f"Encrypt=no;"  # 'no' is okay for local dev without SSL
)

db_conn = None

def get_db_connection():
    """Establishes or re-establishes the database connection."""
    global db_conn
    try:
        if db_conn is None or db_conn.closed:
            db_conn = pyodbc.connect(connection_string, timeout=5)
        return db_conn
    except Exception as e:
        print(f"Error connecting to database: {e}")
        return None

# --- Application Logic ---

@app.route('/create', methods=['POST'])
def create_link():
    """
    Receives a POST request with {'longUrl': '...'}
    and creates a new short link in the database.
    """
    data = request.get_json()
    if not data or 'longUrl' not in data:
        return jsonify({"error": "longUrl is required"}), 400

    long_url = data['longUrl']
    
    # Generate a simple short code
    # We'll use 7 characters: e.g., 'fT7d8Xq'
    short_code = shortuuid.uuid()[:7]

    try:
        conn = get_db_connection()
        if not conn:
            return jsonify({"error": "Database connection failed"}), 500
            
        cursor = conn.cursor()
        
        # Insert the new link into the database
        sql_insert = "INSERT INTO Links (ShortCode, LongUrl, ClickCount) VALUES (?, ?, 0)"
        cursor.execute(sql_insert, (short_code, long_url))
        conn.commit()

        # Assuming your service will be exposed at a domain,
        # but for now, we'll just return the code.
        return jsonify({
            "shortCode": short_code,
            "longUrl": long_url
        }), 201

    except pyodbc.IntegrityError:
        # This happens if the shortCode (Primary Key) already exists.
        # We could retry, but for now, we'll just error.
        return jsonify({"error": "Short code collision, please try again"}), 500
    except Exception as e:
        return jsonify({"error": str(e)}), 500

@app.route('/health', methods=['GET'])
def health_check():
    """A simple health check endpoint for Kubernetes."""
    try:
        conn = get_db_connection()
        if not conn:
            return "DB Connection Failed", 503
        cursor = conn.cursor()
        cursor.execute("SELECT 1")
        return "OK", 200
    except Exception as e:
        return str(e), 503

if __name__ == '__main__':
    # This is for local testing, but K8s will use a proper server (gunicorn)
    app.run(host='0.0.0.0', port=5000)

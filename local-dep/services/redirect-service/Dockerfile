# --- Stage 1: Build Stage ---
# We use a full-fat image to build our dependencies,
# as it has all the necessary C build tools for 'pyodbc'.
FROM python:3.11-slim as builder

WORKDIR /app

# Install 'pyodbc' build dependencies
RUN apt-get update && \
    apt-get install -y gcc g++ unixodbc-dev

# Install python packages as wheels
COPY requirements.txt .
RUN pip wheel --no-cache-dir --wheel-dir=/app/wheels -r requirements.txt


# --- Stage 2: Final Stage ---
# We use a slim image for the final container to keep it small.
FROM python:3.11-slim

WORKDIR /app

# Install the Microsoft ODBC Driver 18 for SQL Server
# This method uses Microsoft's official .deb package, which handles all keys.
RUN apt-get update && \
    # Install curl (and sudo, which dpkg needs)
    apt-get install -y curl sudo && \
    # Download the official Microsoft repo config package
    curl -sSL -O https://packages.microsoft.com/config/debian/12/packages-microsoft-prod.deb && \
    # Install the config package
    dpkg -i packages-microsoft-prod.deb && \
    # Clean up the .deb file
    rm packages-microsoft-prod.deb && \
    # Update apt lists again (now with the new MS repo)
    apt-get update && \
    # Install the driver (ACCEPT_EULA=Y is critical)
    ACCEPT_EULA=Y apt-get install -y msodbcsql18 unixodbc && \
    # Final cleanup
    rm -rf /var/lib/apt/lists/*

# Copy the pre-built wheels from the 'builder' stage
COPY --from=builder /app/wheels /app/wheels

# Copy our application code
COPY app.py .

# Install the packages from the local wheels (this is fast and offline)
RUN pip install --no-cache --no-index --find-links=/app/wheels /app/wheels/*

# Expose the port our app runs on
EXPOSE 5001

# Command to run the application
CMD ["gunicorn", "--bind", "0.0.0.0:5001", "app:app"]

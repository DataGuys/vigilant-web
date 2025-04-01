# Use an official Python runtime as a parent image
FROM python:3.8-slim

# Install tor and any additional OS packages needed
RUN apt-get update && apt-get install -y tor && rm -rf /var/lib/apt/lists/*

# Set the working directory in the container
WORKDIR /app

# Copy requirements file and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application files (app.py, observer.py, config folder, etc.)
COPY . .

# Expose port 8080 for the Flask web app
EXPOSE 8080

# Start the Tor service in the background, wait for it to initialize, then run the Flask app
CMD ["sh", "-c", "tor & sleep 10 && python app.py"]

#!/bin/bash

echo "Setting up XReplyAgent with persistent MongoDB storage..."

# Create data directories
mkdir -p ./data/mongodb
mkdir -p ./data/mongodb_logs

# Set proper permissions for MongoDB
sudo chown -R 999:999 ./data/mongodb
sudo chown -R 999:999 ./data/mongodb_logs

echo "Directories created:"
echo "- MongoDB data: ./data/mongodb"
echo "- MongoDB logs: ./data/mongodb_logs"

# Start services
echo "Starting services with docker-compose..."
docker-compose up --build

echo "Services started!"
echo "API: http://localhost:3002"
echo "MongoDB: localhost:27017"
echo "Data persisted in: ./data/mongodb"
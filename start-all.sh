#!/bin/bash

echo "========================================"
echo "Starting all microservices..."
echo "========================================"

# Create network if not exists
docker network create gomall 2>/dev/null

# Remove old containers if they exist
docker rm -f cart checkout email frontend order payment product user 2>/dev/null

# Start all microservices
echo "Starting cart service..."
docker run -d --name cart --network gomall -p 8883:8883 gomall-cart:latest

echo "Starting checkout service..."
docker run -d --name checkout --network gomall -p 8884:8884 gomall-checkout:latest

echo "Starting email service..."
docker run -d --name email --network gomall -p 8885:8885 gomall-email:latest

echo "Starting frontend service..."
docker run -d --name frontend --network gomall -p 8080:8080 gomall-frontend:latest

echo "Starting order service..."
docker run -d --name order --network gomall -p 8886:8886 gomall-order:latest

echo "Starting payment service..."
docker run -d --name payment --network gomall -p 8887:8887 gomall-payment:latest

echo "Starting product service..."
docker run -d --name product --network gomall -p 8888:8888 gomall-product:latest

echo "Starting user service..."
docker run -d --name user --network gomall -p 8889:8889 gomall-user:latest

echo ""
echo "========================================"
echo "All services started!"
echo "========================================"
docker ps | grep gomall
echo ""
echo "Frontend is available at: http://localhost:8080"

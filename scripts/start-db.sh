#!/bin/bash

echo "Starting PostgreSQL with Docker Compose..."
docker-compose up -d

echo "Waiting for PostgreSQL to be ready..."
sleep 3

echo "Checking PostgreSQL status..."
docker ps | grep grpc-postgres

echo ""
echo "PostgreSQL is ready!"
echo "Connection info:"
echo "  Host: localhost"
echo "  Port: 5432"
echo "  Database: grpc_dev"
echo "  User: postgres"
echo "  Password: postgres"

#!/bin/bash

# Function to cleanup background processes on exit
cleanup() {
    echo "🛑 Stopping services..."
    kill $(jobs -p) 2>/dev/null
    exit
}

trap cleanup EXIT INT TERM

echo "🚀 Starting development servers..."
echo ""

# Start PostgreSQL database
echo "🗄️  Starting PostgreSQL database..."
./scripts/start-db.sh

# Wait for database to be ready
sleep 2
echo ""

# Start Go server
echo "📦 Starting Go server on :8080..."
cd server && go run main.go &

# Wait a bit for server to start
sleep 2

# Start Vite dev server
echo "⚡ Starting Vite dev server on :5173..."
cd client && npm run dev &

# Wait for all background processes
wait

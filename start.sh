#!/bin/bash

# Quick start script for go-asgi

set -e

echo "🚀 Starting go-asgi development environment..."

# Check if Python virtual environment exists
if [ ! -d "app/venv" ]; then
    echo "📦 Creating Python virtual environment..."
    cd app
    python3 -m venv venv
    source venv/bin/activate
    pip install -r requirements.txt
    cd ..
else
    echo "✅ Python virtual environment found"
fi

# Start FastAPI backend in background
echo "🐍 Starting FastAPI backend on port 8001..."
cd app
source venv/bin/activate
uvicorn main:app --port 8001 --reload &
FASTAPI_PID=$!
cd ..

# Wait for FastAPI to start
echo "⏳ Waiting for FastAPI to start..."
sleep 2

# Start Go server
echo "🔷 Starting Go server on port 8000..."
cd server
go run main.go &
GO_PID=$!
cd ..

echo ""
echo "✨ go-asgi is running!"
echo "   Go Server:    http://localhost:8000"
echo "   FastAPI Docs: http://localhost:8001/docs"
echo ""
echo "Press Ctrl+C to stop all servers"

# Trap Ctrl+C and kill both processes
trap "echo ''; echo '🛑 Stopping servers...'; kill $FASTAPI_PID $GO_PID 2>/dev/null; exit" INT

# Wait for processes
wait

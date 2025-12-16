# go-asgi

Learning project: Building an ASGI-compatible server in Go to understand web server internals and the async Python ecosystem.

## Project Goal

Understand how production web servers work by implementing one from scratch. This project bridges Go's high-performance networking with Python's async web frameworks (FastAPI), exploring the ASGI protocol that connects them.

## Learning Objectives

**Core Concepts:**
- How ASGI servers (Uvicorn, Hypercorn) communicate with Python web apps
- Network programming in Go (HTTP handling, connection pooling, concurrency)
- Cross-language integration patterns (Go ↔ Python)
- Production architecture (reverse proxies, load balancing)

**Technical Skills:**
- Implementing HTTP servers with Go's `net/http` package
- Understanding the ASGI specification
- Benchmarking and performance profiling
- Systems programming and protocol implementation

## Architecture Plan

### Phase 1: Reverse Proxy (Current)
Build a Go HTTP server that forwards requests to a Python ASGI server (Uvicorn).

```
Client → Go Server (port 8000) → Uvicorn + FastAPI (port 8001)
```

**Why start here:** This is how production systems work (nginx → Gunicorn). Learn real-world patterns before diving into complex FFI.

### Phase 2: Native ASGI (Future)
Implement the ASGI protocol directly in Go, calling Python functions via FFI.

```
Client → Go Server → Python interpreter (via cgo/FFI) → FastAPI
```

**Why this matters:** Deep understanding of how servers and frameworks communicate. More complex but highly educational.

## Current Progress

**✅ Completed:**
- Go reverse proxy server with connection pooling
- FastAPI sample app with test endpoints
- Request/response logging middleware
- Project structure and documentation

**🚧 Next Steps:**
- Add benchmarking suite (compare vs Uvicorn/Gunicorn)
- Implement load balancing across multiple workers
- Research FFI options for Phase 2 (cgo + Python C API)
- Add metrics and observability

## Quick Start

```bash
# Install Python dependencies
cd app
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Run FastAPI backend
uvicorn main:app --port 8001

# Run Go server (in another terminal)
cd ../server
go run main.go
```

Visit `http://localhost:8000` - requests flow through Go to FastAPI.


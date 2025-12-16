from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
import time
import asyncio

app = FastAPI(
    title="go-asgi Sample App",
    description="A sample FastAPI application for testing the go-asgi server",
    version="0.1.0"
)

@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "message": "Hello from FastAPI via go-asgi!",
        "server": "FastAPI + Uvicorn",
        "proxy": "go-asgi reverse proxy"
    }

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "timestamp": time.time()
    }

@app.get("/async/{delay}")
async def async_endpoint(delay: int):
    """
    Async endpoint that simulates work
    Tests async handling through the proxy
    """
    await asyncio.sleep(delay)
    return {
        "message": f"Completed after {delay} seconds",
        "async": True
    }

@app.post("/echo")
async def echo(request: Request):
    """Echo back the request body"""
    body = await request.json()
    return {
        "echoed": body,
        "method": request.method,
        "headers": dict(request.headers)
    }

@app.get("/info")
async def server_info(request: Request):
    """Return information about the request"""
    return {
        "client": request.client.host,
        "method": request.method,
        "url": str(request.url),
        "headers": dict(request.headers),
        "query_params": dict(request.query_params)
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	// Port for the Go server to listen on
	ServerPort = ":8000"
	// Backend ASGI server (Uvicorn running FastAPI)
	BackendURL = "http://localhost:8001"
)

func main() {
	// Parse the backend URL
	backendURL, err := url.Parse(BackendURL)
	if err != nil {
		log.Fatalf("Failed to parse backend URL: %v", err)
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// Customize the proxy transport for better performance
	proxy.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	// Custom error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Bad Gateway - ASGI backend unavailable", http.StatusBadGateway)
	}

	// Create HTTP server with the proxy handler
	server := &http.Server{
		Addr:         ServerPort,
		Handler:      loggingMiddleware(proxy),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 go-asgi server starting on http://localhost%s", ServerPort)
	log.Printf("📡 Proxying to ASGI backend at %s", BackendURL)
	log.Printf("💡 Make sure your FastAPI app is running on port 8001")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// loggingMiddleware logs each request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Log request
		log.Printf("→ %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		// Log response time
		duration := time.Since(start)
		log.Printf("← %s %s completed in %v", r.Method, r.URL.Path, duration)
	})
}

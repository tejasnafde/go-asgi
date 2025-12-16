package main

import (
	"bytes"
	"io"
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
		http.Error(w, "Bad Gateway - go-asgi Server Unavailable", http.StatusBadGateway)
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

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs each request with detailed information
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get IST timezone
		ist, _ := time.LoadLocation("Asia/Kolkata")
		istTime := start.In(ist)

		// Check if system timezone is different from IST
		_, systemOffset := start.Zone()
		_, istOffset := istTime.Zone()

		var timestamp string
		if systemOffset == istOffset {
			// Same timezone, show only one timestamp
			timestamp = istTime.Format("15:04:05")
		} else {
			// Different timezones, show both with labels
			systemTime := start.Format("15:04:05")
			timestamp = systemTime + " (Local) | " + istTime.Format("15:04:05") + " (IST)"
		}

		// Extract IP address (remove port if present)
		ip := r.RemoteAddr
		if colonIndex := len(ip) - 1; colonIndex >= 0 {
			for i := len(ip) - 1; i >= 0; i-- {
				if ip[i] == ':' {
					ip = ip[:i]
					break
				}
			}
		}

		// Log incoming request
		log.SetFlags(0) // Remove default timestamp
		log.Printf("→ [%s] %s %s | IP: %s", timestamp, r.Method, r.URL.Path, ip)

		// Log request payload for POST/PUT/PATCH requests
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			body := make([]byte, r.ContentLength)
			if r.ContentLength > 0 && r.ContentLength < 10000 { // Limit to 10KB
				r.Body.Read(body)
				log.Printf("  📦 Payload: %s", string(body))
				// Restore body for downstream handlers
				r.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Log response with status code and duration
		duration := time.Since(start)
		log.Printf("← [%s] %s %s | Status: %d | Duration: %v",
			timestamp, r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

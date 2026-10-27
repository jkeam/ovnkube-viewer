package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jkeam/ovnkube-viewer/pkg/api/handlers"
	"github.com/jkeam/ovnkube-viewer/pkg/api/middleware"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create handlers
	traceHandler := handlers.NewTraceHandler()

	// Setup router
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/v1/trace", traceHandler.HandleTrace)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware
	handler := middleware.Logging(middleware.CORS(mux))

	// Start server
	addr := ":" + port
	log.Printf("Starting OVN Trace Viewer API server on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jkeam/ovnkube-viewer/pkg/executor"
	"github.com/jkeam/ovnkube-viewer/pkg/models"
	"github.com/jkeam/ovnkube-viewer/pkg/parser"
)

// TraceHandler handles trace requests
type TraceHandler struct {
	executor executor.Executor
	parser   *parser.Parser
}

// NewTraceHandler creates a new TraceHandler
func NewTraceHandler() *TraceHandler {
	// Try to create Kubernetes-based executor
	k8sExecutor, err := executor.NewK8sTraceExecutor()
	if err != nil {
		// Fall back to local executor (which uses mock data)
		log.Printf("WARNING: Failed to create K8s executor, using mock executor: %v", err)
		return &TraceHandler{
			executor: executor.NewTraceExecutor(),
			parser:   parser.NewParser(),
		}
	}

	log.Printf("INFO: Using Kubernetes-based OVN trace executor")
	return &TraceHandler{
		executor: k8sExecutor,
		parser:   parser.NewParser(),
	}
}

// HandleTrace processes a trace request
func (h *TraceHandler) HandleTrace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req models.TraceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate request
	if req.Source == "" || req.Destination == "" || req.Protocol == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "source, destination, and protocol are required"})
		return
	}

	// Execute trace
	output, err := h.executor.Execute(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Parse output
	result, err := h.parser.Parse(output)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse trace output: " + err.Error()})
		return
	}

	// Populate metadata
	result.Metadata.Timestamp = time.Now()
	result.Metadata.Source = req.Source
	result.Metadata.Destination = req.Destination
	result.Metadata.Protocol = req.Protocol

	// Populate packet information in flow steps
	for i := range result.Flow {
		result.Flow[i].Packet.SourceIP = req.Source
		result.Flow[i].Packet.DestIP = req.Destination
		result.Flow[i].Packet.Protocol = req.Protocol
		result.Flow[i].Packet.SourcePort = req.SourcePort
		result.Flow[i].Packet.DestPort = req.DestPort
	}

	// Return result
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jkeam/ovnkube-viewer/pkg/executor"
	"github.com/jkeam/ovnkube-viewer/pkg/models"
	"github.com/jkeam/ovnkube-viewer/pkg/parser"
)

// TraceHandler handles trace requests
type TraceHandler struct {
	executor *executor.TraceExecutor
	parser   *parser.Parser
}

// NewTraceHandler creates a new TraceHandler
func NewTraceHandler() *TraceHandler {
	return &TraceHandler{
		executor: executor.NewTraceExecutor(),
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Source == "" || req.Destination == "" || req.Protocol == "" {
		http.Error(w, "source, destination, and protocol are required", http.StatusBadRequest)
		return
	}

	// Execute trace
	output, err := h.executor.Execute(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse output
	result, err := h.parser.Parse(output)
	if err != nil {
		http.Error(w, "Failed to parse trace output", http.StatusInternalServerError)
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

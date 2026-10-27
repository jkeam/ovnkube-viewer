package executor

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jkeam/ovnkube-viewer/pkg/models"
)

// TraceExecutor executes ovn-trace commands
type TraceExecutor struct {
	ovnTracePath string
}

// NewTraceExecutor creates a new TraceExecutor
func NewTraceExecutor() *TraceExecutor {
	return &TraceExecutor{
		ovnTracePath: "ovn-trace", // This should be configurable
	}
}

// Execute runs an ovn-trace command and returns the output
func (e *TraceExecutor) Execute(req *models.TraceRequest) (string, error) {
	// Validate request
	if req.Source == "" || req.Destination == "" {
		return "", fmt.Errorf("source and destination are required")
	}

	// Build ovn-trace command
	// This is a simplified example - actual implementation will need to:
	// 1. Determine the correct datapath
	// 2. Construct proper packet specification
	// 3. Handle different protocols
	args := e.buildTraceCommand(req)

	// For now, return mock data if ovn-trace is not available
	// In production, this would execute the actual command
	if !e.isOVNTraceAvailable() {
		return e.getMockOutput(req), nil
	}

	cmd := exec.Command(e.ovnTracePath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ovn-trace failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// buildTraceCommand builds the ovn-trace command arguments
func (e *TraceExecutor) buildTraceCommand(req *models.TraceRequest) []string {
	args := []string{}

	// Add datapath (simplified - should be determined from source)
	args = append(args, "ls-switch") // logical switch name

	// Build packet specification
	packet := e.buildPacketSpec(req)
	args = append(args, packet)

	return args
}

// buildPacketSpec builds the packet specification string
func (e *TraceExecutor) buildPacketSpec(req *models.TraceRequest) string {
	parts := []string{}

	// Add source and destination
	parts = append(parts, fmt.Sprintf("inport==\"%s\"", req.Source))

	// Add protocol
	proto := strings.ToLower(req.Protocol)
	parts = append(parts, fmt.Sprintf("%s", proto))

	// Add source IP (simplified)
	parts = append(parts, fmt.Sprintf("ip.src==%s", req.Source))
	parts = append(parts, fmt.Sprintf("ip.dst==%s", req.Destination))

	// Add ports if specified
	if req.SourcePort != nil {
		parts = append(parts, fmt.Sprintf("%s.src==%d", proto, *req.SourcePort))
	}
	if req.DestPort != nil {
		parts = append(parts, fmt.Sprintf("%s.dst==%d", proto, *req.DestPort))
	}

	return strings.Join(parts, ",")
}

// isOVNTraceAvailable checks if ovn-trace command is available
func (e *TraceExecutor) isOVNTraceAvailable() bool {
	_, err := exec.LookPath(e.ovnTracePath)
	return err == nil
}

// getMockOutput returns mock ovn-trace output for development
func (e *TraceExecutor) getMockOutput(req *models.TraceRequest) string {
	return fmt.Sprintf(`# ovn-trace mock output
# Tracing from %s to %s via %s

ingress(dp="ovn-cluster", inport="k8s-%s")
 -> ovn_stage_1
 -> ovn_switch_1 (type: switch)
 -> routing decision: forward to router
egress(dp="ovn-cluster", outport="rtos-ovn_cluster_router")
 -> ovn_router_1 (type: router)
 -> routing table lookup: %s
 -> next hop: gateway
ingress(dp="ovn-cluster", inport="gateway-1")
 -> ovn_switch_2 (type: switch)
 -> destination found
egress(dp="ovn-cluster", outport="k8s-%s")
 -> delivery to %s

Trace completed successfully.
`, req.Source, req.Destination, req.Protocol, req.Source, req.Destination, req.Destination, req.Destination)
}

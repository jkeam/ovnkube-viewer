package parser

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/jkeam/ovnkube-viewer/pkg/models"
)

// Parser handles parsing of ovn-trace output
type Parser struct{}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses ovn-trace output and returns a TraceResult
func (p *Parser) Parse(output string) (*models.TraceResult, error) {
	if output == "" {
		return nil, fmt.Errorf("empty trace output")
	}

	scanner := bufio.NewScanner(strings.NewReader(output))

	result := &models.TraceResult{
		Topology: models.TopologyGraph{
			Nodes: []models.TopologyNode{},
			Edges: []models.TopologyEdge{},
		},
		Flow:     []models.FlowStep{},
		Metadata: models.TraceMetadata{
			Status: "success",
		},
	}

	stepNum := 0
	nodeMap := make(map[string]bool)

	// Parse line by line
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// TODO: Implement actual parsing logic based on ovn-trace output format
		// This is a placeholder implementation

		// Example: detect ingress/egress patterns
		if strings.Contains(line, "ingress") || strings.Contains(line, "egress") {
			stepNum++

			// Extract node information
			nodeName := p.extractNodeName(line)
			nodeID := fmt.Sprintf("node-%d", stepNum)

			if !nodeMap[nodeID] {
				nodeType := p.detectNodeType(line)
				node := models.TopologyNode{
					ID:   nodeID,
					Type: nodeType,
					Name: nodeName,
					Metadata: map[string]interface{}{
						"raw": line,
					},
				}
				result.Topology.Nodes = append(result.Topology.Nodes, node)
				nodeMap[nodeID] = true
			}

			// Create flow step
			step := models.FlowStep{
				StepNum:     stepNum,
				NodeID:      nodeID,
				Action:      p.extractAction(line),
				Description: line,
				Packet: models.PacketState{
					Protocol:      "tcp",
					Modifications: []string{},
				},
			}
			result.Flow = append(result.Flow, step)

			// Create edge if there's a previous node
			if stepNum > 1 {
				edge := models.TopologyEdge{
					Source: fmt.Sprintf("node-%d", stepNum-1),
					Target: nodeID,
				}
				result.Topology.Edges = append(result.Topology.Edges, edge)
			}
		}

		// Detect dropped packets
		if strings.Contains(strings.ToLower(line), "drop") {
			result.Metadata.Status = "dropped"
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning output: %w", err)
	}

	return result, nil
}

// extractNodeName extracts the node name from a trace line
func (p *Parser) extractNodeName(line string) string {
	// TODO: Implement proper node name extraction
	parts := strings.Fields(line)
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

// detectNodeType determines the type of OVN component from the line
func (p *Parser) detectNodeType(line string) string {
	lineLower := strings.ToLower(line)

	if strings.Contains(lineLower, "switch") {
		return "switch"
	}
	if strings.Contains(lineLower, "router") {
		return "router"
	}
	if strings.Contains(lineLower, "port") {
		return "port"
	}
	if strings.Contains(lineLower, "acl") {
		return "acl"
	}

	return "unknown"
}

// extractAction extracts the action from a trace line
func (p *Parser) extractAction(line string) string {
	lineLower := strings.ToLower(line)

	if strings.Contains(lineLower, "ingress") {
		return "ingress"
	}
	if strings.Contains(lineLower, "egress") {
		return "egress"
	}
	if strings.Contains(lineLower, "route") {
		return "route"
	}
	if strings.Contains(lineLower, "drop") {
		return "drop"
	}
	if strings.Contains(lineLower, "output") {
		return "output"
	}

	return "forward"
}

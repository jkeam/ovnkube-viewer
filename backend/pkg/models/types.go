package models

import "time"

// TopologyNode represents a node in the OVN topology
type TopologyNode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // "switch", "router", "port", "acl"
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

// PacketState represents the state of a packet at a given step
type PacketState struct {
	SourceIP      string   `json:"sourceIP"`
	DestIP        string   `json:"destIP"`
	Protocol      string   `json:"protocol"`
	SourcePort    *int     `json:"sourcePort,omitempty"`
	DestPort      *int     `json:"destPort,omitempty"`
	Modifications []string `json:"modifications"`
}

// FlowStep represents a single step in the packet flow
type FlowStep struct {
	StepNum     int         `json:"stepNum"`
	NodeID      string      `json:"nodeID"`
	Action      string      `json:"action"`
	Description string      `json:"description"`
	Packet      PacketState `json:"packet"`
}

// TraceMetadata contains metadata about the trace execution
type TraceMetadata struct {
	Timestamp   time.Time `json:"timestamp"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Protocol    string    `json:"protocol"`
	Status      string    `json:"status"` // "success", "dropped", "error"
}

// TopologyEdge represents a connection between nodes
type TopologyEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label,omitempty"`
}

// TopologyGraph represents the complete network topology
type TopologyGraph struct {
	Nodes []TopologyNode `json:"nodes"`
	Edges []TopologyEdge `json:"edges"`
}

// TraceResult is the complete result of a trace operation
type TraceResult struct {
	Topology TopologyGraph `json:"topology"`
	Flow     []FlowStep    `json:"flow"`
	Metadata TraceMetadata `json:"metadata"`
}

// TraceRequest represents a trace execution request
type TraceRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Protocol    string `json:"protocol"`
	SourcePort  *int   `json:"sourcePort,omitempty"`
	DestPort    *int   `json:"destPort,omitempty"`
}

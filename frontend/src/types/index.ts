// Core data models matching backend types

export interface TopologyNode {
  id: string;
  type: 'switch' | 'router' | 'port' | 'acl';
  name: string;
  metadata: Record<string, unknown>;
}

export interface PacketState {
  sourceIP: string;
  destIP: string;
  protocol: string;
  sourcePort?: number;
  destPort?: number;
  modifications: string[];
}

export interface FlowStep {
  stepNum: number;
  nodeID: string;
  action: string;
  description: string;
  packet: PacketState;
}

export interface TraceMetadata {
  timestamp: string;
  source: string;
  destination: string;
  protocol: string;
  status: 'success' | 'dropped' | 'error';
}

export interface TopologyGraph {
  nodes: TopologyNode[];
  edges: Array<{
    source: string;
    target: string;
    label?: string;
  }>;
}

export interface TraceResult {
  topology: TopologyGraph;
  flow: FlowStep[];
  metadata: TraceMetadata;
}

export interface TraceRequest {
  source: string;
  destination: string;
  protocol: string;
  sourcePort?: number;
  destPort?: number;
}

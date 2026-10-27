# OVN-Kubernetes Trace Viewer - Implementation Plan

## Project Overview

A visual viewer for OVN-Kubernetes traffic tracing that transforms text-based ovn-trace output into interactive flow diagrams. The project will evolve from a standalone application to a full OpenShift operator with web console integration.

**Reference Documentation**: [OpenShift OVN-Kubernetes Tracing](https://docs.redhat.com/en/documentation/openshift_container_platform/4.22/html/ovn-kubernetes_network_plugin/ovn-kubernetes-tracing-using-ovntrace)

---

## Architecture Overview

### Phase 1: Standalone Application

```
┌─────────────────────────────────────────────────────────────────┐
│                         PHASE 1: Standalone                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────┐          ┌──────────────────────────────┐    │
│  │   Web UI     │          │      Backend API Server      │    │
│  │  (React)     │◄────────►│         (Go)                 │    │
│  │              │   REST/  │                              │    │
│  │ - Input form │   WS     │ - Parse ovn-trace output     │    │
│  │ - Flow viz   │          │ - Build topology graph       │    │
│  │   (ReactFlow)│          │ - Execute trace commands     │    │
│  └──────────────┘          │ - Transform to viz format    │    │
│                             └──────────────────────────────┘    │
│                                          │                       │
│                                          ▼                       │
│                             ┌──────────────────────────────┐    │
│                             │   OVN/OVS Components         │    │
│                             │   (via ovn-trace CLI)        │    │
│                             └──────────────────────────────┘    │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### Phase 2: OpenShift Operator + Console Plugin

```
┌─────────────────────────────────────────────────────────────────┐
│              PHASE 2: OpenShift Operator + Plugin                │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │          OpenShift Web Console                             │ │
│  │  ┌──────────────────────────────────────────────────────┐  │ │
│  │  │    OVN Trace Viewer Plugin (Dynamic Plugin)          │  │ │
│  │  │    - Integrated navigation                           │  │ │
│  │  │    - Access via Networking → OVN Tracing            │  │ │
│  │  │    - Reuses Phase 1 React components               │  │ │
│  │  └──────────────────────────────────────────────────────┘  │ │
│  └──────┬─────────────────────────────────────────────────────┘ │
│         │                                                         │
│         │ K8s API (CRDs)                                         │
│         ▼                                                         │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              OVN Trace Operator                           │  │
│  │                                                            │  │
│  │  ┌─────────────────┐      ┌─────────────────────────┐   │  │
│  │  │   Controller    │      │   Trace Executor        │   │  │
│  │  │   (Reconcile    │─────►│   - Run ovn-trace       │   │  │
│  │  │   TraceRequest) │      │   - Parse output        │   │  │
│  │  └─────────────────┘      │   - Store results       │   │  │
│  │          │                 └─────────────────────────┘   │  │
│  │          │                                                 │  │
│  │          ▼                                                 │  │
│  │  ┌─────────────────────────────────────────────────────┐ │  │
│  │  │   Custom Resource Definitions (CRDs)                │ │  │
│  │  │   - TraceRequest  (trace parameters)                │ │  │
│  │  │   - TraceResult   (parsed topology + flow)          │ │  │
│  │  └─────────────────────────────────────────────────────┘ │  │
│  └──────────────────────────────────────────────────────────┘  │
│                              │                                   │
│                              ▼                                   │
│                   ┌────────────────────────┐                    │
│                   │  OVN-Kubernetes Pods   │                    │
│                   │  (ovnkube-master/node) │                    │
│                   └────────────────────────┘                    │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### Data Flow

```
Input (User) → Parse → Model → Transform → Visualize
                ↓         ↓         ↓          ↓
            ovn-trace  Topology  Graph    ReactFlow
             output     Graph    Layout    Render
```

---

## Technology Stack

### Backend
- **Language**: Go
- **Framework**: Standard library + gorilla/mux or gin for routing
- **Parser**: Custom lexer/parser for ovn-trace output
- **Graph**: Custom topology graph builder
- **Operator**: Operator SDK (for Phase 3+)

### Frontend
- **Framework**: React + TypeScript
- **Build Tool**: Vite
- **Visualization**: ReactFlow
- **UI Components**: 
  - Phase 1: Tailwind CSS or Material-UI
  - Phase 4: PatternFly (OpenShift console)
- **State Management**: React Context + hooks

### Deployment
- **Containerization**: Docker
- **Orchestration**: Kubernetes/OpenShift
- **Package Manager**: Helm (Phase 2+)
- **Operator Framework**: Operator SDK (Phase 3+)

---

## Data Models

### Core Types

```go
type TopologyNode struct {
    ID       string                    // Unique identifier
    Type     string                    // "switch", "router", "port", "acl"
    Name     string                    // Component name
    Metadata map[string]interface{}    // Additional properties
}

type FlowStep struct {
    StepNum     int
    NodeID      string
    Action      string                 // "ingress", "egress", "route", "drop", etc.
    Description string
    Packet      PacketState
}

type PacketState struct {
    SourceIP      string
    DestIP        string
    Protocol      string
    SourcePort    int
    DestPort      int
    Modifications []string            // NAT, TTL changes, etc.
}

type TraceResult struct {
    Topology *Graph
    Flow     []FlowStep
    Metadata TraceMetadata
}

type TraceMetadata struct {
    Timestamp   time.Time
    Source      string
    Destination string
    Protocol    string
    Status      string              // "success", "dropped", "error"
}
```

### CRDs (Phase 3)

```yaml
apiVersion: ovn.openshift.io/v1alpha1
kind: TraceRequest
metadata:
  name: pod-to-service-trace
spec:
  source:
    namespace: default
    pod: web-pod
  destination:
    service: backend-svc
    port: 8080
  protocol: tcp
status:
  state: Completed|Running|Failed
  resultRef: trace-result-xyz
  
---
apiVersion: ovn.openshift.io/v1alpha1
kind: TraceResult
metadata:
  name: trace-result-xyz
spec:
  topology:
    nodes: [...]
    edges: [...]
  flow: [...]
  metadata:
    timestamp: "2026-10-26T14:00:00Z"
    status: "success"
```

---

## Implementation Phases

### Phase 0: Research & Prototyping (1-2 weeks)

**Goal**: Understand ovn-trace output format and validate parsing approach

**Tasks**:
1. Collect sample ovn-trace outputs from various scenarios:
   - Pod-to-pod within namespace
   - Pod-to-service
   - Ingress/egress traffic
   - NetworkPolicy enforcement scenarios
   - Cross-namespace traffic
   - Drop scenarios

2. Build prototype parser (Go):
   - Identify patterns in ovn-trace output
   - Extract topology elements (switches, routers, ports)
   - Map flow steps to topology nodes
   - Handle edge cases (drops, NAT, load balancing)

3. Create simple proof-of-concept visualization:
   - Static HTML + D3.js or Mermaid
   - Render one sample trace as a graph
   - Validate that visualization conveys meaningful information

**Deliverable**: Working parser + static visualization of sample traces

---

### Phase 1: Standalone MVP (4-6 weeks)

**Goal**: Production-ready standalone application

#### Milestone 1.1: Backend Core (2 weeks)

- **Parser library**: Robust ovn-trace parser with tests
- **Graph builder**: Convert parsed data to topology graph
- **API server**: REST endpoints for trace execution
  - `POST /api/v1/trace` - Execute trace with parameters
  - `GET /api/v1/trace/:id` - Retrieve trace results
  - `GET /api/v1/topology` - Get cached topology
- **OVN integration**: Execute ovn-trace commands safely
- **Testing**: Unit tests for parser, integration tests for API

#### Milestone 1.2: Frontend Core (2 weeks)

- **React app setup**: TypeScript, Vite build
- **Trace input form**: 
  - Source/destination specification
  - Protocol/port configuration
  - Advanced options (datapath, detailed output)
- **Flow visualization**: 
  - ReactFlow integration
  - Custom node types for OVN components
  - Edge styling for different flow types
  - Interactive node selection/highlighting
- **Step-by-step view**: Timeline/stepper showing flow progression
- **Responsive design**: Works on desktop/tablet

#### Milestone 1.3: Integration & Polish (1-2 weeks)

- **Live tracing**: WebSocket for real-time trace execution
- **Error handling**: User-friendly error messages
- **Export functionality**: Save traces as JSON/SVG/PNG
- **Documentation**: User guide, API documentation
- **Containerization**: Dockerfile, docker-compose for easy deployment
- **Configuration**: Environment-based config (API endpoint, auth, etc.)

**Deployment**:
```yaml
# docker-compose.yaml
services:
  backend:
    image: ovnkube-viewer-backend:latest
    ports: ["8080:8080"]
    volumes: ["/var/run/openvswitch:/var/run/openvswitch:ro"]
  
  frontend:
    image: ovnkube-viewer-frontend:latest
    ports: ["3000:3000"]
    environment:
      VITE_API_URL: http://backend:8080
```

**MVP Feature Set**:

Must Have:
1. Parse ovn-trace output (basic scenarios)
2. Visualize flow as interactive graph
3. Show packet transformations at each step
4. Basic input form (source/dest, protocol)
5. Export trace as JSON
6. Containerized deployment

Should Have:
7. Real-time trace execution
8. Error handling and validation
9. Responsive design
10. Basic documentation

Nice to Have:
11. Trace history (in-memory)
12. Dark mode
13. Export as SVG/PNG
14. Sample traces for demo

---

### Phase 2: Kubernetes Native (4-6 weeks)

**Goal**: Deploy on Kubernetes, add persistence and multi-user support

#### Milestone 2.1: Kubernetes Deployment (1-2 weeks)

- **Helm chart**: Parameterized deployment
- **Service/Ingress**: Expose via Kubernetes services
- **RBAC**: ServiceAccount with permissions to access OVN resources
- **ConfigMaps/Secrets**: External configuration
- **Persistent storage**: Optional PVC for trace history

#### Milestone 2.2: Enhanced Features (2 weeks)

- **Trace history**: Store traces in database (PostgreSQL or in-memory)
- **Multi-tenancy**: Namespace isolation
- **Batch tracing**: Queue multiple trace requests
- **Comparison view**: Compare two traces side-by-side
- **Filtering**: Filter topology by namespace, labels

#### Milestone 2.3: Authentication & Authorization (1-2 weeks)

- **ServiceAccount tokens**: Kubernetes-native auth
- **RBAC policies**: Define who can trace what
- **Audit logging**: Record trace requests for compliance

---

### Phase 3: OpenShift Operator (6-8 weeks)

**Goal**: Full operator pattern with CRDs and automated lifecycle management

#### Milestone 3.1: Operator Foundation (2-3 weeks)

- **Operator SDK setup**: Initialize Go-based operator
- **CRD definitions**: TraceRequest, TraceResult
- **Controller reconciliation loop**:
  - Watch TraceRequest resources
  - Execute traces via Jobs/Pods
  - Store results in TraceResult resources
  - Update status and conditions

#### Milestone 3.2: Operator Features (2 weeks)

- **Scheduled tracing**: CronJob-based periodic traces
- **Webhooks**: Validation/mutation webhooks for TraceRequest
- **Metrics**: Prometheus metrics for operator health
- **Events**: Emit Kubernetes events for trace lifecycle
- **Garbage collection**: Auto-cleanup old TraceResults

#### Milestone 3.3: Advanced Operator Capabilities (2-3 weeks)

- **Topology caching**: Operator maintains cluster topology
- **Diff detection**: Notify on topology changes
- **Policy validation**: Suggest NetworkPolicy issues
- **Auto-remediation**: Optional auto-fixes for common issues

---

### Phase 4: OpenShift Console Plugin (4-6 weeks)

**Goal**: Integrated UI in OpenShift web console

#### Milestone 4.1: Plugin Foundation (2 weeks)

- **Dynamic plugin setup**: OpenShift Dynamic Plugin SDK
- **Plugin registration**: ConsolePlugin CR, Nginx serving
- **Navigation integration**: Add "OVN Tracing" to Networking menu
- **PatternFly components**: Use OpenShift design system
- **Auth integration**: Use console's auth tokens

#### Milestone 4.2: Console Integration (2 weeks)

- **Resource pages**: View TraceRequest/TraceResult in console
- **Action menus**: "Trace from this Pod" on Pod details
- **List views**: Table of recent traces with filters
- **Form integration**: Create TraceRequest via console form
- **Status indicators**: Show trace state in resource lists

#### Milestone 4.3: Advanced Console Features (2 weeks)

- **Embedded visualization**: Inline flow diagrams in resource pages
- **Quick actions**: One-click trace from Pod/Service pages
- **Topology integration**: Link to console topology view
- **Notifications**: Toast notifications for trace completion
- **Help integration**: Context-sensitive help

**Console Plugin Deployment**:
```yaml
apiVersion: console.openshift.io/v1alpha1
kind: ConsolePlugin
metadata:
  name: ovn-trace-viewer
spec:
  displayName: OVN Trace Viewer
  service:
    name: ovn-trace-viewer-plugin
    namespace: ovn-kubernetes
    port: 9443
    basePath: /
```

---

## Project Directory Structure

```
ovnkube-viewer/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── pkg/
│   │   ├── parser/           # ovn-trace parser
│   │   │   ├── lexer.go
│   │   │   ├── parser.go
│   │   │   └── parser_test.go
│   │   ├── graph/            # Topology graph builder
│   │   │   ├── builder.go
│   │   │   ├── layout.go
│   │   │   └── types.go
│   │   ├── executor/         # Execute ovn-trace commands
│   │   │   └── trace.go
│   │   ├── api/              # HTTP API handlers
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── server.go
│   │   └── models/           # Data models
│   │       └── types.go
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── TraceForm/
│   │   │   ├── FlowVisualization/
│   │   │   ├── TopologyView/
│   │   │   └── shared/
│   │   ├── hooks/
│   │   ├── utils/
│   │   ├── types/
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── Dockerfile
│
├── operator/                  # Phase 3
│   ├── api/
│   │   └── v1alpha1/
│   │       ├── tracerequest_types.go
│   │       └── traceresult_types.go
│   ├── controllers/
│   │   ├── tracerequest_controller.go
│   │   └── topology_controller.go
│   ├── config/
│   │   ├── crd/
│   │   ├── rbac/
│   │   └── manager/
│   ├── Dockerfile
│   └── Makefile
│
├── console-plugin/           # Phase 4
│   ├── src/
│   │   ├── components/       # Reuse from frontend/
│   │   ├── extensions.ts
│   │   └── index.ts
│   ├── console-extensions.json
│   ├── package.json
│   └── Dockerfile
│
├── deploy/
│   ├── standalone/
│   │   └── docker-compose.yaml
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── ingress.yaml
│   ├── helm/
│   │   └── ovnkube-viewer/
│   │       ├── Chart.yaml
│   │       ├── values.yaml
│   │       └── templates/
│   └── openshift/
│       ├── operator.yaml
│       └── console-plugin.yaml
│
├── examples/
│   ├── sample-traces/        # Sample ovn-trace outputs
│   └── demo-scenarios/       # Test scenarios
│
├── docs/
│   ├── architecture.md
│   ├── user-guide.md
│   ├── api-reference.md
│   ├── operator-guide.md
│   ├── development.md
│   └── implementation-plan.md  # This file
│
├── .github/
│   └── workflows/
│       ├── backend-ci.yaml
│       ├── frontend-ci.yaml
│       └── operator-ci.yaml
│
├── README.md
├── LICENSE
└── Makefile
```

---

## Key Technical Decisions

### 1. Parser Implementation
- **Choice**: Custom parser in Go
- **Rationale**: ovn-trace output is specific, no standard parser exists. Go provides strong parsing capabilities and aligns with K8s ecosystem.

### 2. Visualization Library
- **Choice**: ReactFlow
- **Rationale**: React-native, excellent TypeScript support, built-in interactive features (zoom, pan, dragging), easy customization.

### 3. Backend Language
- **Choice**: Go
- **Rationale**: Standard for K8s operators, excellent performance, single binary deployment, strong typing.

### 4. State Management
- **Phase 1**: In-memory (ephemeral)
- **Phase 2**: Optional PostgreSQL (persistent)
- **Phase 3**: Kubernetes CRDs (operator pattern)
- **Rationale**: Progressive enhancement based on deployment complexity.

### 5. Graph Layout Algorithm
- **Choice**: Automatic with Dagre (hierarchical layout)
- **Rationale**: OVN flows are naturally hierarchical. Dagre produces clean left-to-right or top-to-bottom flows.

### 6. Console Plugin Strategy
- **Choice**: Build standalone first, wrap in console plugin
- **Rationale**: Broader applicability, 80% component reuse, better modularity.

---

## OpenShift-Specific Considerations

### 1. OVN-Kubernetes Access
- OVN databases run in `openshift-ovn-kubernetes` namespace
- Need privileged access to ovn-nbctl/ovn-trace
- **Approach**: ServiceAccount with cluster-reader role + exec into ovnkube-master pods or DaemonSet with host network

### 2. Console Plugin Requirements
- Must serve over HTTPS (use OpenShift service CA)
- Must register via ConsolePlugin CR
- Must use PatternFly components for consistency
- Must handle console OAuth tokens

### 3. Security Context Constraints (SCC)
- Need appropriate SCC for OVN access
- Likely need `privileged` or custom SCC
- Document security implications

### 4. Multi-tenancy
- Users should only trace within authorized namespaces
- RBAC: Create custom roles (ovn-trace-viewer, ovn-trace-admin)
- Operator validates permissions before executing traces

### 5. Performance Considerations
- OVN databases can be large in big clusters
- Cache topology data (refresh periodically)
- Limit concurrent trace requests
- Add request timeouts

---

## Risk Mitigation

### Risk 1: ovn-trace Output Format Changes
- **Mitigation**: Version parser, add format detection
- **Testing**: Collect samples from multiple OCP versions

### Risk 2: Complex Topologies → Unreadable Graphs
- **Mitigation**: 
  - Implement filtering (by namespace, component type)
  - Add collapsible node groups
  - Provide multiple layout options

### Risk 3: Performance with Large Clusters
- **Mitigation**:
  - Lazy loading of topology
  - Paginate trace history
  - Add request queuing

### Risk 4: Operator SDK Breaking Changes
- **Mitigation**: Pin Operator SDK version, regular updates

### Risk 5: Console Plugin API Changes
- **Mitigation**: Track OpenShift release notes, version compatibility matrix

---

## Success Metrics

### Phase 1
- Successfully parse 95%+ of common trace scenarios
- Render graph in <2s for typical traces
- User can trace and visualize without documentation

### Phase 2
- Deploy on Kubernetes in <5 minutes
- Support 100+ concurrent users
- 99% uptime

### Phase 3
- Operator passes OpenShift certification
- CR API is stable (v1)
- < 50MB memory footprint per operator pod

### Phase 4
- Plugin loads in console without errors
- Integrated with 80%+ of relevant console pages
- Matches OpenShift UX patterns

---

## Next Steps

### Immediate Actions

1. **Set up project structure**:
   ```bash
   mkdir -p backend/{cmd/server,pkg/{parser,graph,executor,api,models}}
   mkdir -p frontend/src/{components,hooks,utils,types}
   mkdir -p deploy/{standalone,kubernetes,helm,openshift}
   mkdir -p examples/{sample-traces,demo-scenarios}
   ```

2. **Initialize Go module**:
   ```bash
   cd backend && go mod init github.com/jkeam/ovnkube-viewer
   ```

3. **Initialize React app**:
   ```bash
   cd frontend && npm create vite@latest . -- --template react-ts
   ```

4. **Collect sample ovn-trace outputs**:
   - From test OCP cluster
   - Various scenarios (pod-pod, svc, ingress, drop)
   - Store in `examples/sample-traces/`

5. **Start with parser prototype**:
   - Build lexer to tokenize ovn-trace output
   - Identify key patterns (ingress/egress, routing, ACLs)
   - Test with sample outputs

6. **Create basic visualization POC**:
   - Simple React app
   - Load sample JSON (manually parsed trace)
   - Render with ReactFlow
   - Validate graph conveys information

### Critical Files for Initial Implementation

1. `backend/pkg/parser/parser.go` - Core ovn-trace output parser
2. `backend/pkg/graph/builder.go` - Topology graph construction
3. `backend/pkg/models/types.go` - Shared data models
4. `frontend/src/components/FlowVisualization/FlowGraph.tsx` - Main visualization component
5. `frontend/src/hooks/useTrace.ts` - API integration hook

---

## Timeline Summary

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| Phase 0 | 1-2 weeks | Parser prototype + static viz POC |
| Phase 1 | 4-6 weeks | Standalone containerized app |
| Phase 2 | 4-6 weeks | Kubernetes-native with Helm |
| Phase 3 | 6-8 weeks | OpenShift operator with CRDs |
| Phase 4 | 4-6 weeks | Console plugin integration |
| **Total** | **19-28 weeks** | Full OpenShift-integrated solution |

---

## Resources

- [OVN-Kubernetes Documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/4.22/html/ovn-kubernetes_network_plugin/ovn-kubernetes-tracing-using-ovntrace)
- [Operator SDK](https://sdk.operatorframework.io/)
- [OpenShift Console Plugin SDK](https://github.com/openshift/console)
- [ReactFlow Documentation](https://reactflow.dev/)
- [PatternFly Design System](https://www.patternfly.org/)

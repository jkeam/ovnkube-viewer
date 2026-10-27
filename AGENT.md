# AGENT.md - Context for AI Agents

This document provides context for AI agents (including Claude) working on the OVN-Kubernetes Trace Viewer project.

## Project Overview

**Name**: OVN-Kubernetes Trace Viewer  
**Purpose**: Visual viewer that transforms text-based ovn-trace output into interactive flow diagrams  
**Current Phase**: Phase 1 - Standalone MVP (30% complete)  
**Started**: October 26, 2026  
**Tech Stack**: Go 1.21+ (backend), React + TypeScript + ReactFlow on Node.js 24.22.1+ LTS (frontend), Docker

### What This Project Does

Helps users visualize and debug network traffic flow in OVN-Kubernetes environments by:
1. Taking ovn-trace command output (text)
2. Parsing the topology and flow steps
3. Rendering an interactive graph showing packet flow through OVN components
4. Displaying step-by-step packet transformations

### Evolution Path

- **Phase 1** (Current): Standalone containerized application
- **Phase 2** (Future): Kubernetes-native with Helm charts
- **Phase 3** (Future): OpenShift Operator with CRDs
- **Phase 4** (Future): OpenShift Console Plugin

## Current State

### ✅ What's Built

1. **Complete project structure** - All directories and scaffolding in place
2. **Backend (Go)** - API server with basic REST endpoints, mock executor, parser foundation
3. **Frontend (React)** - Input form, ReactFlow visualization, basic UI
4. **Deployment** - Dockerfiles, docker-compose, Makefile, run.sh script
5. **Documentation** - README, GETTING_STARTED, PHASE1_STATUS, implementation plan

### 🚧 What's In Progress

1. **Parser** - Basic structure exists, needs enhancement for real ovn-trace output
2. **Graph builder** - Not yet implemented, needs Dagre layout algorithm
3. **Testing** - No tests yet, TDD needed going forward

### ❌ What's Not Built Yet

1. **Real OVN integration** - Currently using mock data
2. **Comprehensive parser** - Only basic placeholder parsing
3. **Graph layout** - Nodes are positioned simply, not intelligently
4. **Tests** - No unit tests, integration tests, or e2e tests
5. **WebSocket** - Only HTTP REST, no real-time updates

## Architecture

### Data Flow

```
User Input → API Server → Trace Executor → ovn-trace command
                                ↓
                          Raw text output
                                ↓
                            Parser
                                ↓
                          TraceResult (JSON)
                                ↓
                          Frontend/ReactFlow
                                ↓
                        Interactive Graph
```

### Key Components

**Backend**:
- `cmd/server/main.go` - Entry point, HTTP server setup
- `pkg/models/types.go` - Core data structures (TraceResult, TopologyNode, FlowStep)
- `pkg/parser/parser.go` - Parses ovn-trace text output (NEEDS WORK)
- `pkg/executor/trace.go` - Executes ovn-trace commands (currently mocks)
- `pkg/api/handlers/trace.go` - HTTP request handler for /api/v1/trace
- `pkg/api/middleware/` - CORS and logging middleware

**Frontend**:
- `src/App.tsx` - Main component, state management
- `src/components/TraceForm/` - Input form for trace parameters
- `src/components/FlowVisualization/` - ReactFlow graph rendering
- `src/types/index.ts` - TypeScript type definitions (mirror backend)

**Deployment**:
- `run.sh` - Convenience script for running the app (USE THIS!)
- `Makefile` - Build and dev commands
- `deploy/standalone/docker-compose.yaml` - Container orchestration

## Important Files

### Must-Read Before Changes

1. **docs/implementation-plan.md** - Complete roadmap, all phases detailed
2. **PHASE1_STATUS.md** - Current progress, next milestones
3. **pkg/models/types.go** - Data model changes affect both backend and frontend
4. **src/types/index.ts** - Frontend types must match backend models

### Sample Data

- `examples/sample-traces/pod-to-pod.txt` - Example ovn-trace output
- `examples/sample-traces/pod-to-service.txt` - Service with load balancing example

These are **real-ish** sample outputs that the parser needs to handle.

## Development Workflow

### Running the App

**Easiest way**:
```bash
./run.sh init      # First time only
./run.sh dev       # Development mode with hot reload
./run.sh docker    # Docker Compose mode
./run.sh stop      # Stop everything
```

**Testing locally**:
```bash
# Backend only
cd backend && go run ./cmd/server

# Frontend only
cd frontend && npm run dev

# Build
make backend-build
make frontend-build
```

### Making Changes

**Backend changes**:
1. Edit Go files in `backend/pkg/`
2. Go auto-reloads in dev mode
3. Test with `curl` or frontend
4. Add tests in `*_test.go` files (PLEASE ADD TESTS!)

**Frontend changes**:
1. Edit TypeScript files in `frontend/src/`
2. Vite hot-reloads automatically
3. Check browser console for errors
4. Types must match backend models

**Data model changes**:
1. Update `backend/pkg/models/types.go`
2. Update `frontend/src/types/index.ts` to match
3. Update any handlers/components using those types
4. Consider backwards compatibility

## Key Technical Decisions

### Why Go for Backend?

- Standard in Kubernetes ecosystem (Phase 3 will be an operator)
- Excellent performance for parsing
- Single binary deployment
- Strong typing

### Why ReactFlow?

- Best React library for interactive graphs
- Built-in zoom, pan, node dragging
- Easy custom node types
- TypeScript support

### Why Mock Executor?

- Development without needing actual OVN environment
- Faster iteration
- Can still build and test entire stack
- Real executor will drop in later

### Why No Database?

Phase 1 is stateless/ephemeral. Phase 2 will add optional persistence.

## Common Tasks

### Adding a New API Endpoint

1. Add handler function to `backend/pkg/api/handlers/`
2. Register route in `backend/cmd/server/main.go`
3. Update frontend to call new endpoint
4. Document in API docs (when they exist)

### Enhancing the Parser

**CRITICAL**: This is the next major task!

1. Study real ovn-trace outputs in `examples/sample-traces/`
2. Identify patterns (ingress, egress, routing, drops, ACLs)
3. Build lexer to tokenize line-by-line
4. Extract node types, actions, packet states
5. Create comprehensive test suite
6. Test against multiple scenarios

See `pkg/parser/parser.go` TODOs for specific areas.

### Adding a New React Component

1. Create in `frontend/src/components/ComponentName/`
2. Include `ComponentName.tsx` and `ComponentName.css`
3. Import and use in parent component
4. Follow existing patterns (props interface, state management)

### Adding Sample Traces

1. Get real ovn-trace output from OCP cluster
2. Add to `examples/sample-traces/scenario-name.txt`
3. Add comments explaining the scenario
4. Use in parser tests

## Patterns and Conventions

### Go Code

- Package names: lowercase, no underscores
- Exported types: PascalCase
- Private types: camelCase
- Error handling: explicit returns, no panics in handlers
- JSON tags: camelCase (matches JavaScript convention)

### TypeScript Code

- File naming: PascalCase for components (TraceForm.tsx)
- Interfaces: PascalCase with Props suffix for component props
- State: useState hooks, minimize prop drilling
- Types: Keep in sync with backend Go structs

### Git

- Not initialized yet - will need git init
- .gitignore is ready
- No commit conventions established yet

## Next Priorities (in order)

### Immediate (This Week)

1. **Collect real ovn-trace samples** from working OCP cluster
2. **Analyze ovn-trace format** - document patterns and structure
3. **Build robust parser** - handle all common scenarios
4. **Add parser tests** - TDD from here forward

### Short Term (Next 2 Weeks)

1. **Implement graph builder** - convert parsed data to graph
2. **Add Dagre layout** - intelligent node positioning
3. **Enhance visualization** - custom node types, better styling
4. **Real OVN executor** - replace mock with actual ovn-trace execution

### Medium Term (3-6 Weeks)

1. **Comprehensive testing** - unit, integration, e2e
2. **Error handling** - user-friendly messages
3. **Export functionality** - JSON, SVG, PNG
4. **Documentation** - API docs, user guide

## Important Context

### OVN-Kubernetes Background

- **OVN** = Open Virtual Network (part of Open vSwitch)
- **ovn-trace** = CLI tool that simulates packet flow through OVN logical networks
- Output shows: ingress, egress, routing decisions, ACL checks, NAT, load balancing
- Used for debugging connectivity issues in Kubernetes networking

### Target Users

- **Platform engineers** debugging network issues in OpenShift/K8s
- **Support teams** troubleshooting customer problems
- **Developers** understanding how their apps' traffic flows
- **Anyone** who finds ovn-trace text output hard to parse

### Why This Matters

ovn-trace output is powerful but hard to read. A visual graph makes it:
- Easier to spot where packets are dropped
- Clearer how NetworkPolicies affect traffic
- Faster to understand complex multi-hop flows
- Better for sharing with teammates

## Common Pitfalls

### Parser Pitfalls

- ovn-trace format may vary by OVN version - need to handle multiple formats
- Some lines are metadata, some are flow steps - distinguish carefully
- Dropped packets need special handling - status should be "dropped"
- NAT and load balancing add complexity - packet state changes

### Graph Pitfalls

- Large topologies can be overwhelming - consider filtering/collapsing
- Node positioning affects readability - Dagre helps but may need tweaking
- Edge labels can overlap - ReactFlow has options for this
- Performance with 100+ nodes - ReactFlow handles it, but test

### Frontend Pitfalls

- Types must match backend exactly - easy to get out of sync
- ReactFlow requires specific node/edge structure - follow their docs
- State management gets complex - consider Context or Zustand if needed

## Testing Strategy

### Backend Tests (TODO)

```go
// pkg/parser/parser_test.go
func TestParsePodToPod(t *testing.T) {
    p := NewParser()
    output, _ := os.ReadFile("../../examples/sample-traces/pod-to-pod.txt")
    result, err := p.Parse(string(output))
    
    assert.Nil(t, err)
    assert.Equal(t, "success", result.Metadata.Status)
    assert.Greater(t, len(result.Flow), 0)
}
```

### Frontend Tests (TODO)

- Component tests with React Testing Library
- Mock API responses
- Test user interactions (form submit, graph interactions)

### Integration Tests (TODO)

- Full API request → response cycle
- Test with docker-compose
- Verify frontend can consume backend responses

## Debugging Tips

### Backend Issues

```bash
# Check if server is running
curl http://localhost:8080/health

# Test trace endpoint
curl -X POST http://localhost:8080/api/v1/trace \
  -H "Content-Type: application/json" \
  -d '{"source":"pod1","destination":"pod2","protocol":"tcp"}'

# View logs in Docker mode
./run.sh logs
```

### Frontend Issues

- Open browser DevTools → Network tab
- Check for 404s or CORS errors
- Console tab shows React errors
- ReactFlow has debug mode in docs

### Port Conflicts

```bash
./run.sh status    # Shows what's using ports
./run.sh stop      # Stops everything
```

## Resources

### Documentation

- [OVN Tracing Docs](https://docs.redhat.com/en/documentation/openshift_container_platform/4.22/html/ovn-kubernetes_network_plugin/ovn-kubernetes-tracing-using-ovntrace)
- [ReactFlow Docs](https://reactflow.dev/)
- [Dagre Layout](https://github.com/dagrejs/dagre)

### Internal Docs

- `docs/implementation-plan.md` - Complete 4-phase roadmap
- `GETTING_STARTED.md` - Setup and running guide
- `PHASE1_STATUS.md` - Current progress tracking
- `RUN_SCRIPT_GUIDE.md` - run.sh command reference

## Questions to Ask the User

When the user asks you to work on this project:

1. **"Do you have access to a real OVN-Kubernetes cluster?"** - Helps determine if you can get real ovn-trace samples
2. **"What specific scenario are you working on?"** - Pod-to-pod, service, ingress, etc.
3. **"Are you focused on parser, visualization, or both?"** - Helps prioritize
4. **"Do you want tests added?"** - Should be yes, but ask
5. **"Should this work offline or need real OVN?"** - Mock vs real executor

## Working with This Project

### First Time as Agent

1. Read this file (you're doing it!)
2. Check `PHASE1_STATUS.md` for current state
3. Run `./run.sh status` to see what's running
4. Ask user what they want to work on
5. Check relevant files before making changes

### Continuing Work

1. Check `PHASE1_STATUS.md` for what's done/in-progress
2. Review recent changes if user provides context
3. Ensure types stay in sync (backend ↔ frontend)
4. Add tests for new functionality
5. Update `PHASE1_STATUS.md` when completing milestones

### Making Changes

1. **Understand impact** - Data model changes affect multiple files
2. **Keep types in sync** - Go structs ↔ TypeScript interfaces
3. **Test locally** - Use `./run.sh dev` for quick iteration
4. **Document decisions** - Update this file or PHASE1_STATUS.md
5. **Think phases** - Will this work in Phase 2, 3, 4?

## Quick Reference

### File Purposes

```
backend/
├── cmd/server/main.go           # HTTP server setup
├── pkg/models/types.go          # Data structures (SOURCE OF TRUTH)
├── pkg/parser/parser.go         # ovn-trace parser (NEEDS WORK)
├── pkg/executor/trace.go        # Run ovn-trace (currently mocked)
└── pkg/api/handlers/trace.go    # API request handler

frontend/
├── src/App.tsx                  # Main React component
├── src/types/index.ts           # Types (MUST MATCH backend)
├── src/components/TraceForm/    # User input form
└── src/components/FlowVisualization/  # Graph rendering

deploy/
└── standalone/docker-compose.yaml    # Local deployment

examples/
└── sample-traces/               # Real ovn-trace outputs (USE THESE)

docs/
└── implementation-plan.md       # Full roadmap (READ THIS)
```

### Quick Commands

```bash
./run.sh help      # Show all commands
./run.sh dev       # Start development
./run.sh status    # Check what's running
./run.sh stop      # Stop everything
make backend-test  # Run Go tests (when they exist)
```

### Data Flow Example

```
User fills form → POST /api/v1/trace
                       ↓
                  trace.HandleTrace()
                       ↓
                  executor.Execute()  (returns mock ovn-trace output)
                       ↓
                  parser.Parse()      (converts to TraceResult)
                       ↓
                  JSON response
                       ↓
                  Frontend receives TraceResult
                       ↓
                  ReactFlow renders graph
```

## Success Criteria for Phase 1

When you're working on this project, Phase 1 is complete when:

- [ ] Parser handles 95%+ of common ovn-trace scenarios
- [ ] Graph renders in <2s for typical traces
- [ ] User can execute trace without reading docs
- [ ] Docker Compose deployment works smoothly
- [ ] Test coverage >70%
- [ ] Basic documentation is complete

Currently at ~30% overall.

## Final Notes

This is an **active development project**. The foundation is solid, but there's significant work ahead on the parser and graph builder. The mock data lets you develop the full stack, but eventually this needs to integrate with real OVN environments.

The project will evolve through 4 phases over ~6 months. Phase 1 is about proving the concept works. Phase 2-4 add Kubernetes/OpenShift integration.

**Most Important**: Keep backend and frontend types in sync. They're the contract between the two systems.

---

**Last Updated**: October 26, 2026  
**Phase**: 1 (Standalone MVP)  
**Progress**: 30%  
**Next Milestone**: Enhanced Parser (Milestone 1.1)

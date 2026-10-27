# Getting Started with OVN-Kubernetes Trace Viewer

This guide will help you get the Phase 1 standalone application up and running.

## Prerequisites

- **Go** 1.21 or later
- **Node.js** 24.22.1 or later (LTS)
- **npm** 10.0.0 or later
- **Docker** and **Docker Compose** (optional, for containerized deployment)
- **Make** (optional, but recommended)

## Quick Start

### Option 1: Using run.sh Script (Easiest)

The `run.sh` script provides the easiest way to get started:

1. **Initialize the project**:
   ```bash
   ./run.sh init
   ```

2. **Run with Docker** (recommended):
   ```bash
   ./run.sh docker
   ```
   
   Or **run in development mode**:
   ```bash
   ./run.sh dev
   ```

3. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

**Available commands**:
```bash
./run.sh help      # Show all available commands
./run.sh status    # Check if services are running
./run.sh logs      # View Docker logs
./run.sh stop      # Stop all services
./run.sh clean     # Clean build artifacts
```

### Option 2: Using Make

1. **Initialize the project** (install dependencies):
   ```bash
   make init
   ```

2. **Run in development mode**:
   
   In one terminal, start the backend:
   ```bash
   make backend-run
   ```

   In another terminal, start the frontend:
   ```bash
   make frontend-dev
   ```

3. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

### Option 3: Using Docker Compose

1. **Build and start services**:
   ```bash
   make docker-build
   make docker-up
   ```

2. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

3. **View logs**:
   ```bash
   make docker-logs
   ```

4. **Stop services**:
   ```bash
   make docker-down
   ```

### Option 4: Manual Setup

#### Backend

```bash
cd backend
go mod tidy
go run ./cmd/server
```

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Testing the Application

1. Open http://localhost:3000 in your browser

2. Fill in the trace form:
   - **Source**: `pod/web-pod` or `10.128.0.5`
   - **Destination**: `pod/db-pod` or `10.128.0.10`
   - **Protocol**: `tcp`
   - **Dest Port**: `5432`

3. Click "Execute Trace"

4. View the generated flow visualization and step-by-step trace

## Current Status

### ✅ Implemented (Phase 1 - MVP Foundation)

- Backend Go API server with REST endpoints
- Frontend React app with TypeScript
- Basic trace form for input
- Flow visualization with ReactFlow
- Mock trace executor (for development without OVN access)
- Basic parser (proof-of-concept)
- Docker containerization
- Docker Compose deployment

### 🚧 In Progress / TODO

- **Parser Enhancement**: Robust ovn-trace output parsing
  - Handle all OVN component types (switches, routers, ACLs)
  - Extract packet transformations (NAT, TTL, etc.)
  - Detect dropped packets and error conditions
  
- **Graph Builder**: Intelligent topology layout
  - Automatic node positioning (Dagre algorithm)
  - Edge routing and labeling
  - Collapsible node groups
  
- **OVN Integration**: Real ovn-trace execution
  - Detect OVN environment
  - Build proper ovn-trace commands
  - Handle different datapaths
  
- **Testing**: Comprehensive test coverage
  - Unit tests for parser
  - Integration tests for API
  - Frontend component tests
  
- **Documentation**:
  - API documentation
  - User guide
  - Sample traces collection

## Development Workflow

### Backend Development

```bash
# Run tests
make backend-test

# Build binary
make backend-build

# Run the server
make backend-run
```

### Frontend Development

```bash
# Install dependencies
cd frontend && npm install

# Run dev server (with hot reload)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Adding Sample Traces

Add real ovn-trace outputs to `examples/sample-traces/` for parser testing:

```bash
# From an OVN-Kubernetes environment
ovn-trace <datapath> <packet-spec> > examples/sample-traces/my-scenario.txt
```

## Project Structure

```
ovnkube-viewer/
├── backend/              # Go API server
│   ├── cmd/server/      # Main entry point
│   ├── pkg/
│   │   ├── models/      # Data models
│   │   ├── parser/      # ovn-trace parser
│   │   ├── executor/    # Trace executor
│   │   └── api/         # HTTP handlers & middleware
│   └── Dockerfile
├── frontend/            # React app
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── types/       # TypeScript types
│   │   └── App.tsx     # Main app
│   └── Dockerfile
├── deploy/
│   └── standalone/
│       └── docker-compose.yaml
└── examples/
    └── sample-traces/   # Sample ovn-trace outputs
```

## Next Steps

1. **Collect real ovn-trace samples**: From a working OVN-Kubernetes cluster
2. **Enhance parser**: Build comprehensive parsing logic
3. **Improve visualization**: Better graph layouts, custom node types
4. **Add OVN integration**: Connect to real OVN environment
5. **Testing**: Add comprehensive test coverage

## Troubleshooting

### Backend won't start
- Check if port 8080 is already in use
- Verify Go is installed: `go version`
- Run `go mod tidy` to install dependencies

### Frontend won't start
- Check if port 3000 is already in use
- Verify Node.js is installed: `node --version`
- Delete `node_modules` and run `npm install` again

### Docker issues
- Ensure Docker daemon is running
- Try rebuilding: `make docker-down && make docker-build && make docker-up`

## Resources

- [Implementation Plan](implementation-plan.md)
- [OVN-Kubernetes Tracing Docs](https://docs.redhat.com/en/documentation/openshift_container_platform/4.22/html/ovn-kubernetes_network_plugin/ovn-kubernetes-tracing-using-ovntrace)
- [ReactFlow Docs](https://reactflow.dev/)

## Contributing

See the main [README.md](../README.md) for contribution guidelines.

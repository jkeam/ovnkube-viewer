# OVN-Kubernetes Trace Viewer

A visual viewer for OVN-Kubernetes traffic tracing that transforms text-based ovn-trace output into interactive flow diagrams.

## Overview

This project provides a graphical interface to visualize network traffic flow through OVN-Kubernetes components, making it easier to:
- Debug network connectivity issues
- Understand packet flow through logical switches and routers
- Visualize NetworkPolicy enforcement
- Trace traffic between pods, services, and external endpoints

## Project Status

🚀 **Phase 1 MVP - In Development** - Basic standalone application is scaffolded and ready for development!

See [GETTING_STARTED.md](GETTING_STARTED.md) for setup instructions and [Implementation Plan](docs/implementation-plan.md) for detailed roadmap.

### Current Implementation Status

- ✅ Project structure created
- ✅ Backend API server (Go) with basic REST endpoints
- ✅ Frontend React app with TypeScript, ReactFlow, and PatternFly
- ✅ Red Hat UX design system implementation
- ✅ Mock trace executor for development
- ✅ Basic parser foundation
- ✅ Docker containerization
- ✅ Docker Compose setup for easy deployment
- 🚧 Parser enhancement (robust ovn-trace parsing)
- 🚧 Graph layout optimization
- 🚧 Real OVN integration
- 🚧 Comprehensive testing

## Roadmap

### Phase 1: Standalone Application (MVP)
- Parse ovn-trace output
- Interactive graph visualization
- REST API backend (Go)
- React frontend with ReactFlow
- Docker deployment

### Phase 2: Kubernetes Native
- Helm chart deployment
- Kubernetes RBAC integration
- Persistent trace history
- Multi-tenancy support

### Phase 3: OpenShift Operator
- Custom Resource Definitions (TraceRequest, TraceResult)
- Operator lifecycle management
- Scheduled tracing
- Topology caching

### Phase 4: OpenShift Console Plugin
- Integrated into OpenShift web console
- PatternFly UI components
- One-click tracing from resource pages
- Native console authentication

## Architecture

```
┌──────────────┐          ┌──────────────────────────────┐
│   Web UI     │          │      Backend API Server      │
│  (React)     │◄────────►│         (Go)                 │
│              │          │                              │
│ - Input form │          │ - Parse ovn-trace output     │
│ - Flow viz   │          │ - Build topology graph       │
│   (ReactFlow)│          │ - Execute trace commands     │
└──────────────┘          └──────────────┬───────────────┘
                                          │
                                          ▼
                          ┌──────────────────────────────┐
                          │   OVN/OVS Components         │
                          │   (via ovn-trace CLI)        │
                          └──────────────────────────────┘
```

## Technology Stack

- **Backend**: Go 1.21+
- **Frontend**: React + TypeScript + ReactFlow + PatternFly (Node.js 24.22.1+ LTS)
- **Design System**: PatternFly 5.4 + Red Hat UX
- **Deployment**: Docker, Kubernetes, Helm, OpenShift Operator
- **Operator**: Operator SDK

## Quick Start

**Easiest way** - Use the run script:

```bash
# Initialize and run with Docker (recommended)
./run.sh init
./run.sh docker

# Or run in development mode
./run.sh dev

# View help for all options
./run.sh help
```

**Alternative** - Use Make commands:

```bash
# Initialize dependencies
make init

# Start with Docker Compose
make docker-build
make docker-up

# Or run in development mode (2 terminals)
make backend-run     # Terminal 1
make frontend-dev    # Terminal 2
```

Then open http://localhost:3000 in your browser.

See [GETTING_STARTED.md](GETTING_STARTED.md) for detailed instructions.

## Documentation

- [Getting Started Guide](GETTING_STARTED.md) - Setup and development instructions
- [Run Script Guide](RUN_SCRIPT_GUIDE.md) - run.sh command reference
- [Node.js Upgrade Guide](NODE_UPGRADE.md) - How to upgrade to Node.js 24 LTS
- [PatternFly Design](PATTERNFLY_DESIGN.md) - Design system implementation
- [Phase 1 Status](PHASE1_STATUS.md) - Current progress and next steps
- [Implementation Plan](docs/implementation-plan.md) - Comprehensive development roadmap
- [AGENT.md](AGENT.md) - Context for AI agents working on this project
- [Architecture](docs/architecture.md) - Technical architecture details (coming soon)
- [User Guide](docs/user-guide.md) - How to use the viewer (coming soon)
- [Development Guide](docs/development.md) - Contributing and development setup (coming soon)

## Reference

This project is based on the OVN-Kubernetes tracing functionality documented at:
https://docs.redhat.com/en/documentation/openshift_container_platform/4.22/html/ovn-kubernetes_network_plugin/ovn-kubernetes-tracing-using-ovntrace

## License

TBD

## Contributing

TBD

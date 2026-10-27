# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Changed

- **Go version requirement updated to 1.26.10** - Matches the Go version in UBI9 go-toolset:9.8
  - Updated `backend/go.mod` from 1.27.3 to 1.26.10
  - Go 1.26.10 is Red Hat's version (Red Hat 1.26.10-1.el9_8)

- **Node.js version requirement updated to 24.21.0** - Matches the Node.js version in UBI9 nodejs-24:9.8
  - Updated `frontend/package.json` engines from >=24.22.1 to >=24.21.0
  - Node.js 24.21.0 is Red Hat's version in nodejs-24 container
  - npm version: 11.19.0

- **All containers rebased to Red Hat UBI9** - All Dockerfiles now use Red Hat Universal Base Images:
  
  **Backend:**
  - Build stage: `registry.redhat.io/ubi9/go-toolset:9.8-1793107762`
  - Runtime stage: `registry.redhat.io/ubi9/ubi-micro:9.8-1791895366`
  
  **Frontend & Plugin:**
  - Build stage: `registry.redhat.io/ubi9/nodejs-24:9.8-1793067767`
  - Runtime stage: `registry.redhat.io/ubi9/nginx-124:1-113`
  
  **Benefits:**
  - Enterprise support and security updates
  - OpenShift optimized
  - All containers run as non-root (UID 1001)
  - Smaller attack surface
  
  **Breaking Changes:**
  - Frontend now listens on port 8080 (was 80)
  - Service manifests updated accordingly
  - nginx configuration paths changed to `/opt/app-root/src`
  
  See [docs/ubi-base-images.md](ubi-base-images.md) for details

### Added

- **OpenShift Console Plugin** - Full dynamic plugin integration
  - Appears in OpenShift Console "Observe" menu
  - Native console authentication
  - Plugin deployment manifests and documentation
  - See [deploy/openshift/CONSOLE-PLUGIN-GUIDE.md](../deploy/openshift/CONSOLE-PLUGIN-GUIDE.md)

- **Kustomize Deployment** - OpenShift deployment with Kustomize
  - Base manifests with dev and prod overlays
  - ServiceAccount with elevated OVN permissions
  - ClusterRole for pod exec and network resources
  - See [deploy/openshift/DEPLOYMENT-GUIDE.md](../deploy/openshift/DEPLOYMENT-GUIDE.md)

- **Documentation**
  - [Console Plugin Guide](../deploy/openshift/CONSOLE-PLUGIN-GUIDE.md)
  - [Console Plugin Quick Start](../deploy/openshift/QUICKSTART-CONSOLE-PLUGIN.md)
  - [Deployment Guide](../deploy/openshift/DEPLOYMENT-GUIDE.md)
  - [Building Containers Guide](building-containers.md)
  - [UBI Base Images Guide](ubi-base-images.md)
  - [Build Quick Reference](BUILD-QUICK-REFERENCE.md)

- **Scripts**
  - `deploy/openshift/deploy.sh` - Deployment helper script
  - `deploy/openshift/enable-plugin.sh` - Enable/disable console plugin

### Security

- Backend container now runs as non-root (UID 1001)
- Static Go binary with no dynamic dependencies
- Minimal runtime image (ubi-micro, ~40MB)
- No shell in runtime container

## [0.1.0] - Initial Release

### Added

- Backend API server in Go
  - REST endpoints for trace operations
  - OVN trace execution via pod exec
  - Mock trace executor for development

- Frontend React application
  - PatternFly UI components
  - Red Hat UX design system
  - Basic trace input form
  - Flow visualization placeholder

- Docker deployment
  - Docker Compose configuration
  - Multi-stage Dockerfiles
  - Local development setup

- Documentation
  - README with project overview
  - Getting started guide
  - Implementation plan
  - PatternFly design system docs

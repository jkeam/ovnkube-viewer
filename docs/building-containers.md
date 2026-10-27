# Building Containers Guide

Complete guide for building all container images for the OVN Trace Viewer.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Prerequisites](#prerequisites)
3. [Using Make](#using-make)
4. [Manual Build with Podman](#manual-build-with-podman)
5. [Manual Build with Docker](#manual-build-with-docker)
6. [Build for Different Registries](#build-for-different-registries)
7. [Troubleshooting](#troubleshooting)

## Quick Start

```bash
# Set your registry
export REGISTRY=quay.io/yourorg

# Build all three images
make openshift-build REGISTRY=$REGISTRY TAG=latest

# Push to registry
make openshift-push REGISTRY=$REGISTRY TAG=latest
```

## Prerequisites

### Required Tools

Choose one container builder:
- **Podman** (recommended for OpenShift/RHEL) - `podman version`
- **Docker** - `docker version`

### Check Installed Tools

```bash
# Check podman
podman version

# OR check docker
docker version
```

### Login to Registry

Before pushing images, login to your container registry:

```bash
# For Quay.io
podman login quay.io
# Enter username and password

# For Docker Hub
podman login docker.io

# For OpenShift internal registry
podman login $(oc get route default-route -n openshift-image-registry -o jsonpath='{.spec.host}')
```

## Using Make

### 1. Build All Images

```bash
# Build with default settings (quay.io/yourorg)
make openshift-build

# Build with custom registry and tag
make openshift-build REGISTRY=quay.io/mycompany TAG=v1.0.0

# Build with different tag
make openshift-build TAG=dev
```

This builds three images:
1. Backend Go service
2. Frontend React app (standalone with nginx)
3. Console plugin (dynamic plugin for OpenShift Console)

### 2. Push All Images

```bash
# Push to registry
make openshift-push REGISTRY=quay.io/mycompany TAG=v1.0.0
```

### 3. Build and Push Together

```bash
# One command to build and push
make openshift-build openshift-push REGISTRY=quay.io/mycompany TAG=v1.0.0
```

### 4. All-in-One: Build, Push, Deploy

```bash
# Complete workflow
make openshift-all OVERLAY=dev REGISTRY=quay.io/mycompany TAG=dev
```

## Manual Build with Podman

### Build Backend

```bash
# Navigate to project root
cd /path/to/ovnkube-viewer

# Build backend image
podman build \
  -t quay.io/yourorg/ovnkube-viewer-backend:latest \
  -f backend/Dockerfile \
  backend/

# Verify image
podman images | grep ovnkube-viewer-backend
```

### Build Frontend (Standalone)

```bash
# Build frontend image
podman build \
  -t quay.io/yourorg/ovnkube-viewer-frontend:latest \
  -f frontend/Dockerfile \
  frontend/

# Verify image
podman images | grep ovnkube-viewer-frontend
```

### Build Console Plugin

```bash
# Build plugin image
podman build \
  -t quay.io/yourorg/ovnkube-viewer-plugin:latest \
  -f frontend/Dockerfile.plugin \
  frontend/

# Verify image
podman images | grep ovnkube-viewer-plugin
```

### Push Images

```bash
# Login first
podman login quay.io

# Push backend
podman push quay.io/yourorg/ovnkube-viewer-backend:latest

# Push frontend
podman push quay.io/yourorg/ovnkube-viewer-frontend:latest

# Push plugin
podman push quay.io/yourorg/ovnkube-viewer-plugin:latest
```

## Manual Build with Docker

Same as podman, just replace `podman` with `docker`:

```bash
# Build backend
docker build \
  -t quay.io/yourorg/ovnkube-viewer-backend:latest \
  -f backend/Dockerfile \
  backend/

# Build frontend
docker build \
  -t quay.io/yourorg/ovnkube-viewer-frontend:latest \
  -f frontend/Dockerfile \
  frontend/

# Build plugin
docker build \
  -t quay.io/yourorg/ovnkube-viewer-plugin:latest \
  -f frontend/Dockerfile.plugin \
  frontend/

# Push all
docker push quay.io/yourorg/ovnkube-viewer-backend:latest
docker push quay.io/yourorg/ovnkube-viewer-frontend:latest
docker push quay.io/yourorg/ovnkube-viewer-plugin:latest
```

## Build for Different Registries

### Quay.io (Default)

```bash
export REGISTRY=quay.io/yourorg
make openshift-build openshift-push REGISTRY=$REGISTRY TAG=latest
```

### Docker Hub

```bash
export REGISTRY=docker.io/yourusername
make openshift-build openshift-push REGISTRY=$REGISTRY TAG=latest
```

### OpenShift Internal Registry

```bash
# Get internal registry URL
INTERNAL_REGISTRY=$(oc get route default-route -n openshift-image-registry -o jsonpath='{.spec.host}')

# Login
podman login -u $(oc whoami) -p $(oc whoami -t) $INTERNAL_REGISTRY

# Build and push
export REGISTRY=$INTERNAL_REGISTRY/ovnkube-viewer
make openshift-build openshift-push REGISTRY=$REGISTRY TAG=latest
```

### Private Harbor Registry

```bash
export REGISTRY=harbor.example.com/ovnkube-viewer
podman login harbor.example.com
make openshift-build openshift-push REGISTRY=$REGISTRY TAG=latest
```

### GitHub Container Registry (ghcr.io)

```bash
export REGISTRY=ghcr.io/yourusername
echo $GITHUB_TOKEN | podman login ghcr.io -u yourusername --password-stdin
make openshift-build openshift-push REGISTRY=$REGISTRY TAG=latest
```

## Build with Different Tags

### Development Build

```bash
make openshift-build REGISTRY=quay.io/yourorg TAG=dev
```

### Semantic Versioning

```bash
# Build specific version
make openshift-build REGISTRY=quay.io/yourorg TAG=v1.2.3

# Also tag as latest
podman tag quay.io/yourorg/ovnkube-viewer-backend:v1.2.3 quay.io/yourorg/ovnkube-viewer-backend:latest
podman tag quay.io/yourorg/ovnkube-viewer-frontend:v1.2.3 quay.io/yourorg/ovnkube-viewer-frontend:latest
podman tag quay.io/yourorg/ovnkube-viewer-plugin:v1.2.3 quay.io/yourorg/ovnkube-viewer-plugin:latest

# Push both
make openshift-push TAG=v1.2.3
make openshift-push TAG=latest
```

### Git Commit SHA

```bash
# Use git commit as tag
GIT_SHA=$(git rev-parse --short HEAD)
make openshift-build REGISTRY=quay.io/yourorg TAG=$GIT_SHA
```

### Branch Name

```bash
# Use current branch as tag
BRANCH=$(git branch --show-current)
make openshift-build REGISTRY=quay.io/yourorg TAG=$BRANCH
```

## Multi-Architecture Builds

Build for multiple architectures (amd64, arm64):

### Using Podman Manifest

```bash
# Create manifest
podman manifest create quay.io/yourorg/ovnkube-viewer-backend:latest

# Build for amd64
podman build \
  --platform linux/amd64 \
  --manifest quay.io/yourorg/ovnkube-viewer-backend:latest \
  -f backend/Dockerfile \
  backend/

# Build for arm64
podman build \
  --platform linux/arm64 \
  --manifest quay.io/yourorg/ovnkube-viewer-backend:latest \
  -f backend/Dockerfile \
  backend/

# Push manifest
podman manifest push quay.io/yourorg/ovnkube-viewer-backend:latest
```

### Using Docker Buildx

```bash
# Create builder
docker buildx create --name multiarch --use

# Build and push for multiple platforms
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t quay.io/yourorg/ovnkube-viewer-backend:latest \
  -f backend/Dockerfile \
  --push \
  backend/
```

## Verify Images

### List Local Images

```bash
# Podman
podman images | grep ovnkube-viewer

# Docker
docker images | grep ovnkube-viewer
```

### Inspect Image

```bash
# Podman
podman inspect quay.io/yourorg/ovnkube-viewer-backend:latest

# Check image size
podman images quay.io/yourorg/ovnkube-viewer-backend:latest --format "{{.Size}}"

# Check layers
podman history quay.io/yourorg/ovnkube-viewer-backend:latest
```

### Test Image Locally

```bash
# Run backend container
podman run -p 8080:8080 quay.io/yourorg/ovnkube-viewer-backend:latest

# Run frontend container
podman run -p 3000:80 quay.io/yourorg/ovnkube-viewer-frontend:latest

# Test health
curl http://localhost:8080/health
curl http://localhost:3000
```

## Optimize Build Times

### Use Build Cache

Podman and Docker automatically cache layers. To leverage this:

1. Don't change Dockerfiles unnecessarily
2. Order instructions from least to most frequently changing
3. Copy dependency files before source code

### Parallel Builds

Build images in parallel (requires modern shell):

```bash
# Using background jobs
podman build -t quay.io/yourorg/ovnkube-viewer-backend:latest -f backend/Dockerfile backend/ &
podman build -t quay.io/yourorg/ovnkube-viewer-frontend:latest -f frontend/Dockerfile frontend/ &
podman build -t quay.io/yourorg/ovnkube-viewer-plugin:latest -f frontend/Dockerfile.plugin frontend/ &
wait
```

### Use --squash

Reduce image layers:

```bash
podman build --squash \
  -t quay.io/yourorg/ovnkube-viewer-backend:latest \
  -f backend/Dockerfile \
  backend/
```

## Clean Up

### Remove Built Images

```bash
# Remove specific images
podman rmi quay.io/yourorg/ovnkube-viewer-backend:latest
podman rmi quay.io/yourorg/ovnkube-viewer-frontend:latest
podman rmi quay.io/yourorg/ovnkube-viewer-plugin:latest

# Remove all ovnkube-viewer images
podman rmi $(podman images | grep ovnkube-viewer | awk '{print $3}')

# Remove dangling images
podman image prune -f
```

### Remove Build Cache

```bash
# Podman
podman system prune -af

# Docker
docker system prune -af
```

## Troubleshooting

### Build Fails: "Cannot find module"

**Frontend builds:**

```bash
# Install dependencies first
cd frontend
npm install
cd ..

# Then build
make openshift-build
```

### Build Fails: "go.mod not found"

**Backend builds:**

```bash
# Initialize Go modules
cd backend
go mod tidy
cd ..

# Then build
make openshift-build
```

### Build Fails: "Permission denied"

```bash
# On SELinux systems (RHEL, Fedora), use :Z flag
podman build --security-opt label=disable \
  -t quay.io/yourorg/ovnkube-viewer-backend:latest \
  -f backend/Dockerfile \
  backend/
```

### Push Fails: "Unauthorized"

```bash
# Login to registry
podman login quay.io

# Verify credentials
podman login --get-login quay.io
```

### Image Too Large

**Check what's taking space:**

```bash
podman history quay.io/yourorg/ovnkube-viewer-backend:latest
```

**Optimize Dockerfile:**
- Use multi-stage builds (already done)
- Remove unnecessary files
- Use .dockerignore
- Don't install dev dependencies in production

### Build is Slow

**Use build cache:**

```bash
# Podman preserves cache automatically
# To disable cache for fresh build:
podman build --no-cache \
  -t quay.io/yourorg/ovnkube-viewer-backend:latest \
  -f backend/Dockerfile \
  backend/
```

### Can't Connect to Podman Socket

```bash
# Start podman socket
systemctl --user start podman.socket

# Or use sudo
sudo systemctl start podman.socket
```

### Different Builder (Podman vs Docker)

If Makefile uses podman but you have docker:

```bash
# Option 1: Install podman
sudo dnf install podman  # Fedora/RHEL
sudo apt install podman  # Ubuntu

# Option 2: Create alias
alias podman=docker

# Option 3: Edit Makefile to use docker
# Change 'podman' to 'docker' in Makefile
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Build Images

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Login to Quay
        uses: docker/login-action@v2
        with:
          registry: quay.io
          username: ${{ secrets.QUAY_USERNAME }}
          password: ${{ secrets.QUAY_PASSWORD }}
      
      - name: Build and Push
        env:
          REGISTRY: quay.io/${{ secrets.QUAY_ORG }}
          TAG: ${{ github.sha }}
        run: |
          make openshift-build openshift-push REGISTRY=$REGISTRY TAG=$TAG
```

### GitLab CI Example

```yaml
build:
  stage: build
  image: quay.io/podman/stable
  script:
    - podman login -u $QUAY_USER -p $QUAY_PASSWORD quay.io
    - make openshift-build openshift-push REGISTRY=quay.io/yourorg TAG=$CI_COMMIT_SHA
```

## Next Steps

After building images:

1. **Update manifests** with your registry:
   ```bash
   # Edit deploy/openshift/base/kustomization.yaml
   # Update image names with your registry
   ```

2. **Push images** to registry:
   ```bash
   make openshift-push REGISTRY=quay.io/yourorg TAG=latest
   ```

3. **Deploy** to OpenShift:
   ```bash
   make openshift-deploy OVERLAY=dev
   ```

4. **Verify** deployment:
   ```bash
   oc get pods -n ovnkube-viewer-dev
   ```

## Summary of Commands

```bash
# Complete workflow
export REGISTRY=quay.io/yourorg
export TAG=v1.0.0

# 1. Build all images
make openshift-build REGISTRY=$REGISTRY TAG=$TAG

# 2. Push to registry
make openshift-push REGISTRY=$REGISTRY TAG=$TAG

# 3. Deploy to OpenShift
make openshift-deploy OVERLAY=dev

# 4. Enable console plugin
./deploy/openshift/enable-plugin.sh enable
```

## References

- [Podman Documentation](https://docs.podman.io/)
- [Docker Documentation](https://docs.docker.com/)
- [OpenShift Container Images](https://docs.openshift.com/container-platform/latest/openshift_images/index.html)
- [Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)

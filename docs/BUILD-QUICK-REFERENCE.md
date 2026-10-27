# Build Quick Reference

Quick commands to build all containers.

## 🚀 Fastest Way (Using Make)

```bash
# 1. Set your registry
export REGISTRY=quay.io/yourorg

# 2. Build all three images (backend, frontend, plugin)
make openshift-build REGISTRY=$REGISTRY TAG=latest

# 3. Push to registry
make openshift-push REGISTRY=$REGISTRY TAG=latest
```

**Or in one line:**
```bash
make openshift-build openshift-push REGISTRY=quay.io/yourorg TAG=latest
```

## 📦 What Gets Built

Three container images:

1. **Backend** - `quay.io/yourorg/ovnkube-viewer-backend:latest`
   - Go API server
   - Executes OVN traces
   - Based on Red Hat UBI9 (Universal Base Image)

2. **Frontend** - `quay.io/yourorg/ovnkube-viewer-frontend:latest`
   - React app + Nginx
   - Standalone mode

3. **Plugin** - `quay.io/yourorg/ovnkube-viewer-plugin:latest`
   - Console plugin bundle
   - OpenShift Console integration

**Note:** Backend uses Red Hat UBI9 base images. See [docs/ubi-base-images.md](docs/ubi-base-images.md) for details.

## 🔑 Prerequisites

**Login to your image registry:**
```bash
# For pushing images
podman login quay.io
# Enter username and password
```

**Optional - Login to Red Hat registry** (for better rate limits on UBI base images):
```bash
podman login registry.redhat.io
# Not required - UBI images are freely available
# But login provides better rate limits
```

## 🏷️ Build with Different Tags

```bash
# Development
make openshift-build REGISTRY=quay.io/yourorg TAG=dev

# Production version
make openshift-build REGISTRY=quay.io/yourorg TAG=v1.0.0

# Git commit SHA
make openshift-build REGISTRY=quay.io/yourorg TAG=$(git rev-parse --short HEAD)
```

## 🔨 Manual Build (Podman)

If you prefer to build manually:

```bash
# Backend
podman build -t quay.io/yourorg/ovnkube-viewer-backend:latest \
  -f backend/Dockerfile backend/

# Frontend
podman build -t quay.io/yourorg/ovnkube-viewer-frontend:latest \
  -f frontend/Dockerfile frontend/

# Plugin
podman build -t quay.io/yourorg/ovnkube-viewer-plugin:latest \
  -f frontend/Dockerfile.plugin frontend/

# Push all
podman push quay.io/yourorg/ovnkube-viewer-backend:latest
podman push quay.io/yourorg/ovnkube-viewer-frontend:latest
podman push quay.io/yourorg/ovnkube-viewer-plugin:latest
```

## 🐳 Using Docker Instead

Replace `podman` with `docker`:

```bash
make openshift-build REGISTRY=quay.io/yourorg TAG=latest
# Uses podman by default

# Or manually with docker
docker build -t quay.io/yourorg/ovnkube-viewer-backend:latest -f backend/Dockerfile backend/
docker build -t quay.io/yourorg/ovnkube-viewer-frontend:latest -f frontend/Dockerfile frontend/
docker build -t quay.io/yourorg/ovnkube-viewer-plugin:latest -f frontend/Dockerfile.plugin frontend/
```

## 🔍 Verify Images

```bash
# List images
podman images | grep ovnkube-viewer

# Check image size
podman images quay.io/yourorg/ovnkube-viewer-backend:latest

# Test backend locally
podman run -p 8080:8080 quay.io/yourorg/ovnkube-viewer-backend:latest
curl http://localhost:8080/health
```

## 📝 Update Image References

After building, update your Kustomize configuration:

**Edit:** `deploy/openshift/base/kustomization.yaml`

```yaml
images:
  - name: ovnkube-viewer-backend
    newName: quay.io/yourorg/ovnkube-viewer-backend  # ← Update this
    newTag: latest
  - name: ovnkube-viewer-frontend
    newName: quay.io/yourorg/ovnkube-viewer-frontend  # ← Update this
    newTag: latest
  - name: ovnkube-viewer-plugin
    newName: quay.io/yourorg/ovnkube-viewer-plugin  # ← Update this
    newTag: latest
```

## 🚢 Complete Workflow

```bash
# 1. Build
export REGISTRY=quay.io/yourorg
make openshift-build REGISTRY=$REGISTRY TAG=dev

# 2. Push
make openshift-push REGISTRY=$REGISTRY TAG=dev

# 3. Deploy
make openshift-deploy OVERLAY=dev

# 4. Enable plugin
./deploy/openshift/enable-plugin.sh enable
```

## ⚡ One-Liner for Impatient People

```bash
export REGISTRY=quay.io/yourorg && \
  podman login quay.io && \
  make openshift-build openshift-push REGISTRY=$REGISTRY TAG=latest && \
  make openshift-deploy OVERLAY=dev && \
  ./deploy/openshift/enable-plugin.sh enable
```

## 🔧 Troubleshooting

**Build fails - npm dependencies:**
```bash
cd frontend && npm install && cd ..
make openshift-build
```

**Build fails - go modules:**
```bash
cd backend && go mod tidy && cd ..
make openshift-build
```

**Push fails - not logged in:**
```bash
podman login quay.io
make openshift-push
```

**Wrong registry in manifests:**
```bash
# Edit deploy/openshift/base/kustomization.yaml
# Change quay.io/yourorg to your actual registry
```

## 📚 More Details

See [docs/building-containers.md](docs/building-containers.md) for:
- Different registries (Docker Hub, Harbor, ghcr.io)
- Multi-architecture builds
- CI/CD integration
- Advanced troubleshooting

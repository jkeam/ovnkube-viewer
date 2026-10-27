# UBI Migration Summary

All containers have been migrated from Alpine Linux to Red Hat Universal Base Images (UBI9).

## 🎯 What Changed

### Backend Container

**Before:**
```dockerfile
FROM golang:1.21-alpine AS builder
FROM alpine:latest
```

**After:**
```dockerfile
FROM registry.redhat.io/ubi9/go-toolset:9.8-1793107762 AS builder
FROM registry.redhat.io/ubi9/ubi-micro:9.8-1791895366
```

### Frontend Container

**Before:**
```dockerfile
FROM node:24-alpine AS builder
FROM nginx:alpine
```

**After:**
```dockerfile
FROM registry.redhat.io/ubi9/nodejs-24:9.8-1793067767 AS builder
FROM registry.redhat.io/ubi9/nginx-126:9.8-1793065662
```

### Console Plugin Container

**Before:**
```dockerfile
FROM node:24-alpine AS builder
FROM nginx:alpine
```

**After:**
```dockerfile
FROM registry.redhat.io/ubi9/nodejs-24:9.8-1793067767 AS builder
FROM registry.redhat.io/ubi9/nginx-126:9.8-1793065662
```

## 🔄 Breaking Changes

### Port Changes

| Component | Old Port | New Port | Reason |
|-----------|----------|----------|--------|
| Frontend | 80 | 8080 | UBI nginx runs on non-privileged port |
| Plugin | 9443 | 9443 | No change |
| Backend | 8080 | 8080 | No change |

### Updated Files

The following files were automatically updated for the port change:

1. **`deploy/openshift/base/frontend-deployment.yaml`**
   - Container port: 80 → 8080
   - Removed hardcoded user ID (now uses UBI default 1001)

2. **`deploy/openshift/base/frontend-service.yaml`**
   - Service port: 80 → 8080

3. **`deploy/openshift/base/plugin-deployment.yaml`**
   - Removed hardcoded user ID

4. **`deploy/standalone/docker-compose.yaml`**
   - Port mapping: `3000:80` → `3000:8080`

5. **`frontend/nginx.conf`**
   - Listen port: 80 → 8080
   - Document root: `/usr/share/nginx/html` → `/opt/app-root/src`

6. **`frontend/nginx.plugin.conf`**
   - Document root: `/usr/share/nginx/html` → `/opt/app-root/src`

## ✅ No Action Required

If you're deploying fresh, everything just works. The manifests have been updated.

## 🔧 Migration Steps (For Existing Deployments)

If you have an existing deployment:

### 1. Rebuild Images

```bash
export REGISTRY=quay.io/yourorg

# Rebuild all three images
make openshift-build REGISTRY=$REGISTRY TAG=latest

# Push to registry
make openshift-push REGISTRY=$REGISTRY TAG=latest
```

### 2. Update Deployment

```bash
# Delete existing deployment
oc delete -k deploy/openshift/overlays/dev

# Redeploy with new images
oc apply -k deploy/openshift/overlays/dev
```

### 3. Verify

```bash
# Check pods are running
oc get pods -n ovnkube-viewer-dev

# Check frontend port
oc get svc ovnkube-viewer-frontend -n ovnkube-viewer-dev
# Should show port 8080

# Test frontend
oc port-forward -n ovnkube-viewer-dev svc/ovnkube-viewer-frontend 8080:8080
curl http://localhost:8080
```

## 🎁 Benefits

### Security

- ✅ All containers run as non-root (UID 1001)
- ✅ Minimal runtime images (ubi-micro: 40MB)
- ✅ No shell in backend runtime (smaller attack surface)
- ✅ Regular security updates from Red Hat

### Enterprise

- ✅ Red Hat Enterprise Linux base
- ✅ Subscription support available
- ✅ Compliance ready
- ✅ CVE scanning and patching

### OpenShift

- ✅ Optimized for OpenShift
- ✅ Better integration
- ✅ Meets restrictive SCC requirements
- ✅ Consistent with RHEL ecosystem

## 📊 Image Sizes

| Component | Alpine | UBI9 | Change |
|-----------|--------|------|--------|
| Backend | ~20MB | ~40MB | +20MB (ubi-micro) |
| Frontend | ~50MB | ~180MB | +130MB (nginx-126) |
| Plugin | ~50MB | ~180MB | +130MB (nginx-126) |

**Note:** Slightly larger images, but benefits outweigh size increase.

## 🔍 What to Test

### Backend

```bash
# Build
podman build -t backend-test:latest -f backend/Dockerfile backend/

# Run
podman run -p 8080:8080 backend-test:latest

# Test
curl http://localhost:8080/health
# Should return: {"status":"ok"}
```

### Frontend

```bash
# Build
podman build -t frontend-test:latest -f frontend/Dockerfile frontend/

# Run
podman run -p 8080:8080 frontend-test:latest

# Test
curl http://localhost:8080
# Should return HTML
```

### Plugin

```bash
# Build
podman build -t plugin-test:latest -f frontend/Dockerfile.plugin frontend/

# Run
podman run -p 9443:9443 plugin-test:latest

# Test
curl http://localhost:9443/plugin-manifest.json
# Should return plugin manifest
```

## 🐛 Troubleshooting

### Build Fails: Cannot pull UBI image

**Issue:** Network or authentication issue

**Solution:**
```bash
# UBI images are free, but login for better rate limits
podman login registry.redhat.io
```

### Container Won't Start: Permission denied

**Issue:** SELinux or file permissions

**Solution:**
```bash
# Check logs
oc logs <pod-name>

# Verify user
oc exec <pod-name> -- id
# Should show: uid=1001
```

### Frontend Returns 403/404

**Issue:** Nginx path misconfiguration

**Solution:**
```bash
# Check nginx config
oc exec <frontend-pod> -- cat /etc/nginx/conf.d/default.conf

# Verify files exist
oc exec <frontend-pod> -- ls -la /opt/app-root/src
```

### Port Mismatch Errors

**Issue:** Service still pointing to old port

**Solution:**
```bash
# Check service
oc get svc ovnkube-viewer-frontend -o yaml | grep port

# Should show port 8080, not 80
# If wrong, reapply manifests:
oc apply -k deploy/openshift/overlays/dev
```

## 📚 Documentation

See these guides for more details:

- **[UBI Base Images Guide](docs/ubi-base-images.md)** - Complete UBI documentation
- **[Building Containers](docs/building-containers.md)** - Build instructions
- **[CHANGELOG](docs/CHANGELOG.md)** - All changes

## 🔐 Authentication

UBI images are **freely available** - no Red Hat subscription required:

```bash
# Pull UBI images (no auth needed)
podman pull registry.redhat.io/ubi9/ubi-micro:9.8-1791895366
podman pull registry.redhat.io/ubi9/go-toolset:9.8-1793107762
podman pull registry.redhat.io/ubi9/nodejs-24:9.8-1793067767
podman pull registry.redhat.io/ubi9/nginx-126:9.8-1793065662

# Build without authentication
make openshift-build REGISTRY=quay.io/yourorg TAG=latest
```

Optional: Login for better rate limits:
```bash
podman login registry.redhat.io
```

## ✨ Quick Commands

```bash
# Complete migration workflow
export REGISTRY=quay.io/yourorg

# 1. Rebuild images
make openshift-build REGISTRY=$REGISTRY TAG=ubi

# 2. Push to registry
make openshift-push REGISTRY=$REGISTRY TAG=ubi

# 3. Update deployment
oc apply -k deploy/openshift/overlays/dev

# 4. Verify
oc get pods -n ovnkube-viewer-dev
oc get svc -n ovnkube-viewer-dev
```

## 🎯 Summary

All three containers now use Red Hat UBI9 base images for better security, enterprise support, and OpenShift optimization. The only breaking change is the frontend port (80 → 8080), which has been updated in all manifests.

**Next Steps:**
1. Rebuild your images
2. Update your deployment
3. Test functionality
4. Enjoy better security! 🎉

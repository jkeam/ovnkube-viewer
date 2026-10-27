# Red Hat UBI Base Images

All containers now use Red Hat Universal Base Images (UBI) instead of Alpine Linux.

## What Changed

### Backend (Go Service)

**Previous (Alpine-based):**
```dockerfile
FROM golang:1.21-alpine AS builder
# ...
FROM alpine:latest
```

**Current (UBI-based):**
```dockerfile
FROM registry.redhat.io/ubi9/go-toolset:9.8-1793107762 AS builder
# ...
FROM registry.redhat.io/ubi9/ubi-micro:9.8-1791895366
```

### Frontend and Plugin (Node.js + Nginx)

**Previous (Alpine-based):**
```dockerfile
FROM node:24-alpine AS builder
# ...
FROM nginx:alpine
```

**Current (UBI-based):**
```dockerfile
FROM registry.redhat.io/ubi9/nodejs-24:9.8-1793067767 AS builder
# ...
FROM registry.redhat.io/ubi9/nginx-124:1-113
```

## Why UBI?

### Benefits

1. **Enterprise Support** - Red Hat subscription support
2. **Security** - Regular security updates and CVE patches
3. **Compliance** - Meets enterprise compliance requirements
4. **Smaller Attack Surface** - ubi-micro is minimal (40MB)
5. **OpenShift Optimized** - Better integration with OpenShift
6. **RHEL Compatible** - Same packages as RHEL
7. **Free to Use** - No subscription required to use or redistribute

### Image Sizes

| Image | Size | Use Case |
|-------|------|----------|
| ubi9 | ~220MB | Full RHEL userspace |
| ubi9-minimal | ~100MB | Minimal package manager (microdnf) |
| ubi9-micro | ~40MB | No package manager, runtime only |

We use **ubi-micro** for the runtime as our Go binary is statically compiled.

## Image Details

### Backend Images

#### Build Stage: go-toolset

**Image:** `registry.redhat.io/ubi9/go-toolset:9.8-1793107762`

- **Go Version:** 1.26.10 (Red Hat 1.26.10-1.el9_8)
- **Base:** RHEL 9
- **User:** Non-root (UID 1001)
- **Workdir:** `/opt/app-root/src`
- **Includes:** Go toolchain, git, make, gcc

#### Runtime Stage: ubi-micro

**Image:** `registry.redhat.io/ubi9/ubi-micro:9.8-1791895366`

- **Size:** ~40MB
- **Minimal:** No package manager, no shell
- **Includes:** Essential libraries only
- **Perfect for:** Static Go binaries

### Frontend/Plugin Images

#### Build Stage: nodejs-24

**Image:** `registry.redhat.io/ubi9/nodejs-24:9.8-1793067767`

- **Node.js Version:** 24.x LTS
- **Base:** RHEL 9
- **User:** Non-root (UID 1001)
- **Workdir:** `/opt/app-root/src`
- **Includes:** npm, node, build tools

#### Runtime Stage: nginx-124

**Image:** `registry.redhat.io/ubi9/nginx-124:1-113`

- **Nginx Version:** 1.24.x
- **Base:** RHEL 9
- **User:** Non-root (UID 1001)
- **Default Port:** 8080 (non-privileged)
- **Document Root:** `/opt/app-root/src`
- **Includes:** Nginx web server

## Building with UBI

### Standard Build

```bash
# Build as usual
podman build -t ovnkube-viewer-backend:latest -f backend/Dockerfile backend/

# Or with Make
make openshift-build REGISTRY=quay.io/yourorg TAG=latest
```

### Authentication

UBI images are **freely available** and don't require authentication to pull:

```bash
# No login needed for UBI images
podman pull registry.redhat.io/ubi9/ubi-micro:9.8-1791895366

# Build will work without authentication
podman build -t ovnkube-viewer-backend:latest -f backend/Dockerfile backend/
```

However, for better rate limits and access to all images, you can authenticate:

```bash
# Login with Red Hat account (optional)
podman login registry.redhat.io
```

### Testing Locally

```bash
# Build the image
podman build -t ovnkube-viewer-backend:latest -f backend/Dockerfile backend/

# Run the container
podman run -p 8080:8080 ovnkube-viewer-backend:latest

# Test health endpoint
curl http://localhost:8080/health
```

## Key Differences from Alpine

### User Permissions

**Alpine:**
- Runs as root by default
- Manual user creation needed

**UBI:**
- Runs as non-root (UID 1001)
- Security best practice built-in

### Package Management

**Alpine:**
- Uses `apk` package manager
- Adds ca-certificates: `apk --no-cache add ca-certificates`

**UBI-micro:**
- No package manager
- ca-certificates already included
- Can't install additional packages

### File Paths

**Alpine:**
- Backend: `/root/` for application
- Frontend: `/usr/share/nginx/html` for web content
- Standard Linux paths

**UBI:**
- All services: `/opt/app-root/src` for application/content
- Non-root user paths
- Consistent across all UBI images

### Port Changes

**Alpine Nginx:**
- Listens on port 80 (requires root or capabilities)
- Frontend service exposed on port 80

**UBI Nginx:**
- Listens on port 8080 (non-privileged port)
- Frontend service exposed on port 8080
- Plugin service still on port 9443
- No special capabilities needed

**Important:** Service manifests updated to use port 8080 for frontend.

## Security Considerations

### Running as Non-Root

The runtime container runs as UID 1001 (non-root):

```dockerfile
USER 1001
```

This is a security best practice and required by many OpenShift clusters.

### Static Binary

We compile with `CGO_ENABLED=0` to create a static binary:

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
```

This ensures:
- No dynamic library dependencies
- Works in minimal runtime (ubi-micro)
- Better security (smaller attack surface)

### No Shell

ubi-micro doesn't include a shell, which means:
- ✅ Smaller attack surface
- ✅ Can't exec into container for malicious purposes
- ⚠️ Can't `kubectl exec` for debugging (use ephemeral containers instead)

## Troubleshooting

### Build Fails: "permission denied"

**Issue:** Go toolset runs as non-root

**Solution:** Files are copied with correct permissions automatically

```dockerfile
# This works correctly
COPY . .
RUN go build -o server ./cmd/server
```

### Build Fails: "cannot pull image"

**Issue:** Rate limiting or network issues

**Solutions:**

1. **Retry the build** - Temporary network issue
2. **Login to registry** - Better rate limits
   ```bash
   podman login registry.redhat.io
   ```
3. **Check network** - Ensure access to registry.redhat.io

### Runtime: "exec /bin/sh: no such file or directory"

**Issue:** ubi-micro doesn't have a shell

**Solution:** Use the binary directly in CMD

```dockerfile
# ✅ Correct
CMD ["/app/server"]

# ❌ Wrong - no shell
CMD ["sh", "-c", "/app/server"]
```

### Debugging Running Container

**Issue:** Can't exec into ubi-micro container (no shell)

**Solutions:**

1. **Check logs:**
   ```bash
   kubectl logs <pod-name>
   ```

2. **Use ephemeral debug container:**
   ```bash
   kubectl debug -it <pod-name> --image=registry.redhat.io/ubi9/ubi:latest
   ```

3. **Port forward and test:**
   ```bash
   kubectl port-forward <pod-name> 8080:8080
   curl http://localhost:8080/health
   ```

## Image Updates

### Checking for Updates

Check for newer UBI images:

```bash
# List available tags
skopeo list-tags docker://registry.redhat.io/ubi9/go-toolset
skopeo list-tags docker://registry.redhat.io/ubi9/ubi-micro
```

### Updating Images

To update to a newer UBI version:

1. **Check release notes** at [Red Hat Container Catalog](https://catalog.redhat.com/software/containers/search)

2. **Update Dockerfile:**
   ```dockerfile
   FROM registry.redhat.io/ubi9/go-toolset:<new-version>
   # ...
   FROM registry.redhat.io/ubi9/ubi-micro:<new-version>
   ```

3. **Test build:**
   ```bash
   podman build -t test-backend:latest -f backend/Dockerfile backend/
   podman run -p 8080:8080 test-backend:latest
   ```

4. **Update documentation** with new version numbers

## Verification

### Verify UBI Base

Check that your image is using UBI:

```bash
# Inspect the image
podman inspect ovnkube-viewer-backend:latest | jq '.[0].Config.Labels'

# Should see Red Hat labels:
# "com.redhat.component"
# "io.buildah.version"
# "vcs-ref"
```

### Verify Static Binary

Check that the binary is static:

```bash
# Run container and check binary
podman run --rm ovnkube-viewer-backend:latest sh -c "ldd /app/server" 2>&1 || echo "Static binary (good!)"

# Should output: "Static binary (good!)" or "not a dynamic executable"
```

### Verify Non-Root

Check that container runs as non-root:

```bash
# Check user ID
podman run --rm ovnkube-viewer-backend:latest id

# Should output: uid=1001 gid=0(root) groups=0(root)
```

## CI/CD Considerations

### Image Caching

UBI images are large. Use layer caching:

```yaml
# GitHub Actions example
- name: Cache UBI images
  uses: actions/cache@v3
  with:
    path: /var/lib/containers
    key: ubi-images-${{ hashFiles('**/Dockerfile') }}
```

### Scanning

Scan UBI images for vulnerabilities:

```bash
# Using Trivy
trivy image ovnkube-viewer-backend:latest

# Using podman/skopeo
skopeo inspect docker://registry.redhat.io/ubi9/ubi-micro:9.8-1791895366
```

## Resources

- [Red Hat Container Catalog](https://catalog.redhat.com/software/containers/search)
- [UBI Documentation](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/9/html/building_running_and_managing_containers/index)
- [UBI FAQ](https://developers.redhat.com/articles/ubi-faq)
- [Go Toolset Container](https://catalog.redhat.com/software/containers/ubi9/go-toolset)
- [UBI Micro](https://catalog.redhat.com/software/containers/ubi9/ubi-micro)

## Alternatives

If you need different capabilities:

| Need | Image | Size |
|------|-------|------|
| Package manager | ubi9-minimal | ~100MB |
| Full RHEL | ubi9 | ~220MB |
| Debugging tools | ubi9 | ~220MB |
| CGo dependencies | ubi9-minimal | ~100MB |

Update the runtime stage in `backend/Dockerfile` accordingly.

## Support

- **Community:** [Red Hat Developer](https://developers.redhat.com/)
- **Documentation:** [access.redhat.com](https://access.redhat.com/)
- **Registry:** [registry.redhat.io](https://registry.redhat.io/)
- **Catalog:** [catalog.redhat.com](https://catalog.redhat.com/)

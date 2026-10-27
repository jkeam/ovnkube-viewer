# Console Plugin Quick Start

Get the OVN Trace Viewer running in your OpenShift Console in 5 minutes.

## Prerequisites

- OpenShift 4.10+ cluster
- Cluster admin access
- `oc` CLI logged in

## Step 1: Build Images (5 minutes)

```bash
# Set your registry
export REGISTRY=quay.io/yourorg

# Build all images (backend, frontend, plugin)
make openshift-build REGISTRY=$REGISTRY TAG=dev

# Push to registry
make openshift-push REGISTRY=$REGISTRY TAG=dev
```

**Note**: Update `yourorg` to your actual registry organization.

## Step 2: Update Image References (30 seconds)

Edit `deploy/openshift/base/kustomization.yaml`:

```yaml
images:
  - name: ovnkube-viewer-backend
    newName: quay.io/yourorg/ovnkube-viewer-backend  # Update this
    newTag: latest
  - name: ovnkube-viewer-frontend
    newName: quay.io/yourorg/ovnkube-viewer-frontend  # Update this
    newTag: latest
  - name: ovnkube-viewer-plugin
    newName: quay.io/yourorg/ovnkube-viewer-plugin  # Update this
    newTag: latest
```

## Step 3: Deploy to OpenShift (1 minute)

```bash
# Deploy to development
make openshift-deploy OVERLAY=dev

# Wait for pods to be ready
oc get pods -n ovnkube-viewer-dev -w
```

## Step 4: Enable Plugin in Console (30 seconds)

```bash
# Enable the plugin
cd deploy/openshift
./enable-plugin.sh enable

# Check status
./enable-plugin.sh status
```

## Step 5: Access the Plugin (10 seconds)

1. Open your OpenShift Console
2. Click **Observe** in the left sidebar
3. Click **OVN Trace Viewer**

🎉 Done! You should now see the OVN Trace Viewer in your console.

## Quick Test

Try a trace:

1. Enter source: Any pod name or IP
2. Enter destination: Any pod name, service name, or IP
3. Protocol: TCP
4. Click **Trace**
5. View the network flow visualization

## Troubleshooting

### Plugin not appearing?

```bash
# Check plugin status
./enable-plugin.sh status

# Check console pods restarted
oc get pods -n openshift-console

# Check plugin service
oc get svc ovnkube-viewer-plugin -n ovnkube-viewer-dev
```

### Plugin shows error?

```bash
# Check plugin logs
oc logs -n ovnkube-viewer-dev -l component=console-plugin

# Check backend logs
oc logs -n ovnkube-viewer-dev -l component=backend

# Verify backend is accessible
oc exec -n ovnkube-viewer-dev deployment/dev-ovnkube-viewer-plugin -- \
  wget -qO- http://dev-ovnkube-viewer-backend:8080/health
```

## One-Line Deploy

For the impatient:

```bash
export REGISTRY=quay.io/yourorg && \
  make openshift-build openshift-push openshift-deploy OVERLAY=dev REGISTRY=$REGISTRY TAG=dev && \
  ./deploy/openshift/enable-plugin.sh enable
```

Then open your OpenShift Console → Observe → OVN Trace Viewer

## Next Steps

- Review [Console Plugin Guide](CONSOLE-PLUGIN-GUIDE.md) for detailed documentation
- See [Deployment Guide](DEPLOYMENT-GUIDE.md) for troubleshooting
- Check [README.md](README.md) for architecture details

## Disable/Remove

```bash
# Disable plugin (keeps deployment)
./enable-plugin.sh disable

# Remove deployment
make openshift-undeploy OVERLAY=dev
```

## Production Deployment

For production, use versioned tags:

```bash
export REGISTRY=quay.io/yourorg
export VERSION=v1.0.0

# Build with version tag
make openshift-build REGISTRY=$REGISTRY TAG=$VERSION

# Update prod overlay to use $VERSION
# Then deploy
make openshift-deploy OVERLAY=prod REGISTRY=$REGISTRY TAG=$VERSION

# Enable plugin
./deploy/openshift/enable-plugin.sh enable
```

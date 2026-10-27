# OpenShift Deployment Guide

Complete guide for deploying OVN-Kubernetes Trace Viewer on OpenShift.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Prerequisites](#prerequisites)
3. [Building Images](#building-images)
4. [Deploying to Development](#deploying-to-development)
5. [Deploying to Production](#deploying-to-production)
6. [Post-Deployment](#post-deployment)
7. [Troubleshooting](#troubleshooting)

## Quick Start

Deploy to development environment in 3 steps:

```bash
# 1. Build and push images
make openshift-build openshift-push TAG=dev REGISTRY=quay.io/yourorg

# 2. Deploy to OpenShift
make openshift-deploy OVERLAY=dev

# 3. Get the route URL
oc get route -n ovnkube-viewer-dev
```

## Prerequisites

### Required Tools

- `oc` CLI (OpenShift CLI) - [Install Guide](https://docs.openshift.com/container-platform/latest/cli_reference/openshift_cli/getting-started-cli.html)
- `podman` or `docker` for building images
- `kubectl` with kustomize support (usually bundled with `oc`)

Verify tools:

```bash
oc version
podman version
oc kustomize version
```

### OpenShift Cluster Requirements

- OpenShift 4.10+ (with OVN-Kubernetes networking)
- Cluster admin access (for ClusterRole/ClusterRoleBinding creation)
- Access to an image registry (e.g., Quay.io, Docker Hub, OpenShift internal registry)

### Access Requirements

Login to your OpenShift cluster:

```bash
oc login <your-cluster-url> --token=<your-token>
```

Verify cluster access:

```bash
oc whoami
oc cluster-info
```

Verify OVN-Kubernetes is running:

```bash
oc get pods -n openshift-ovn-kubernetes
```

You should see `ovnkube-master` and `ovnkube-node` pods running.

## Building Images

### Option 1: Using Make (Recommended)

```bash
# Set your registry (update to your actual registry)
export REGISTRY=quay.io/yourorg

# Build images
make openshift-build TAG=latest

# Push to registry
make openshift-push TAG=latest

# Or build and push in one command
make openshift-build openshift-push TAG=latest
```

### Option 2: Manual Build

```bash
# Build backend
cd backend
podman build -t quay.io/yourorg/ovnkube-viewer-backend:latest .

# Build frontend
cd ../frontend
podman build -t quay.io/yourorg/ovnkube-viewer-frontend:latest .

# Push images
podman login quay.io
podman push quay.io/yourorg/ovnkube-viewer-backend:latest
podman push quay.io/yourorg/ovnkube-viewer-frontend:latest
```

### Using OpenShift Internal Registry

If you want to use OpenShift's internal registry:

```bash
# Get the internal registry URL
INTERNAL_REGISTRY=$(oc get route default-route -n openshift-image-registry -o jsonpath='{.spec.host}')

# Login to the registry
podman login -u $(oc whoami) -p $(oc whoami -t) $INTERNAL_REGISTRY

# Build with internal registry
export REGISTRY=$INTERNAL_REGISTRY/ovnkube-viewer
make openshift-build openshift-push TAG=latest
```

Then update the kustomization.yaml to use the internal registry URL.

## Deploying to Development

### Step 1: Review the Configuration

Check the dev overlay configuration:

```bash
cat deploy/openshift/overlays/dev/kustomization.yaml
```

Key settings for dev:
- 1 replica for each component
- Debug logging enabled
- Dev-tagged images
- Namespace: `ovnkube-viewer-dev`

### Step 2: Customize Image References

Edit `deploy/openshift/overlays/dev/kustomization.yaml`:

```yaml
images:
  - name: quay.io/yourorg/ovnkube-viewer-backend
    newTag: dev  # or your preferred tag
  - name: quay.io/yourorg/ovnkube-viewer-frontend
    newTag: dev
```

Or edit `deploy/openshift/base/kustomization.yaml` to set the default registry:

```yaml
images:
  - name: ovnkube-viewer-backend
    newName: quay.io/yourorg/ovnkube-viewer-backend  # Update yourorg
    newTag: latest
  - name: ovnkube-viewer-frontend
    newName: quay.io/yourorg/ovnkube-viewer-frontend  # Update yourorg
    newTag: latest
```

### Step 3: Validate Manifests

Dry-run to see what will be created:

```bash
make openshift-validate OVERLAY=dev
```

Or manually:

```bash
oc kustomize deploy/openshift/overlays/dev
```

### Step 4: Deploy

```bash
make openshift-deploy OVERLAY=dev
```

Or manually:

```bash
oc apply -k deploy/openshift/overlays/dev
```

### Step 5: Verify Deployment

```bash
# Check deployment status
make openshift-status OVERLAY=dev

# Watch pods come up
oc get pods -n ovnkube-viewer-dev -w

# Check logs
oc logs -n ovnkube-viewer-dev -l component=backend --tail=50
oc logs -n ovnkube-viewer-dev -l component=frontend --tail=50
```

### Step 6: Access the Application

Get the route URL:

```bash
oc get route -n ovnkube-viewer-dev
```

Open the URL in your browser.

## Deploying to Production

### Step 1: Tag Images with Version

```bash
# Tag with semantic version
VERSION=v1.0.0
podman tag quay.io/yourorg/ovnkube-viewer-backend:latest quay.io/yourorg/ovnkube-viewer-backend:$VERSION
podman tag quay.io/yourorg/ovnkube-viewer-frontend:latest quay.io/yourorg/ovnkube-viewer-frontend:$VERSION

# Push versioned images
podman push quay.io/yourorg/ovnkube-viewer-backend:$VERSION
podman push quay.io/yourorg/ovnkube-viewer-frontend:$VERSION
```

### Step 2: Update Production Configuration

Edit `deploy/openshift/overlays/prod/kustomization.yaml`:

```yaml
images:
  - name: quay.io/yourorg/ovnkube-viewer-backend
    newTag: v1.0.0  # Update to your version
  - name: quay.io/yourorg/ovnkube-viewer-frontend
    newTag: v1.0.0
```

### Step 3: Review Production Settings

```bash
cat deploy/openshift/overlays/prod/kustomization.yaml
```

Key settings for prod:
- 3 replicas for each component
- Higher resource limits
- Info-level logging
- Namespace: `ovnkube-viewer`

### Step 4: Validate Production Manifests

```bash
make openshift-validate OVERLAY=prod
```

### Step 5: Deploy to Production

```bash
make openshift-deploy OVERLAY=prod
```

### Step 6: Monitor Rollout

```bash
# Watch the rollout
oc rollout status deployment/ovnkube-viewer-backend -n ovnkube-viewer
oc rollout status deployment/ovnkube-viewer-frontend -n ovnkube-viewer

# Check status
make openshift-status OVERLAY=prod

# Check pod health
oc get pods -n ovnkube-viewer -l app=ovnkube-viewer
```

## Post-Deployment

### Testing the Application

1. **Get the route URL:**
   ```bash
   ROUTE_URL=$(oc get route ovnkube-viewer -n ovnkube-viewer -o jsonpath='{.spec.host}')
   echo "https://$ROUTE_URL"
   ```

2. **Test backend health:**
   ```bash
   curl -k https://$ROUTE_URL/api/health
   ```

3. **Access the UI:**
   Open `https://$ROUTE_URL` in your browser

4. **Test a trace:**
   - Navigate to the UI
   - Enter source and destination pod IPs or service names
   - Select protocol (TCP/UDP)
   - Click "Trace"
   - Verify the trace results are displayed

### Monitoring

#### View Logs

```bash
# Backend logs
oc logs -f deployment/ovnkube-viewer-backend -n ovnkube-viewer

# Frontend logs
oc logs -f deployment/ovnkube-viewer-frontend -n ovnkube-viewer

# All logs
oc logs -f -l app=ovnkube-viewer -n ovnkube-viewer
```

#### Check Resource Usage

```bash
# Pod metrics (requires metrics-server)
oc adm top pods -n ovnkube-viewer

# Node metrics
oc adm top nodes
```

#### View Events

```bash
oc get events -n ovnkube-viewer --sort-by='.lastTimestamp'
```

### Scaling

Scale deployments as needed:

```bash
# Scale backend
oc scale deployment/ovnkube-viewer-backend --replicas=5 -n ovnkube-viewer

# Scale frontend
oc scale deployment/ovnkube-viewer-frontend --replicas=5 -n ovnkube-viewer
```

### Updating the Application

To update to a new version:

```bash
# 1. Build and push new version
VERSION=v1.1.0
make openshift-build openshift-push TAG=$VERSION

# 2. Update kustomization.yaml with new tag
# Edit deploy/openshift/overlays/prod/kustomization.yaml

# 3. Apply the update
make openshift-deploy OVERLAY=prod

# 4. Watch the rollout
oc rollout status deployment/ovnkube-viewer-backend -n ovnkube-viewer
```

## Troubleshooting

### Pods Not Starting

**Check pod status:**
```bash
oc get pods -n ovnkube-viewer
oc describe pod <pod-name> -n ovnkube-viewer
```

**Common issues:**

1. **Image pull errors:**
   - Verify image exists: `podman pull <image-url>`
   - Check image pull secrets if using private registry
   - Verify registry URL in kustomization.yaml

2. **Insufficient permissions:**
   - Verify ServiceAccount exists: `oc get sa ovnkube-viewer -n ovnkube-viewer`
   - Check ClusterRoleBinding: `oc get clusterrolebinding ovnkube-viewer`

3. **Resource limits:**
   - Check node resources: `oc adm top nodes`
   - Adjust resource requests/limits in deployment

### Permission Denied When Accessing OVN

**Verify OVN access:**

```bash
# Get an ovnkube-master pod name
OVN_POD=$(oc get pods -n openshift-ovn-kubernetes -l app=ovnkube-master -o jsonpath='{.items[0].metadata.name}')

# Test exec access from backend pod
BACKEND_POD=$(oc get pods -n ovnkube-viewer -l component=backend -o jsonpath='{.items[0].metadata.name}')

# Check if backend can list OVN pods
oc exec -n ovnkube-viewer $BACKEND_POD -- kubectl get pods -n openshift-ovn-kubernetes
```

**If exec fails:**

1. Verify ClusterRole has `pods/exec` permission
2. Check that ServiceAccount is bound to ClusterRole
3. Verify the backend pod is using the correct ServiceAccount

### Cannot Access Application via Route

**Check route:**
```bash
oc get route -n ovnkube-viewer
oc describe route ovnkube-viewer -n ovnkube-viewer
```

**Verify route DNS:**
```bash
ROUTE_HOST=$(oc get route ovnkube-viewer -n ovnkube-viewer -o jsonpath='{.spec.host}')
nslookup $ROUTE_HOST
```

**Test directly to service:**
```bash
# Port-forward to test service
oc port-forward -n ovnkube-viewer svc/ovnkube-viewer-frontend 8080:80

# In another terminal
curl http://localhost:8080
```

### Backend Cannot Reach Frontend or Vice Versa

**Test service connectivity:**

```bash
# From backend pod
BACKEND_POD=$(oc get pods -n ovnkube-viewer -l component=backend -o jsonpath='{.items[0].metadata.name}')
oc exec -n ovnkube-viewer $BACKEND_POD -- wget -qO- http://ovnkube-viewer-frontend

# From frontend pod
FRONTEND_POD=$(oc get pods -n ovnkube-viewer -l component=frontend -o jsonpath='{.items[0].metadata.name}')
oc exec -n ovnkube-viewer $FRONTEND_POD -- wget -qO- http://ovnkube-viewer-backend:8080/health
```

### High Memory/CPU Usage

**Check resource usage:**
```bash
oc adm top pods -n ovnkube-viewer
```

**Adjust resources:**

Edit the deployment directly:
```bash
oc edit deployment ovnkube-viewer-backend -n ovnkube-viewer
```

Or update the kustomization overlay and redeploy.

### Viewing Full Manifests

To see the complete generated manifests:

```bash
# Development
oc kustomize deploy/openshift/overlays/dev > /tmp/dev-manifests.yaml
cat /tmp/dev-manifests.yaml

# Production
oc kustomize deploy/openshift/overlays/prod > /tmp/prod-manifests.yaml
cat /tmp/prod-manifests.yaml
```

## Cleanup

### Remove Development Deployment

```bash
make openshift-undeploy OVERLAY=dev
```

Or manually:

```bash
oc delete -k deploy/openshift/overlays/dev
```

### Remove Production Deployment

```bash
make openshift-undeploy OVERLAY=prod
```

### Complete Cleanup (Including Cluster Resources)

```bash
# Delete namespaces
oc delete namespace ovnkube-viewer
oc delete namespace ovnkube-viewer-dev

# Delete cluster-scoped resources
oc delete clusterrole ovnkube-viewer
oc delete clusterrolebinding ovnkube-viewer
oc delete clusterrolebinding dev-ovnkube-viewer
```

## Advanced Topics

### Using a Custom Domain

Edit the route to use a custom domain:

```yaml
# Add to route.yaml
spec:
  host: ovnkube-viewer.example.com
  to:
    kind: Service
    name: ovnkube-viewer-frontend
```

### Adding Network Policies

Create NetworkPolicy to restrict traffic:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: ovnkube-viewer-backend
  namespace: ovnkube-viewer
spec:
  podSelector:
    matchLabels:
      component: backend
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              component: frontend
      ports:
        - protocol: TCP
          port: 8080
  egress:
    - to:
        - namespaceSelector:
            matchLabels:
              name: openshift-ovn-kubernetes
```

### Enabling Horizontal Pod Autoscaling

```bash
# Create HPA for backend
oc autoscale deployment ovnkube-viewer-backend \
  --min=2 --max=10 \
  --cpu-percent=80 \
  -n ovnkube-viewer
```

### Monitoring with Prometheus

The application exposes standard health endpoints. To add Prometheus monitoring:

1. Add ServiceMonitor resource
2. Instrument backend with Prometheus metrics
3. Create Grafana dashboards

## Next Steps

- Review [README.md](README.md) for architecture details
- Check [../../docs/implementation-plan.md](../../docs/implementation-plan.md) for feature roadmap
- Explore Phase 2 features (Helm charts, persistence)
- Plan Phase 3 migration to Operator pattern

## Support

For issues and questions:
- Check troubleshooting section above
- Review OpenShift logs and events
- Consult [OVN-Kubernetes documentation](https://docs.openshift.com/container-platform/latest/networking/ovn_kubernetes_network_provider/about-ovn-kubernetes.html)

# OpenShift Deployment with Kustomize

This directory contains OpenShift manifests for deploying the OVN-Kubernetes Trace Viewer using Kustomize.

## Directory Structure

```
openshift/
├── base/                           # Base manifests
│   ├── namespace.yaml              # Namespace definition
│   ├── serviceaccount.yaml         # ServiceAccount for the application
│   ├── clusterrole.yaml            # ClusterRole with OVN access permissions
│   ├── clusterrolebinding.yaml     # Binds ClusterRole to ServiceAccount
│   ├── configmap.yaml              # Application configuration
│   ├── backend-deployment.yaml     # Backend Go service deployment
│   ├── backend-service.yaml        # Backend service
│   ├── frontend-deployment.yaml    # Frontend React app deployment
│   ├── frontend-service.yaml       # Frontend service
│   ├── route.yaml                  # OpenShift Route for external access
│   └── kustomization.yaml          # Base kustomization
├── overlays/
│   ├── dev/                        # Development overlay
│   │   └── kustomization.yaml
│   └── prod/                       # Production overlay
│       └── kustomization.yaml
└── README.md                       # This file
```

## Permissions

The application requires elevated permissions to examine cluster network traffic. The ServiceAccount is granted the following cluster-wide permissions:

### Core Resources
- **Pods**: `get`, `list`, `watch`, `exec` - Required to access ovnkube-master and ovnkube-node pods to execute `ovn-trace` commands
- **Services & Endpoints**: `get`, `list`, `watch` - Required for service-based tracing
- **Namespaces**: `get`, `list`, `watch` - Required for multi-namespace tracing
- **Nodes**: `get`, `list`, `watch` - Required for node-level network information

### Network Resources
- **NetworkPolicies**: `get`, `list`, `watch` - Required to understand traffic filtering rules
- **OpenShift EgressNetworkPolicies**: `get`, `list`, `watch` - Required for egress policy analysis
- **HostSubnets**: `get`, `list`, `watch` - Required for node network configuration
- **NetNamespaces**: `get`, `list`, `watch` - Required for namespace network isolation info

### OpenShift Routes
- **Routes**: `get`, `list`, `watch` - Required for service discovery

## Prerequisites

1. OpenShift cluster with OVN-Kubernetes networking
2. `oc` CLI tool installed and configured
3. Cluster admin access (required to create ClusterRole and ClusterRoleBinding)
4. Container images pushed to an accessible registry

## Building and Pushing Images

Before deploying, you need to build and push the container images:

```bash
# Build backend image
cd backend
podman build -t quay.io/yourorg/ovnkube-viewer-backend:latest .

# Build frontend image
cd ../frontend
podman build -t quay.io/yourorg/ovnkube-viewer-frontend:latest .

# Push images
podman push quay.io/yourorg/ovnkube-viewer-backend:latest
podman push quay.io/yourorg/ovnkube-viewer-frontend:latest
```

**Note**: Update `yourorg` in the image names and in the kustomization files to match your registry organization.

## Deployment

### Option 1: Deploy Development Environment

```bash
# Review the generated manifests
oc kustomize deploy/openshift/overlays/dev

# Apply the manifests
oc apply -k deploy/openshift/overlays/dev

# Check deployment status
oc get pods -n ovnkube-viewer-dev
oc get route -n ovnkube-viewer-dev
```

### Option 2: Deploy Production Environment

```bash
# Review the generated manifests
oc kustomize deploy/openshift/overlays/prod

# Apply the manifests
oc apply -k deploy/openshift/overlays/prod

# Check deployment status
oc get pods -n ovnkube-viewer
oc get route -n ovnkube-viewer
```

### Option 3: Deploy Base (No Overlay)

```bash
# Deploy base manifests directly
oc apply -k deploy/openshift/base
```

## Accessing the Application

After deployment, get the route URL:

```bash
# For dev
oc get route -n ovnkube-viewer-dev

# For prod
oc get route -n ovnkube-viewer
```

Access the application in your browser using the route hostname.

## Configuration

### Customizing Images

Update the image references in the overlay kustomization.yaml files:

```yaml
images:
  - name: quay.io/yourorg/ovnkube-viewer-backend
    newTag: your-tag
  - name: quay.io/yourorg/ovnkube-viewer-frontend
    newTag: your-tag
```

### Customizing Configuration

The ConfigMap can be customized in overlays using configMapGenerator:

```yaml
configMapGenerator:
  - name: ovnkube-viewer-config
    behavior: merge
    literals:
      - CORS_ALLOWED_ORIGINS=https://your-domain.com
      - LOG_LEVEL=info
      - TRACE_TIMEOUT=60s
```

### OVN Configuration

By default, the application looks for OVN components in the `openshift-ovn-kubernetes` namespace with labels:
- `app=ovnkube-master` for master pods
- `app=ovnkube-node` for node pods

These can be customized in the ConfigMap if your cluster uses different labels.

## Scaling

### Development
- Backend: 1 replica
- Frontend: 1 replica

### Production
- Backend: 3 replicas
- Frontend: 3 replicas

Scale manually:

```bash
oc scale deployment/ovnkube-viewer-backend --replicas=5 -n ovnkube-viewer
```

## Monitoring

### Check Pod Status

```bash
oc get pods -n ovnkube-viewer -l app=ovnkube-viewer
```

### View Logs

```bash
# Backend logs
oc logs -f deployment/ovnkube-viewer-backend -n ovnkube-viewer

# Frontend logs
oc logs -f deployment/ovnkube-viewer-frontend -n ovnkube-viewer
```

### Health Checks

The backend exposes a health endpoint:

```bash
oc exec -n ovnkube-viewer deployment/ovnkube-viewer-backend -- wget -qO- http://localhost:8080/health
```

## Security Considerations

### Elevated Permissions
This application requires cluster-wide read access and exec permissions to function. These are necessary to:
- Execute `ovn-trace` commands inside ovnkube pods
- Discover pods, services, and network policies for trace operations
- Access OVN databases for network topology information

### Network Policies
Consider implementing NetworkPolicies to restrict:
- Ingress traffic to only necessary sources
- Egress traffic to only OVN components and Kubernetes API

### RBAC Audit
Regularly review the ClusterRole permissions to ensure they align with your security policies:

```bash
oc get clusterrole ovnkube-viewer -o yaml
```

## Troubleshooting

### Pods Not Starting

```bash
# Check events
oc get events -n ovnkube-viewer --sort-by='.lastTimestamp'

# Check pod details
oc describe pod -n ovnkube-viewer -l component=backend
```

### Permission Denied Errors

Verify the ServiceAccount has the correct ClusterRoleBinding:

```bash
oc get clusterrolebinding ovnkube-viewer -o yaml
```

### Cannot Access OVN Pods

Verify OVN pods are running:

```bash
oc get pods -n openshift-ovn-kubernetes
```

Test exec access to an OVN pod:

```bash
oc exec -n openshift-ovn-kubernetes <ovnkube-master-pod> -- ovn-nbctl show
```

### Image Pull Failures

Ensure the image registry is accessible and images exist:

```bash
podman pull quay.io/yourorg/ovnkube-viewer-backend:latest
```

## Cleanup

### Remove Development Deployment

```bash
oc delete -k deploy/openshift/overlays/dev
```

### Remove Production Deployment

```bash
oc delete -k deploy/openshift/overlays/prod
```

### Remove Base Deployment

```bash
oc delete -k deploy/openshift/base
```

### Complete Cleanup

```bash
# Delete namespace (this will delete all resources in it)
oc delete namespace ovnkube-viewer

# Delete cluster-scoped resources
oc delete clusterrole ovnkube-viewer
oc delete clusterrolebinding ovnkube-viewer
```

## Development Workflow

1. Make changes to source code
2. Build and push new images with a development tag
3. Update the image tag in `overlays/dev/kustomization.yaml`
4. Apply the changes:
   ```bash
   oc apply -k deploy/openshift/overlays/dev
   ```
5. Watch the rollout:
   ```bash
   oc rollout status deployment/dev-ovnkube-viewer-backend -n ovnkube-viewer-dev
   ```

## Production Deployment Workflow

1. Tag images with semantic version:
   ```bash
   podman tag quay.io/yourorg/ovnkube-viewer-backend:latest quay.io/yourorg/ovnkube-viewer-backend:v1.0.0
   podman push quay.io/yourorg/ovnkube-viewer-backend:v1.0.0
   ```
2. Update version in `overlays/prod/kustomization.yaml`
3. Apply to production:
   ```bash
   oc apply -k deploy/openshift/overlays/prod
   ```
4. Monitor the deployment:
   ```bash
   oc get pods -n ovnkube-viewer -w
   ```

## Additional Resources

- [Kustomize Documentation](https://kustomize.io/)
- [OpenShift Documentation](https://docs.openshift.com/)
- [OVN-Kubernetes Tracing](https://docs.redhat.com/en/documentation/openshift_container_platform/4.22/html/ovn-kubernetes_network_plugin/ovn-kubernetes-tracing-using-ovntrace)

# OpenShift Console Plugin Guide

This guide explains how to deploy the OVN Trace Viewer as an OpenShift Console dynamic plugin, integrating it directly into the OpenShift web console.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Prerequisites](#prerequisites)
4. [Building the Plugin](#building-the-plugin)
5. [Deploying the Plugin](#deploying-the-plugin)
6. [Enabling the Plugin](#enabling-the-plugin)
7. [Accessing the Plugin](#accessing-the-plugin)
8. [Development](#development)
9. [Troubleshooting](#troubleshooting)

## Overview

The OVN Trace Viewer can be deployed in two modes:

1. **Standalone Mode**: Accessed via an OpenShift Route (covered in DEPLOYMENT-GUIDE.md)
2. **Console Plugin Mode**: Integrated into the OpenShift Console (this guide)

Console plugin mode provides:
- Native integration into the OpenShift web console
- Appears in the "Observe" menu
- Shared authentication and authorization with console
- Consistent UI/UX with other OpenShift features
- No separate route needed - access through the main console

## Architecture

```
┌─────────────────────────────────────────────────────┐
│         OpenShift Console (Port 443)                │
│  ┌───────────────────────────────────────────────┐  │
│  │  Navigation: Observe > OVN Trace Viewer       │  │
│  └───────────────────────────────────────────────┘  │
│                       │                              │
│                       ▼                              │
│  ┌───────────────────────────────────────────────┐  │
│  │  Plugin Service (Port 9443)                   │  │
│  │  - Serves plugin JavaScript bundles           │  │
│  │  - Serves plugin manifest                     │  │
│  │  - CORS enabled for console integration       │  │
│  └───────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                       │
                       ▼ (Proxied API calls)
          ┌─────────────────────────────┐
          │  Backend API (Port 8080)    │
          │  - Execute OVN traces       │
          │  - Parse results            │
          └─────────────────────────────┘
                       │
                       ▼
          ┌─────────────────────────────┐
          │  OVN Pods (ovnkube-master)  │
          │  - ovn-trace execution      │
          └─────────────────────────────┘
```

### Components

1. **Console Plugin Service** (`ovnkube-viewer-plugin`)
   - Nginx serving plugin assets
   - Port 9443 (standard for console plugins)
   - Serves `plugin-manifest.json` and JavaScript bundles

2. **Backend API Service** (`ovnkube-viewer-backend`)
   - Unchanged from standalone mode
   - Proxied through ConsolePlugin CR

3. **Frontend (Standalone)** (`ovnkube-viewer-frontend`)
   - Optional - can be deployed alongside plugin for standalone access

4. **ConsolePlugin Custom Resource**
   - Registers plugin with OpenShift Console
   - Defines service endpoint and proxy rules

## Prerequisites

### Required

- OpenShift 4.10+ with Console Operator
- Cluster admin access (to create ConsolePlugin CR)
- `oc` CLI tool
- `npm` and `node` (for building the plugin)
- Container image registry access

### Check Console Operator

```bash
# Verify Console Operator is running
oc get pods -n openshift-console-operator

# Check Console version
oc get clusteroperator console
```

## Building the Plugin

### Step 1: Install Dependencies

```bash
cd frontend
npm install
```

### Step 2: Build Plugin Bundle

```bash
# Build the console plugin
npm run build:plugin

# This creates dist/plugin/ with:
# - plugin-manifest.json
# - JavaScript bundles
# - CSS files
```

### Step 3: Build Container Image

```bash
# From project root
export REGISTRY=quay.io/yourorg

# Build plugin image
podman build -t ${REGISTRY}/ovnkube-viewer-plugin:latest \
  -f frontend/Dockerfile.plugin \
  frontend/

# Push to registry
podman push ${REGISTRY}/ovnkube-viewer-plugin:latest
```

Or use Make:

```bash
make openshift-build REGISTRY=quay.io/yourorg TAG=latest
```

This builds all three images:
- Backend
- Frontend (standalone)
- Console Plugin

## Deploying the Plugin

### Step 1: Update Image References

Edit `deploy/openshift/base/kustomization.yaml`:

```yaml
images:
  - name: ovnkube-viewer-plugin
    newName: quay.io/yourorg/ovnkube-viewer-plugin
    newTag: latest
```

### Step 2: Deploy with Kustomize

```bash
# Deploy to development
oc apply -k deploy/openshift/overlays/dev

# Or production
oc apply -k deploy/openshift/overlays/prod
```

Or use Make:

```bash
make openshift-deploy OVERLAY=dev
```

### Step 3: Verify Deployment

```bash
# Check plugin pod
oc get pods -n ovnkube-viewer -l component=console-plugin

# Check plugin service
oc get svc ovnkube-viewer-plugin -n ovnkube-viewer

# Check ConsolePlugin CR
oc get consoleplugin ovnkube-viewer

# Verify plugin manifest is accessible
PLUGIN_POD=$(oc get pods -n ovnkube-viewer -l component=console-plugin -o jsonpath='{.items[0].metadata.name}')
oc exec -n ovnkube-viewer $PLUGIN_POD -- wget -qO- http://localhost:9443/plugin-manifest.json
```

## Enabling the Plugin

The ConsolePlugin CR is deployed, but you need to enable it in the Console configuration.

### Option 1: Enable via Console UI (Recommended)

1. Log into OpenShift Console as cluster-admin
2. Navigate to: **Administration** → **Cluster Settings** → **Configuration** → **Console** → **Console plugins**
3. Click **Edit**
4. Check the box for **ovnkube-viewer**
5. Click **Save**

The console will reload automatically.

### Option 2: Enable via CLI

```bash
# Patch the console operator configuration
oc patch consoles.operator.openshift.io cluster \
  --type=json \
  --patch '[{"op": "add", "path": "/spec/plugins", "value": ["ovnkube-viewer"]}]'

# If plugins already exist, use:
oc patch consoles.operator.openshift.io cluster \
  --type=json \
  --patch '[{"op": "add", "path": "/spec/plugins/-", "value": "ovnkube-viewer"}]'
```

### Step 3: Verify Plugin is Enabled

```bash
# Check console configuration
oc get consoles.operator.openshift.io cluster -o jsonpath='{.spec.plugins}'

# Should see: ["ovnkube-viewer"]
```

### Step 4: Wait for Console Rollout

The console pods will restart to load the plugin:

```bash
# Watch console pods restart
oc get pods -n openshift-console -w

# Wait for all pods to be ready
oc wait --for=condition=ready pod -l app=console -n openshift-console --timeout=300s
```

## Accessing the Plugin

### Step 1: Open OpenShift Console

Navigate to your OpenShift Console URL (usually `https://console-openshift-console.apps.<cluster-domain>`)

### Step 2: Navigate to Plugin

1. In the left sidebar, click **Observe**
2. Click **OVN Trace Viewer** (appears after "Monitoring")

### Step 3: Use the Trace Viewer

The plugin provides the same functionality as standalone mode:
- Enter source and destination (pod names, IPs, or service names)
- Select protocol (TCP/UDP)
- Click "Trace" to visualize traffic flow
- View results in interactive graph

## Development

### Local Development with Console Bridge

For local development, you can serve the plugin and connect it to a running OpenShift cluster:

```bash
# In terminal 1: Start backend
cd backend
go run ./cmd/server

# In terminal 2: Build and serve plugin
cd frontend
npm run build:plugin
npm run http-server

# The plugin is now available at http://localhost:9001
```

Configure console to use local plugin:

```bash
# Create a development ConsolePlugin pointing to localhost
cat <<EOF | oc apply -f -
apiVersion: console.openshift.io/v1
kind: ConsolePlugin
metadata:
  name: ovnkube-viewer-dev
spec:
  displayName: 'OVN Trace Viewer (Dev)'
  service:
    name: ovnkube-viewer-plugin
    namespace: ovnkube-viewer
    port: 9001
    basePath: /
EOF
```

**Note**: For security, the console must be configured to allow loading plugins from localhost during development.

### Watch Mode

For iterative development:

```bash
# Terminal 1: Backend with auto-reload
cd backend
air  # or: go run ./cmd/server

# Terminal 2: Frontend watch build
cd frontend
npm run build:plugin -- --watch

# Terminal 3: Serve plugin
cd frontend
npm run http-server
```

### Testing Plugin Integration

```bash
# Check plugin manifest
curl http://localhost:9001/plugin-manifest.json

# Verify CORS headers
curl -I http://localhost:9001/plugin-manifest.json
```

## Troubleshooting

### Plugin Not Appearing in Console

**Check if plugin is enabled:**

```bash
oc get consoles.operator.openshift.io cluster -o jsonpath='{.spec.plugins}'
```

**Check ConsolePlugin CR:**

```bash
oc get consoleplugin ovnkube-viewer -o yaml
```

**Check console operator logs:**

```bash
oc logs -n openshift-console-operator deployment/console-operator
```

**Check console pod logs:**

```bash
oc logs -n openshift-console -l app=console --tail=100
```

### Plugin Fails to Load

**Check plugin service:**

```bash
# Verify service exists
oc get svc ovnkube-viewer-plugin -n ovnkube-viewer

# Check endpoints
oc get endpoints ovnkube-viewer-plugin -n ovnkube-viewer
```

**Check plugin pod:**

```bash
# Pod status
oc get pods -n ovnkube-viewer -l component=console-plugin

# Pod logs
oc logs -n ovnkube-viewer -l component=console-plugin
```

**Test plugin manifest:**

```bash
# Port-forward to plugin service
oc port-forward -n ovnkube-viewer svc/ovnkube-viewer-plugin 9443:9443

# In another terminal
curl http://localhost:9443/plugin-manifest.json
```

### CORS Errors

**Check nginx configuration:**

```bash
PLUGIN_POD=$(oc get pods -n ovnkube-viewer -l component=console-plugin -o jsonpath='{.items[0].metadata.name}')
oc exec -n ovnkube-viewer $PLUGIN_POD -- cat /etc/nginx/nginx.conf
```

**Verify CORS headers:**

```bash
oc exec -n ovnkube-viewer $PLUGIN_POD -- \
  curl -I http://localhost:9443/plugin-manifest.json | grep -i access-control
```

Should see:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, OPTIONS
```

### Backend API Errors

The console plugin proxies API requests through the ConsolePlugin CR proxy configuration.

**Check proxy configuration:**

```bash
oc get consoleplugin ovnkube-viewer -o jsonpath='{.spec.proxy}'
```

**Test backend from plugin pod:**

```bash
PLUGIN_POD=$(oc get pods -n ovnkube-viewer -l component=console-plugin -o jsonpath='{.items[0].metadata.name}')
oc exec -n ovnkube-viewer $PLUGIN_POD -- \
  wget -qO- http://ovnkube-viewer-backend:8080/health
```

**Check backend service:**

```bash
oc get svc ovnkube-viewer-backend -n ovnkube-viewer
oc get endpoints ovnkube-viewer-backend -n ovnkube-viewer
```

### Plugin Build Fails

**Check Node.js version:**

```bash
node --version  # Should be 24+
npm --version   # Should be 10+
```

**Clean and rebuild:**

```bash
cd frontend
rm -rf node_modules dist
npm install
npm run build:plugin
```

**Check webpack errors:**

```bash
npm run build:plugin -- --verbose
```

### Console Not Loading Plugin Bundle

**Check browser console** (F12 in most browsers):
- Look for failed network requests to plugin bundles
- Check for JavaScript errors

**Verify plugin manifest:**

```bash
curl http://localhost:9443/plugin-manifest.json | jq .
```

Should contain:
```json
{
  "name": "ovnkube-viewer",
  "version": "0.1.0",
  "displayName": "OVN Trace Viewer",
  "exposedModules": {
    "OVNTraceViewerPage": "./OVNTraceViewerPage"
  }
}
```

## Disabling the Plugin

### Via Console UI

1. Navigate to: **Administration** → **Cluster Settings** → **Configuration** → **Console**
2. Click **Console plugins** tab
3. Uncheck **ovnkube-viewer**
4. Click **Save**

### Via CLI

```bash
# Remove plugin from enabled list
oc patch consoles.operator.openshift.io cluster \
  --type=json \
  --patch '[{"op": "remove", "path": "/spec/plugins/0"}]'
```

**Note**: Adjust the array index if multiple plugins are enabled.

## Uninstalling the Plugin

```bash
# Remove deployment
oc delete -k deploy/openshift/overlays/prod

# Or just the plugin components
oc delete consoleplugin ovnkube-viewer
oc delete deployment ovnkube-viewer-plugin -n ovnkube-viewer
oc delete service ovnkube-viewer-plugin -n ovnkube-viewer
```

## Differences from Standalone Mode

| Aspect | Standalone Mode | Console Plugin Mode |
|--------|----------------|---------------------|
| Access | Via dedicated Route | Via OpenShift Console |
| Authentication | Separate | Shared with Console |
| Navigation | Standalone URL | Console sidebar menu |
| Build | Vite | Webpack (dynamic plugin) |
| Port | 80 (nginx) | 9443 (plugin service) |
| Deployment | frontend-deployment | plugin-deployment |

## Best Practices

1. **Use Both Modes**: Deploy both standalone and plugin for flexibility
2. **Version Tags**: Always use specific version tags in production
3. **Resource Limits**: Set appropriate limits based on usage
4. **Monitoring**: Monitor plugin pod health and performance
5. **Updates**: Test plugin updates in dev before production
6. **RBAC**: Ensure backend ServiceAccount has necessary permissions

## References

- [OpenShift Console Plugin Documentation](https://docs.openshift.com/container-platform/latest/web_console/dynamic-plugin/overview-dynamic-plugin.html)
- [Console Plugin Template](https://github.com/openshift/console-plugin-template)
- [Dynamic Plugin SDK](https://github.com/openshift/console/tree/master/frontend/packages/console-dynamic-plugin-sdk)
- [Red Hat Blog: Developing Console Plugins](https://www.redhat.com/en/blog/developing-an-openshift-dynamic-console-plugin-1)

## Next Steps

After successful deployment:
1. Configure monitoring and alerting
2. Set up automated plugin updates via CI/CD
3. Add custom analytics or telemetry
4. Explore additional console extension points
5. Consider creating multiple plugin pages for different workflows

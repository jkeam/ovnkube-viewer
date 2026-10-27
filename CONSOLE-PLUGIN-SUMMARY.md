# OpenShift Console Plugin - Implementation Summary

This document summarizes the OpenShift Console Plugin implementation for the OVN Trace Viewer.

## Overview

The OVN Trace Viewer has been converted into an OpenShift Console dynamic plugin, allowing it to be integrated directly into the OpenShift web console. Users can now access the trace viewer from the "Observe" menu in the console.

## What Was Added

### Frontend Changes

#### 1. Console Plugin Configuration
- **`frontend/console-extensions.json`** - Defines console extension points (navigation, routes)
- **`frontend/webpack.config.ts`** - Webpack configuration for building the dynamic plugin
- **`frontend/package.json`** - Updated with console plugin dependencies and build scripts

#### 2. Plugin Source Code
- **`frontend/src/plugin/`** - New directory for plugin-specific code
  - `OVNTraceViewerPage.tsx` - Main plugin page component
  - `plugin.ts` - Plugin metadata and configuration

#### 3. Plugin Build & Deployment
- **`frontend/Dockerfile.plugin`** - Dockerfile for building the plugin container
- **`frontend/nginx.plugin.conf`** - Nginx configuration for serving plugin assets

### OpenShift Manifests

#### 4. Console Plugin Resources (`deploy/openshift/base/`)
- **`console-plugin.yaml`** - ConsolePlugin CR that registers the plugin with OpenShift Console
- **`plugin-deployment.yaml`** - Deployment for the plugin service (serves plugin assets)
- **`plugin-service.yaml`** - Service exposing the plugin on port 9443

#### 5. Updated Base Kustomization
- **`deploy/openshift/base/kustomization.yaml`** - Added plugin resources and image references

#### 6. Updated Overlays
- **`deploy/openshift/overlays/dev/kustomization.yaml`** - Added plugin image tag for dev
- **`deploy/openshift/overlays/prod/kustomization.yaml`** - Added plugin image tag for prod

### Scripts & Tools

#### 7. Plugin Management Scripts
- **`deploy/openshift/enable-plugin.sh`** - Script to enable/disable/check plugin status
- **`deploy/openshift/deploy.sh`** - Updated with plugin deployment actions

### Documentation

#### 8. Comprehensive Guides
- **`deploy/openshift/CONSOLE-PLUGIN-GUIDE.md`** - Complete guide to console plugin deployment
- **`deploy/openshift/QUICKSTART-CONSOLE-PLUGIN.md`** - 5-minute quick start guide
- **`CONSOLE-PLUGIN-SUMMARY.md`** - This file

#### 9. Updated Documentation
- **`README.md`** - Added console plugin deployment section
- **`deploy/openshift/README.md`** - Added deployment modes section

### Build System Updates

#### 10. Makefile
- Added `frontend-build-plugin` target
- Updated `openshift-build` to build plugin image
- Updated `openshift-push` to push plugin image

## Architecture

### Deployment Architecture

```
┌─────────────────────────────────────────┐
│    OpenShift Console (console.apps)     │
│                                         │
│  User navigates to:                     │
│  Observe → OVN Trace Viewer             │
└────────────┬────────────────────────────┘
             │ Loads plugin from
             ▼
┌─────────────────────────────────────────┐
│  Console Plugin Service (:9443)         │
│  - Serves plugin manifest               │
│  - Serves JavaScript bundles            │
│  - Serves CSS and assets                │
└────────────┬────────────────────────────┘
             │ Plugin makes API calls (proxied)
             ▼
┌─────────────────────────────────────────┐
│  Backend API Service (:8080)            │
│  - /api/trace endpoint                  │
│  - Execute ovn-trace via pods           │
│  - Parse and return results             │
└────────────┬────────────────────────────┘
             │ Exec into OVN pods
             ▼
┌─────────────────────────────────────────┐
│  OVN Pods (openshift-ovn-kubernetes)    │
│  - ovnkube-master                       │
│  - ovnkube-node                         │
└─────────────────────────────────────────┘
```

### Component Breakdown

1. **Console Integration Layer**
   - ConsolePlugin CR registers the plugin
   - Console Operator loads the plugin
   - Plugin appears in navigation

2. **Plugin Service**
   - Nginx serving plugin static assets
   - Port 9443 (standard for console plugins)
   - CORS enabled for console access

3. **Backend Service** (unchanged)
   - REST API for trace operations
   - ServiceAccount with elevated permissions
   - Executes ovn-trace in OVN pods

## Key Features

### ✅ Implemented

- [x] Dynamic plugin SDK integration
- [x] Console navigation entry (Observe menu)
- [x] Plugin deployment manifests
- [x] Service with CORS configuration
- [x] ConsolePlugin custom resource
- [x] Backend API proxy configuration
- [x] Build system (webpack)
- [x] Container image for plugin
- [x] Enable/disable plugin script
- [x] Comprehensive documentation
- [x] Both standalone and plugin modes

### 🚧 Future Enhancements

- [ ] Deep linking from pod details pages
- [ ] Context actions on pod/service pages ("Trace from here")
- [ ] Integration with topology view
- [ ] Quick actions in pod list
- [ ] Real-time trace monitoring
- [ ] Trace history and persistence

## Deployment Modes

The application now supports two deployment modes:

### 1. Console Plugin Mode (Recommended)

**Pros:**
- Native console integration
- Shared authentication/authorization
- Consistent user experience
- No separate route needed
- Appears in standard navigation

**Cons:**
- Requires Console 4.10+
- Requires cluster admin to enable

**Access:**
`OpenShift Console → Observe → OVN Trace Viewer`

### 2. Standalone Mode

**Pros:**
- Works on any Kubernetes/OpenShift
- No console integration required
- Independent deployment
- Dedicated route

**Cons:**
- Separate authentication
- Different URL
- Not integrated into console workflow

**Access:**
`https://ovnkube-viewer-<namespace>.apps.<cluster-domain>`

## Quick Start Summary

### Deploy Console Plugin

```bash
# 1. Build images
export REGISTRY=quay.io/yourorg
make openshift-build REGISTRY=$REGISTRY TAG=dev

# 2. Push images
make openshift-push REGISTRY=$REGISTRY TAG=dev

# 3. Deploy
make openshift-deploy OVERLAY=dev

# 4. Enable plugin
./deploy/openshift/enable-plugin.sh enable

# 5. Access via console
# Navigate to: Observe → OVN Trace Viewer
```

### Deploy Standalone

```bash
# Deploy standalone mode (no console integration)
./deploy/openshift/deploy.sh deploy -o dev -r quay.io/yourorg -t dev

# Get route URL
oc get route -n ovnkube-viewer-dev
```

## File Checklist

Use this checklist to verify all console plugin files are present:

### Frontend
- [ ] `frontend/console-extensions.json`
- [ ] `frontend/webpack.config.ts`
- [ ] `frontend/Dockerfile.plugin`
- [ ] `frontend/nginx.plugin.conf`
- [ ] `frontend/src/plugin/OVNTraceViewerPage.tsx`
- [ ] `frontend/src/plugin/plugin.ts`
- [ ] `frontend/package.json` (updated)

### OpenShift Manifests
- [ ] `deploy/openshift/base/console-plugin.yaml`
- [ ] `deploy/openshift/base/plugin-deployment.yaml`
- [ ] `deploy/openshift/base/plugin-service.yaml`
- [ ] `deploy/openshift/base/kustomization.yaml` (updated)
- [ ] `deploy/openshift/overlays/dev/kustomization.yaml` (updated)
- [ ] `deploy/openshift/overlays/prod/kustomization.yaml` (updated)

### Scripts
- [ ] `deploy/openshift/enable-plugin.sh`
- [ ] `deploy/openshift/deploy.sh` (updated)

### Documentation
- [ ] `deploy/openshift/CONSOLE-PLUGIN-GUIDE.md`
- [ ] `deploy/openshift/QUICKSTART-CONSOLE-PLUGIN.md`
- [ ] `deploy/openshift/README.md` (updated)
- [ ] `README.md` (updated)
- [ ] `CONSOLE-PLUGIN-SUMMARY.md` (this file)

### Build System
- [ ] `Makefile` (updated with plugin targets)

## npm Commands

New npm scripts added to `frontend/package.json`:

```bash
# Build standalone version (Vite)
npm run build:standalone

# Build console plugin (Webpack)
npm run build:plugin

# Serve plugin locally for development
npm run http-server
```

## Make Commands

New and updated Make targets:

```bash
# Build console plugin (included in openshift-build)
make frontend-build-plugin

# Build all OpenShift images (backend, frontend, plugin)
make openshift-build REGISTRY=quay.io/yourorg TAG=latest

# Push all images
make openshift-push REGISTRY=quay.io/yourorg TAG=latest

# Deploy (includes plugin)
make openshift-deploy OVERLAY=dev

# Complete workflow
make openshift-all OVERLAY=dev REGISTRY=quay.io/yourorg TAG=dev
```

## Enable/Disable Plugin

```bash
# Enable plugin in console
./deploy/openshift/enable-plugin.sh enable

# Disable plugin
./deploy/openshift/enable-plugin.sh disable

# Check plugin status
./deploy/openshift/enable-plugin.sh status
```

## Permissions

The console plugin uses the same ServiceAccount and permissions as the standalone mode:

- **ClusterRole**: `ovnkube-viewer`
  - pods, pods/exec, pods/log
  - services, endpoints
  - namespaces, nodes
  - networkpolicies
  - OpenShift network resources (egressnetworkpolicies, hostsubnets, netnamespaces)
  - routes

These permissions are required for:
1. Discovering OVN pods
2. Executing ovn-trace commands
3. Mapping pod/service names to IPs
4. Understanding network policies

## Browser Compatibility

The console plugin works in the same browsers as OpenShift Console:

- Chrome/Chromium 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Testing

### Verify Plugin Deployment

```bash
# Check all components
./deploy/openshift/enable-plugin.sh status

# Should show:
# ✓ ConsolePlugin CR exists
# ✓ Plugin is enabled in console
# ✓ Plugin service exists
# ✓ Plugin pods running: 2
```

### Test Plugin Functionality

1. Navigate to Console → Observe → OVN Trace Viewer
2. Enter source pod/IP
3. Enter destination pod/service/IP
4. Select TCP protocol
5. Click "Trace"
6. Verify trace results display

### Test Backend Integration

```bash
# From plugin pod
PLUGIN_POD=$(oc get pods -n ovnkube-viewer -l component=console-plugin -o jsonpath='{.items[0].metadata.name}')

# Test backend health
oc exec -n ovnkube-viewer $PLUGIN_POD -- \
  wget -qO- http://ovnkube-viewer-backend:8080/health
```

## Troubleshooting Quick Reference

| Issue | Command | Fix |
|-------|---------|-----|
| Plugin not in menu | `./enable-plugin.sh status` | Enable plugin |
| Console not loading | `oc logs -n openshift-console -l app=console` | Check console logs |
| Plugin 404 | `oc get svc ovnkube-viewer-plugin -n ovnkube-viewer` | Check service exists |
| Backend errors | `oc logs -n ovnkube-viewer -l component=backend` | Check backend logs |
| CORS errors | Check browser console | Verify nginx CORS config |

## Resources

- **OpenShift Docs**: [Console Plugins](https://docs.openshift.com/container-platform/latest/web_console/dynamic-plugin/overview-dynamic-plugin.html)
- **Template**: [console-plugin-template](https://github.com/openshift/console-plugin-template)
- **SDK**: [dynamic-plugin-sdk](https://github.com/openshift/console/tree/master/frontend/packages/console-dynamic-plugin-sdk)
- **Blog**: [Developing Console Plugins](https://www.redhat.com/en/blog/developing-an-openshift-dynamic-console-plugin-1)

## What's Next?

Potential enhancements for future releases:

1. **Context Integration**
   - Right-click trace from pod details
   - Quick action buttons on pod lists
   - Integration with topology view

2. **Advanced Features**
   - Trace history and persistence
   - Scheduled/automated traces
   - Trace comparison tool
   - Export traces to various formats

3. **Observability**
   - Prometheus metrics for trace operations
   - Grafana dashboards
   - Alerting on trace failures

4. **Performance**
   - Trace result caching
   - Parallel trace execution
   - Batch trace operations

## Credits

This implementation follows the OpenShift Console Plugin architecture and uses:
- OpenShift Dynamic Plugin SDK
- PatternFly React components
- Webpack for plugin bundling
- Nginx for serving plugin assets

## License

Same as main project (TBD).

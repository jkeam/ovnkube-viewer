# Documentation

This directory contains all documentation for the OVN-Kubernetes Trace Viewer project.

## Quick Start

- **[Getting Started Guide](getting-started.md)** - Setup and development instructions
- **[BUILD-QUICK-REFERENCE.md](BUILD-QUICK-REFERENCE.md)** - Quick commands to build containers

## Building & Deployment

- **[Building Containers](building-containers.md)** - Complete guide for building all container images
- **[UBI Base Images](ubi-base-images.md)** - Red Hat UBI migration guide and details
- **[Run Script Guide](run-script-guide.md)** - run.sh command reference

## OpenShift Deployment

See `../deploy/openshift/` directory for OpenShift-specific documentation:
- [OpenShift README](../deploy/openshift/README.md) - Architecture and configuration
- [Deployment Guide](../deploy/openshift/DEPLOYMENT-GUIDE.md) - Step-by-step deployment
- [Console Plugin Guide](../deploy/openshift/CONSOLE-PLUGIN-GUIDE.md) - OpenShift Console integration
- [Console Plugin Quick Start](../deploy/openshift/QUICKSTART-CONSOLE-PLUGIN.md) - 5-minute quick start

## Design & Planning

- **[Implementation Plan](implementation-plan.md)** - Comprehensive development roadmap
- **[PatternFly Design System](patternfly.md)** - Design system implementation guide

## Reference

- **[Node.js Upgrade Guide](node-upgrade.md)** - How to upgrade to Node.js 24 LTS
- **[CHANGELOG](CHANGELOG.md)** - Project change log

## Project Summaries

- **[Console Plugin Summary](../CONSOLE-PLUGIN-SUMMARY.md)** - Console plugin implementation details
- **[UBI Migration Summary](UBI-MIGRATION-SUMMARY.md)** - Red Hat UBI migration overview

## Documentation Organization

```
ovnkube-viewer/
├── docs/                              # General documentation
│   ├── README.md                      # This file
│   ├── BUILD-QUICK-REFERENCE.md       # ⭐ Quick build commands
│   ├── building-containers.md         # Complete build guide
│   ├── ubi-base-images.md            # UBI migration details
│   ├── UBI-MIGRATION-SUMMARY.md       # UBI migration overview
│   ├── getting-started.md            # Setup guide
│   ├── implementation-plan.md        # Development roadmap
│   ├── patternfly.md                 # Design system
│   ├── run-script-guide.md           # run.sh reference
│   ├── node-upgrade.md               # Node.js upgrade
│   └── CHANGELOG.md                  # Changes
│
├── deploy/openshift/                  # OpenShift deployment docs
│   ├── README.md                      # OpenShift overview
│   ├── DEPLOYMENT-GUIDE.md            # Deployment steps
│   ├── CONSOLE-PLUGIN-GUIDE.md        # Console plugin guide
│   └── QUICKSTART-CONSOLE-PLUGIN.md   # Quick start
│
├── README.md                          # Main project README
└── CONSOLE-PLUGIN-SUMMARY.md          # Plugin implementation
```

## Common Tasks

### Building Images

```bash
# Quick reference
see BUILD-QUICK-REFERENCE.md

# Detailed guide
see building-containers.md
```

### OpenShift Deployment

```bash
# Quick start
see deploy/openshift/QUICKSTART-CONSOLE-PLUGIN.md

# Complete guide
see deploy/openshift/DEPLOYMENT-GUIDE.md
```

### Development Setup

```bash
# Initial setup
see getting-started.md

# Run script
see run-script-guide.md
```

## Contributing to Documentation

When adding documentation:

1. **Place in appropriate directory:**
   - General docs → `docs/`
   - OpenShift specific → `deploy/openshift/`
   - Project summaries → root

2. **Update this README** to include new docs

3. **Use relative links** for internal references

4. **Follow naming conventions:**
   - Lowercase with hyphens: `my-doc-name.md`
   - CAPS for important files: `README.md`, `CHANGELOG.md`, `BUILD-QUICK-REFERENCE.md`
   - Summaries can use CAPS if needed

## Getting Help

- Check the relevant guide above
- See [main README](../README.md) for project overview
- Open an issue at [GitHub Issues](https://github.com/anthropics/claude-code/issues)

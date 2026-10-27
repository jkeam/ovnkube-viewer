# run.sh Quick Reference Guide

A convenient wrapper script to manage the OVN-Kubernetes Trace Viewer application.

## Quick Commands

```bash
./run.sh help      # Show all available commands
./run.sh init      # Initialize project dependencies
./run.sh dev       # Run in development mode
./run.sh docker    # Run with Docker Compose
./run.sh stop      # Stop all services
./run.sh status    # Check service status
```

## Command Reference

### Development Commands

#### `./run.sh dev`
Run both backend and frontend in development mode with hot reload.

- **Backend**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **Press Ctrl+C** to stop both services
- **Best for**: Active development with hot reload

```bash
./run.sh dev
# Open http://localhost:3000
```

#### `./run.sh docker`
Run the application using Docker Compose.

- Builds Docker images if needed
- Starts containers in detached mode
- **Best for**: Production-like testing, easy setup

```bash
./run.sh docker
# Open http://localhost:3000
```

### Management Commands

#### `./run.sh init`
Initialize project dependencies.

- Runs `go mod tidy` for backend
- Runs `npm install` for frontend
- Checks for required tools (Go, Node, npm, Docker)

```bash
./run.sh init
```

#### `./run.sh stop`
Stop all running services.

- Stops Docker containers
- Kills development processes
- Safe to run even if nothing is running

```bash
./run.sh stop
```

#### `./run.sh status`
Check the status of services.

- Shows Docker container status
- Shows development process status
- Shows port availability (3000, 8080)

```bash
./run.sh status
```

### Build & Test Commands

#### `./run.sh build`
Build both backend and frontend.

- Backend binary: `backend/bin/server`
- Frontend build: `frontend/dist/`

```bash
./run.sh build
```

#### `./run.sh test`
Run all tests.

- Backend: `go test ./...`
- Frontend: `npm test` (if tests exist)

```bash
./run.sh test
```

#### `./run.sh clean`
Clean all build artifacts.

- Removes `backend/bin/`
- Removes `frontend/dist/`
- Removes `frontend/node_modules/`
- Cleans Go cache

```bash
./run.sh clean
./run.sh init  # Reinstall after cleaning
```

### Monitoring Commands

#### `./run.sh logs`
View Docker container logs (when running via Docker).

- Shows real-time logs
- Press Ctrl+C to exit

```bash
./run.sh logs
```

## Common Workflows

### First Time Setup

```bash
# 1. Initialize project
./run.sh init

# 2. Run with Docker (easiest)
./run.sh docker

# 3. Open browser
open http://localhost:3000
```

### Development Workflow

```bash
# 1. Start in dev mode (hot reload)
./run.sh dev

# 2. Make changes to code
# 3. See changes automatically reload

# 4. Stop when done
# Press Ctrl+C
```

### Testing Changes

```bash
# 1. Make code changes
# 2. Run tests
./run.sh test

# 3. Build production version
./run.sh build

# 4. Test with Docker
./run.sh docker
```

### Cleaning Up

```bash
# Stop all services
./run.sh stop

# Clean build artifacts
./run.sh clean

# Reinstall if needed
./run.sh init
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using the ports
./run.sh status

# Stop all services
./run.sh stop

# Try again
./run.sh dev  # or ./run.sh docker
```

### Dependencies Missing

```bash
# Reinitialize dependencies
./run.sh clean
./run.sh init
```

### Docker Issues

```bash
# Stop Docker containers
./run.sh stop

# Rebuild from scratch
cd deploy/standalone
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
cd ../..
```

### Development Server Issues

```bash
# Kill all related processes
./run.sh stop

# Clean and rebuild
./run.sh clean
./run.sh init
./run.sh build

# Try again
./run.sh dev
```

## Features

### ✅ What the Script Does

- **Dependency Checking**: Verifies Go, Node, npm, Docker are installed
- **Color Output**: Easy-to-read colored terminal output
- **Process Management**: Handles starting/stopping both services
- **Error Handling**: Gracefully handles errors and missing dependencies
- **Service Monitoring**: Check if services are running
- **Port Detection**: Verifies ports 3000 and 8080 availability

### 🎨 Output Colors

- 🟢 **Green**: Success messages
- 🔴 **Red**: Error messages
- 🟡 **Yellow**: Info messages
- 🔵 **Blue**: Section headers

## Examples

### Running for the First Time

```bash
$ ./run.sh init
╔════════════════════════════════════════════════════════╗
║  OVN-Kubernetes Trace Viewer - Phase 1 MVP            ║
╚════════════════════════════════════════════════════════╝

━━━ Initializing Project ━━━

✓ Go go1.21.0
✓ Node v20.0.0
✓ npm 9.6.0

➜ Installing Go dependencies...
✓ Go dependencies installed
➜ Installing Node.js dependencies...
✓ Node.js dependencies installed
✓ Project initialized successfully!
```

### Starting Development Mode

```bash
$ ./run.sh dev
╔════════════════════════════════════════════════════════╗
║  OVN-Kubernetes Trace Viewer - Phase 1 MVP            ║
╚════════════════════════════════════════════════════════╝

━━━ Running in Development Mode ━━━

➜ Backend will run on: http://localhost:8080
➜ Frontend will run on: http://localhost:3000

➜ Press Ctrl+C to stop both services

➜ Starting backend...
➜ Starting frontend...

✓ Services started!

➜ Open http://localhost:3000 in your browser
```

### Checking Status

```bash
$ ./run.sh status
╔════════════════════════════════════════════════════════╗
║  OVN-Kubernetes Trace Viewer - Phase 1 MVP            ║
╚════════════════════════════════════════════════════════╝

━━━ Service Status ━━━

Docker Containers:
NAME                          COMMAND                  SERVICE     STATUS
ovnkube-viewer-backend       "./server"               backend     Up 5 minutes
ovnkube-viewer-frontend      "nginx -g 'daemon of…"   frontend    Up 5 minutes

Port Status:
✓ Port 8080 (backend) is in use
✓ Port 3000 (frontend) is in use
```

## Tips

1. **Use `dev` for development**: Hot reload makes development faster
2. **Use `docker` for testing**: Closer to production environment
3. **Check `status` first**: Before starting, see what's already running
4. **Always `stop` before switching modes**: Prevent port conflicts
5. **Use `init` after pulling changes**: Ensure dependencies are up to date

## Related Documentation

- [Getting Started Guide](GETTING_STARTED.md) - Detailed setup instructions
- [README](README.md) - Project overview
- [Phase 1 Status](PHASE1_STATUS.md) - Development progress

---

**Pro Tip**: Add an alias to your shell for quick access:
```bash
# Add to ~/.bashrc or ~/.zshrc
alias ovn-dev='cd /path/to/ovnkube-viewer && ./run.sh dev'
alias ovn-docker='cd /path/to/ovnkube-viewer && ./run.sh docker'
alias ovn-stop='cd /path/to/ovnkube-viewer && ./run.sh stop'
```

#!/bin/bash

# OVN-Kubernetes Trace Viewer - Easy Run Script
# Usage: ./run.sh [command]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  OVN-Kubernetes Trace Viewer - Phase 1 MVP            ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
    echo
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}➜ $1${NC}"
}

print_section() {
    echo
    echo -e "${BLUE}━━━ $1 ━━━${NC}"
    echo
}

show_help() {
    print_header
    echo "Usage: ./run.sh [command]"
    echo
    echo "Commands:"
    echo "  dev         - Run in development mode (recommended for development)"
    echo "  docker      - Run with Docker Compose (easiest setup)"
    echo "  stop        - Stop all running services"
    echo "  clean       - Clean build artifacts and dependencies"
    echo "  init        - Initialize project dependencies"
    echo "  build       - Build both backend and frontend"
    echo "  test        - Run tests"
    echo "  logs        - View Docker logs (when running via docker)"
    echo "  status      - Check if services are running"
    echo "  help        - Show this help message"
    echo
    echo "Examples:"
    echo "  ./run.sh docker     # Start with Docker Compose"
    echo "  ./run.sh dev        # Start in development mode"
    echo "  ./run.sh stop       # Stop all services"
    echo
}

check_dependencies() {
    print_section "Checking Dependencies"

    local missing_deps=0

    # Check for required tools
    if command -v go &> /dev/null; then
        print_success "Go $(go version | awk '{print $3}')"
    else
        print_error "Go is not installed"
        missing_deps=1
    fi

    if command -v node &> /dev/null; then
        NODE_VERSION=$(node --version)
        print_success "Node $NODE_VERSION"

        # Check if version is at least 24.22.1
        NODE_MAJOR=$(echo $NODE_VERSION | sed 's/v\([0-9]*\).*/\1/')
        if [ "$NODE_MAJOR" -lt 24 ]; then
            print_error "Node.js version must be 24.22.1 or later (found $NODE_VERSION)"
            echo "  Install Node.js 24 LTS: https://nodejs.org/"
            missing_deps=1
        fi
    else
        print_error "Node.js is not installed"
        missing_deps=1
    fi

    if command -v npm &> /dev/null; then
        print_success "npm $(npm --version)"
    else
        print_error "npm is not installed"
        missing_deps=1
    fi

    if [ "$1" == "docker" ]; then
        if command -v docker &> /dev/null; then
            print_success "Docker $(docker --version | awk '{print $3}' | tr -d ',')"
        else
            print_error "Docker is not installed"
            missing_deps=1
        fi

        if command -v docker-compose &> /dev/null; then
            print_success "Docker Compose $(docker-compose --version | awk '{print $4}' | tr -d ',')"
        else
            print_error "Docker Compose is not installed"
            missing_deps=1
        fi
    fi

    echo

    if [ $missing_deps -eq 1 ]; then
        print_error "Missing required dependencies. Please install them first."
        exit 1
    fi
}

init_project() {
    print_header
    print_section "Initializing Project"

    check_dependencies

    print_info "Installing Go dependencies..."
    cd backend && go mod tidy && cd ..
    print_success "Go dependencies installed"

    print_info "Installing Node.js dependencies..."
    cd frontend && npm install && cd ..
    print_success "Node.js dependencies installed"

    print_success "Project initialized successfully!"
    echo
    print_info "You can now run:"
    echo "  ./run.sh docker    - To start with Docker"
    echo "  ./run.sh dev       - To start in development mode"
    echo
}

run_dev() {
    print_header
    print_section "Running in Development Mode"

    check_dependencies

    # Check if dependencies are installed
    if [ ! -f "backend/go.mod" ] || [ ! -d "frontend/node_modules" ]; then
        print_error "Dependencies not installed. Running init first..."
        init_project
    fi

    print_info "Starting Backend and Frontend..."
    echo
    print_info "Backend will run on: http://localhost:8080"
    print_info "Frontend will run on: http://localhost:3000"
    echo
    print_info "Press Ctrl+C to stop both services"
    echo

    # Function to cleanup background processes
    cleanup() {
        echo
        print_info "Shutting down services..."
        kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
        print_success "Services stopped"
        exit 0
    }

    trap cleanup SIGINT SIGTERM

    # Start backend
    print_info "Starting backend..."
    (cd backend && go run ./cmd/server) &
    BACKEND_PID=$!

    # Wait a bit for backend to start
    sleep 2

    # Start frontend
    print_info "Starting frontend..."
    (cd frontend && npm run dev) &
    FRONTEND_PID=$!

    echo
    print_success "Services started!"
    echo
    print_info "Open http://localhost:3000 in your browser"
    echo

    # Wait for processes
    wait $BACKEND_PID $FRONTEND_PID
}

run_docker() {
    print_header
    print_section "Running with Docker Compose"

    check_dependencies docker

    print_info "Building Docker images..."
    cd deploy/standalone && docker-compose build

    print_info "Starting containers..."
    docker-compose up -d

    echo
    print_success "Services started!"
    echo
    print_info "Frontend: http://localhost:3000"
    print_info "Backend:  http://localhost:8080"
    echo
    print_info "View logs with: ./run.sh logs"
    print_info "Stop with:      ./run.sh stop"
    echo
}

stop_services() {
    print_header
    print_section "Stopping Services"

    # Stop Docker containers if running
    if [ -f "deploy/standalone/docker-compose.yaml" ]; then
        cd deploy/standalone
        if docker-compose ps -q &> /dev/null; then
            print_info "Stopping Docker containers..."
            docker-compose down
            print_success "Docker containers stopped"
        fi
        cd ../..
    fi

    # Kill any running Go processes (backend)
    if pgrep -f "go run ./cmd/server" > /dev/null; then
        print_info "Stopping backend..."
        pkill -f "go run ./cmd/server"
        print_success "Backend stopped"
    fi

    # Kill any running Node processes (frontend)
    if pgrep -f "vite" > /dev/null; then
        print_info "Stopping frontend..."
        pkill -f "vite"
        print_success "Frontend stopped"
    fi

    print_success "All services stopped"
    echo
}

clean_project() {
    print_header
    print_section "Cleaning Project"

    print_info "Removing build artifacts..."
    rm -rf backend/bin
    rm -rf frontend/dist
    rm -rf frontend/node_modules

    print_info "Cleaning Go cache..."
    cd backend && go clean && cd ..

    print_success "Project cleaned"
    echo
    print_info "Run './run.sh init' to reinstall dependencies"
    echo
}

build_project() {
    print_header
    print_section "Building Project"

    print_info "Building backend..."
    cd backend && go build -o bin/server ./cmd/server && cd ..
    print_success "Backend built: backend/bin/server"

    print_info "Building frontend..."
    cd frontend && npm run build && cd ..
    print_success "Frontend built: frontend/dist/"

    echo
    print_success "Build complete!"
    echo
}

run_tests() {
    print_header
    print_section "Running Tests"

    print_info "Running backend tests..."
    cd backend && go test ./... -v && cd ..

    # Frontend tests (if they exist)
    if [ -f "frontend/package.json" ] && grep -q "\"test\"" frontend/package.json; then
        print_info "Running frontend tests..."
        cd frontend && npm test && cd ..
    fi

    print_success "Tests complete!"
    echo
}

show_logs() {
    print_header
    print_section "Docker Logs"

    if [ -f "deploy/standalone/docker-compose.yaml" ]; then
        cd deploy/standalone
        docker-compose logs -f
    else
        print_error "Docker Compose file not found"
        exit 1
    fi
}

check_status() {
    print_header
    print_section "Service Status"

    # Check Docker containers
    if command -v docker &> /dev/null; then
        echo "Docker Containers:"
        if [ -f "deploy/standalone/docker-compose.yaml" ]; then
            cd deploy/standalone
            docker-compose ps
            cd ../..
        else
            echo "  No docker-compose.yaml found"
        fi
        echo
    fi

    # Check backend
    if pgrep -f "go run ./cmd/server" > /dev/null; then
        print_success "Backend (dev) is running"
    else
        print_info "Backend (dev) is not running"
    fi

    # Check frontend
    if pgrep -f "vite" > /dev/null; then
        print_success "Frontend (dev) is running"
    else
        print_info "Frontend (dev) is not running"
    fi

    # Check if ports are in use
    echo
    echo "Port Status:"
    if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
        print_success "Port 8080 (backend) is in use"
    else
        print_info "Port 8080 (backend) is available"
    fi

    if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null 2>&1; then
        print_success "Port 3000 (frontend) is in use"
    else
        print_info "Port 3000 (frontend) is available"
    fi

    echo
}

# Main script logic
case "${1:-help}" in
    dev)
        run_dev
        ;;
    docker)
        run_docker
        ;;
    stop)
        stop_services
        ;;
    clean)
        clean_project
        ;;
    init)
        init_project
        ;;
    build)
        build_project
        ;;
    test)
        run_tests
        ;;
    logs)
        show_logs
        ;;
    status)
        check_status
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo
        show_help
        exit 1
        ;;
esac

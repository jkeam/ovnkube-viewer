#!/bin/bash
set -e

# OpenShift Deployment Helper Script
# This script helps deploy the OVN-Kubernetes Trace Viewer to OpenShift

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Default values
OVERLAY="dev"
REGISTRY="quay.io/yourorg"
TAG="latest"
ACTION=""
SKIP_BUILD=false
SKIP_PUSH=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_usage() {
    cat << EOF
Usage: $0 [ACTION] [OPTIONS]

ACTIONS:
    deploy          Build, push, and deploy to OpenShift
    build           Build container images only
    push            Push container images only
    validate        Validate Kustomize manifests
    apply           Apply manifests to OpenShift (skip build/push)
    status          Check deployment status
    logs            Tail backend logs
    delete          Delete deployment from OpenShift
    enable-plugin   Enable the console plugin
    disable-plugin  Disable the console plugin
    plugin-status   Check console plugin status
    help            Show this help message

OPTIONS:
    -o, --overlay OVERLAY       Overlay to use (dev|prod) [default: dev]
    -r, --registry REGISTRY     Container registry [default: quay.io/yourorg]
    -t, --tag TAG              Image tag [default: latest]
    --skip-build               Skip building images
    --skip-push                Skip pushing images

EXAMPLES:
    # Deploy to development
    $0 deploy -o dev -r quay.io/myorg -t v1.0.0

    # Deploy to production
    $0 deploy -o prod -r quay.io/myorg -t v1.0.0

    # Just apply manifests (images already pushed)
    $0 apply -o dev

    # Check status
    $0 status -o prod

    # View logs
    $0 logs -o dev

EOF
}

check_prerequisites() {
    log_info "Checking prerequisites..."

    if ! command -v oc &> /dev/null; then
        log_error "oc CLI not found. Please install OpenShift CLI."
        exit 1
    fi

    if ! command -v podman &> /dev/null && ! command -v docker &> /dev/null; then
        log_error "Neither podman nor docker found. Please install one of them."
        exit 1
    fi

    # Check if logged into OpenShift
    if ! oc whoami &> /dev/null; then
        log_error "Not logged into OpenShift. Please run 'oc login' first."
        exit 1
    fi

    log_info "Prerequisites check passed ✓"
}

build_images() {
    log_info "Building container images..."

    local builder="podman"
    if ! command -v podman &> /dev/null; then
        builder="docker"
    fi

    log_info "Building backend image: ${REGISTRY}/ovnkube-viewer-backend:${TAG}"
    $builder build -t "${REGISTRY}/ovnkube-viewer-backend:${TAG}" \
        -f "${PROJECT_ROOT}/backend/Dockerfile" \
        "${PROJECT_ROOT}/backend"

    log_info "Building frontend image: ${REGISTRY}/ovnkube-viewer-frontend:${TAG}"
    $builder build -t "${REGISTRY}/ovnkube-viewer-frontend:${TAG}" \
        -f "${PROJECT_ROOT}/frontend/Dockerfile" \
        "${PROJECT_ROOT}/frontend"

    log_info "Images built successfully ✓"
}

push_images() {
    log_info "Pushing container images..."

    local builder="podman"
    if ! command -v podman &> /dev/null; then
        builder="docker"
    fi

    log_info "Pushing ${REGISTRY}/ovnkube-viewer-backend:${TAG}"
    $builder push "${REGISTRY}/ovnkube-viewer-backend:${TAG}"

    log_info "Pushing ${REGISTRY}/ovnkube-viewer-frontend:${TAG}"
    $builder push "${REGISTRY}/ovnkube-viewer-frontend:${TAG}"

    log_info "Images pushed successfully ✓"
}

validate_manifests() {
    log_info "Validating Kustomize manifests for overlay: ${OVERLAY}"

    oc kustomize "${SCRIPT_DIR}/overlays/${OVERLAY}"

    log_info "Manifests validated successfully ✓"
}

apply_manifests() {
    log_info "Applying manifests for overlay: ${OVERLAY}"

    oc apply -k "${SCRIPT_DIR}/overlays/${OVERLAY}"

    log_info "Waiting for deployments to be ready..."
    sleep 5

    if [ "$OVERLAY" = "dev" ]; then
        NAMESPACE="ovnkube-viewer-dev"
        BACKEND_DEPLOYMENT="dev-ovnkube-viewer-backend"
        FRONTEND_DEPLOYMENT="dev-ovnkube-viewer-frontend"
    else
        NAMESPACE="ovnkube-viewer"
        BACKEND_DEPLOYMENT="ovnkube-viewer-backend"
        FRONTEND_DEPLOYMENT="ovnkube-viewer-frontend"
    fi

    oc wait --for=condition=available --timeout=300s \
        deployment/${BACKEND_DEPLOYMENT} -n ${NAMESPACE} || true

    oc wait --for=condition=available --timeout=300s \
        deployment/${FRONTEND_DEPLOYMENT} -n ${NAMESPACE} || true

    log_info "Deployment complete ✓"
    log_info "Getting route URL..."
    oc get route -n ${NAMESPACE}
}

show_status() {
    log_info "Checking deployment status for overlay: ${OVERLAY}"

    if [ "$OVERLAY" = "dev" ]; then
        NAMESPACE="ovnkube-viewer-dev"
    else
        NAMESPACE="ovnkube-viewer"
    fi

    echo ""
    echo "=== Pods ==="
    oc get pods -n ${NAMESPACE} -l app=ovnkube-viewer

    echo ""
    echo "=== Deployments ==="
    oc get deployments -n ${NAMESPACE} -l app=ovnkube-viewer

    echo ""
    echo "=== Services ==="
    oc get services -n ${NAMESPACE} -l app=ovnkube-viewer

    echo ""
    echo "=== Routes ==="
    oc get routes -n ${NAMESPACE}
}

tail_logs() {
    log_info "Tailing backend logs for overlay: ${OVERLAY}"

    if [ "$OVERLAY" = "dev" ]; then
        NAMESPACE="ovnkube-viewer-dev"
        DEPLOYMENT="dev-ovnkube-viewer-backend"
    else
        NAMESPACE="ovnkube-viewer"
        DEPLOYMENT="ovnkube-viewer-backend"
    fi

    oc logs -f deployment/${DEPLOYMENT} -n ${NAMESPACE}
}

delete_deployment() {
    log_warn "Deleting deployment for overlay: ${OVERLAY}"
    read -p "Are you sure? (yes/no): " confirm

    if [ "$confirm" != "yes" ]; then
        log_info "Deletion cancelled"
        exit 0
    fi

    oc delete -k "${SCRIPT_DIR}/overlays/${OVERLAY}"

    log_info "Deployment deleted ✓"
}

enable_plugin() {
    log_info "Enabling console plugin..."
    "${SCRIPT_DIR}/enable-plugin.sh" enable
}

disable_plugin() {
    log_info "Disabling console plugin..."
    "${SCRIPT_DIR}/enable-plugin.sh" disable
}

plugin_status() {
    "${SCRIPT_DIR}/enable-plugin.sh" status
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        deploy|build|push|validate|apply|status|logs|delete|enable-plugin|disable-plugin|plugin-status|help)
            ACTION="$1"
            shift
            ;;
        -o|--overlay)
            OVERLAY="$2"
            shift 2
            ;;
        -r|--registry)
            REGISTRY="$2"
            shift 2
            ;;
        -t|--tag)
            TAG="$2"
            shift 2
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --skip-push)
            SKIP_PUSH=true
            shift
            ;;
        *)
            log_error "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Validate overlay
if [ "$OVERLAY" != "dev" ] && [ "$OVERLAY" != "prod" ]; then
    log_error "Invalid overlay: ${OVERLAY}. Must be 'dev' or 'prod'."
    exit 1
fi

# Execute action
case $ACTION in
    deploy)
        check_prerequisites
        if [ "$SKIP_BUILD" != true ]; then
            build_images
        fi
        if [ "$SKIP_PUSH" != true ]; then
            push_images
        fi
        apply_manifests
        ;;
    build)
        check_prerequisites
        build_images
        ;;
    push)
        check_prerequisites
        push_images
        ;;
    validate)
        check_prerequisites
        validate_manifests
        ;;
    apply)
        check_prerequisites
        apply_manifests
        ;;
    status)
        check_prerequisites
        show_status
        ;;
    logs)
        check_prerequisites
        tail_logs
        ;;
    delete)
        check_prerequisites
        delete_deployment
        ;;
    enable-plugin)
        check_prerequisites
        enable_plugin
        ;;
    disable-plugin)
        check_prerequisites
        disable_plugin
        ;;
    plugin-status)
        check_prerequisites
        plugin_status
        ;;
    help|"")
        print_usage
        exit 0
        ;;
    *)
        log_error "Unknown action: ${ACTION}"
        print_usage
        exit 1
        ;;
esac

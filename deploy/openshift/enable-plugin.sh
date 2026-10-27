#!/bin/bash
set -e

# OpenShift Console Plugin Enable/Disable Script

PLUGIN_NAME="ovnkube-viewer"
ACTION="${1:-enable}"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_usage() {
    cat << EOF
Usage: $0 [enable|disable|status]

Manage the OVN Trace Viewer console plugin in OpenShift.

ACTIONS:
    enable      Enable the console plugin (default)
    disable     Disable the console plugin
    status      Show plugin status

EXAMPLES:
    # Enable the plugin
    $0 enable

    # Disable the plugin
    $0 disable

    # Check status
    $0 status

EOF
}

check_prerequisites() {
    if ! command -v oc &> /dev/null; then
        log_error "oc CLI not found. Please install OpenShift CLI."
        exit 1
    fi

    if ! oc whoami &> /dev/null; then
        log_error "Not logged into OpenShift. Please run 'oc login' first."
        exit 1
    fi

    # Check if user is cluster admin
    if ! oc auth can-i patch consoles.operator.openshift.io/cluster &> /dev/null; then
        log_error "Insufficient permissions. Cluster admin access required."
        exit 1
    fi
}

get_current_plugins() {
    oc get consoles.operator.openshift.io cluster \
        -o jsonpath='{.spec.plugins}' 2>/dev/null || echo "[]"
}

is_plugin_enabled() {
    local plugins=$(get_current_plugins)
    if echo "$plugins" | grep -q "\"$PLUGIN_NAME\""; then
        return 0
    else
        return 1
    fi
}

show_status() {
    log_info "Checking plugin status..."
    echo ""

    # Check if ConsolePlugin CR exists
    if oc get consoleplugin $PLUGIN_NAME &> /dev/null; then
        echo "✓ ConsolePlugin CR exists"
    else
        echo "✗ ConsolePlugin CR not found"
        log_warn "Plugin not deployed. Run: oc apply -k deploy/openshift/overlays/dev"
        return 1
    fi

    # Check if plugin is enabled
    if is_plugin_enabled; then
        echo "✓ Plugin is enabled in console"
    else
        echo "✗ Plugin is not enabled in console"
    fi

    # Check plugin service
    if oc get svc ovnkube-viewer-plugin -n ovnkube-viewer &> /dev/null; then
        echo "✓ Plugin service exists"
    else
        echo "✗ Plugin service not found"
    fi

    # Check plugin pods
    local pod_count=$(oc get pods -n ovnkube-viewer -l component=console-plugin --no-headers 2>/dev/null | wc -l | tr -d ' ')
    if [ "$pod_count" -gt 0 ]; then
        echo "✓ Plugin pods running: $pod_count"
        oc get pods -n ovnkube-viewer -l component=console-plugin
    else
        echo "✗ No plugin pods found"
    fi

    echo ""
    log_info "Current enabled plugins:"
    get_current_plugins | jq -r '.[]' 2>/dev/null || echo "  (none)"
}

enable_plugin() {
    log_info "Enabling console plugin: $PLUGIN_NAME"

    # Check if already enabled
    if is_plugin_enabled; then
        log_warn "Plugin is already enabled"
        return 0
    fi

    # Check if ConsolePlugin CR exists
    if ! oc get consoleplugin $PLUGIN_NAME &> /dev/null; then
        log_error "ConsolePlugin CR '$PLUGIN_NAME' not found."
        log_info "Please deploy the plugin first:"
        log_info "  oc apply -k deploy/openshift/overlays/dev"
        exit 1
    fi

    # Get current plugins
    local current_plugins=$(get_current_plugins)

    # Add plugin to the list
    if [ "$current_plugins" = "[]" ] || [ -z "$current_plugins" ]; then
        # No plugins enabled yet
        log_info "No plugins currently enabled. Initializing plugins list..."
        oc patch consoles.operator.openshift.io cluster \
            --type=json \
            --patch "[{\"op\": \"add\", \"path\": \"/spec/plugins\", \"value\": [\"$PLUGIN_NAME\"]}]"
    else
        # Plugins already exist, append to list
        log_info "Adding to existing plugins list..."
        oc patch consoles.operator.openshift.io cluster \
            --type=json \
            --patch "[{\"op\": \"add\", \"path\": \"/spec/plugins/-\", \"value\": \"$PLUGIN_NAME\"}]"
    fi

    log_info "Plugin enabled successfully ✓"
    log_info "Console pods will restart to load the plugin..."
    log_info "Waiting for console pods to restart..."

    sleep 5
    oc rollout status deployment/console -n openshift-console --timeout=300s || true

    echo ""
    log_info "Plugin enabled! Access it via:"
    log_info "  OpenShift Console → Observe → OVN Trace Viewer"
}

disable_plugin() {
    log_info "Disabling console plugin: $PLUGIN_NAME"

    # Check if enabled
    if ! is_plugin_enabled; then
        log_warn "Plugin is not enabled"
        return 0
    fi

    # Get current plugins as JSON array
    local plugins=$(get_current_plugins)

    # Find the index of our plugin
    local index=$(echo "$plugins" | jq -r "to_entries | .[] | select(.value==\"$PLUGIN_NAME\") | .key")

    if [ -z "$index" ]; then
        log_warn "Plugin not found in enabled list"
        return 0
    fi

    log_info "Removing plugin from enabled list (index: $index)..."
    oc patch consoles.operator.openshift.io cluster \
        --type=json \
        --patch "[{\"op\": \"remove\", \"path\": \"/spec/plugins/$index\"}]"

    log_info "Plugin disabled successfully ✓"
    log_info "Console pods will restart..."
    log_info "Waiting for console pods to restart..."

    sleep 5
    oc rollout status deployment/console -n openshift-console --timeout=300s || true

    echo ""
    log_info "Plugin disabled!"
}

# Main execution
case "$ACTION" in
    enable)
        check_prerequisites
        enable_plugin
        ;;
    disable)
        check_prerequisites
        disable_plugin
        ;;
    status)
        check_prerequisites
        show_status
        ;;
    help|--help|-h)
        print_usage
        exit 0
        ;;
    *)
        log_error "Unknown action: $ACTION"
        print_usage
        exit 1
        ;;
esac

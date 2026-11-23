package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jkeam/ovnkube-viewer/pkg/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// K8sTraceExecutor executes ovn-trace commands in Kubernetes OVN pods
type K8sTraceExecutor struct {
	clientset      *kubernetes.Clientset
	config         *rest.Config
	ovnNamespace   string
	masterPodLabel string
	nodePodLabel   string
}

// NewK8sTraceExecutor creates a new K8sTraceExecutor
func NewK8sTraceExecutor() (*K8sTraceExecutor, error) {
	// Try to load in-cluster config first
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = os.Getenv("HOME") + "/.kube/config"
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load kubernetes config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Get OVN configuration from environment variables
	ovnNamespace := os.Getenv("OVN_NAMESPACE")
	if ovnNamespace == "" {
		ovnNamespace = "openshift-ovn-kubernetes"
	}

	masterPodLabel := os.Getenv("OVN_MASTER_POD_LABEL")
	if masterPodLabel == "" {
		masterPodLabel = "app=ovnkube-master"
	}

	nodePodLabel := os.Getenv("OVN_NODE_POD_LABEL")
	if nodePodLabel == "" {
		nodePodLabel = "app=ovnkube-node"
	}

	return &K8sTraceExecutor{
		clientset:      clientset,
		config:         config,
		ovnNamespace:   ovnNamespace,
		masterPodLabel: masterPodLabel,
		nodePodLabel:   nodePodLabel,
	}, nil
}

// Execute runs an ovn-trace command in an OVN master pod
func (e *K8sTraceExecutor) Execute(req *models.TraceRequest) (string, error) {
	// Validate request
	if req.Source == "" || req.Destination == "" {
		return "", fmt.Errorf("source and destination are required")
	}

	// Resolve source and destination to IP addresses if needed
	sourceIP, err := e.resolveToIP(req.Source)
	if err != nil {
		return "", fmt.Errorf("failed to resolve source: %w", err)
	}
	destIP, err := e.resolveToIP(req.Destination)
	if err != nil {
		return "", fmt.Errorf("failed to resolve destination: %w", err)
	}

	fmt.Printf("INFO: Resolved source %s -> %s, destination %s -> %s\n", req.Source, sourceIP, req.Destination, destIP)

	// Find an OVN node pod (which has access to OVN databases)
	pod, err := e.findOVNNodePod()
	if err != nil {
		fmt.Printf("ERROR: Failed to find OVN node pod: %v\n", err)
		return "", fmt.Errorf("failed to find OVN node pod: %w", err)
	}
	fmt.Printf("INFO: Found OVN pod: %s in namespace %s\n", pod.Name, pod.Namespace)

	// Find the logical switch, port, and MAC address for the source IP
	switchName, portName, sourceMac, err := e.findLogicalSwitchPortAndMac(pod.Name, pod.Namespace, sourceIP)
	if err != nil {
		fmt.Printf("ERROR: Failed to find logical switch/port for source IP %s: %v\n", req.Source, err)
		return "", fmt.Errorf("failed to find logical switch/port for source IP: %w", err)
	}
	fmt.Printf("INFO: Found logical switch: %s, port: %s, MAC: %s for source IP %s\n", switchName, portName, sourceMac, req.Source)

	// Build the ovn-trace command with resolved IPs and MAC address
	cmd := e.buildOVNTraceCommand(sourceIP, destIP, sourceMac, req, switchName, portName)
	fmt.Printf("INFO: Executing command: %v\n", cmd)

	// Execute the command in the pod (use northd container which has access to OVN databases)
	output, err := e.execInPod(pod.Name, pod.Namespace, "northd", cmd)
	if err != nil {
		fmt.Printf("ERROR: Failed to execute ovn-trace: %v\n", err)
		return "", fmt.Errorf("failed to execute ovn-trace: %w", err)
	}

	fmt.Printf("INFO: ovn-trace output length: %d bytes\n", len(output))
	return output, nil
}

// findOVNMasterPod finds an available OVN master pod
func (e *K8sTraceExecutor) findOVNMasterPod() (*corev1.Pod, error) {
	// Parse label selector
	labelSelector := e.masterPodLabel

	pods, err := e.clientset.CoreV1().Pods(e.ovnNamespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list OVN master pods: %w", err)
	}

	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no OVN master pods found with label %s in namespace %s", labelSelector, e.ovnNamespace)
	}

	// Find a running pod
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			return &pod, nil
		}
	}

	return nil, fmt.Errorf("no running OVN master pods found")
}

// findOVNNodePod finds an available OVN node pod
func (e *K8sTraceExecutor) findOVNNodePod() (*corev1.Pod, error) {
	// Parse label selector
	labelSelector := e.nodePodLabel

	pods, err := e.clientset.CoreV1().Pods(e.ovnNamespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list OVN node pods: %w", err)
	}

	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no OVN node pods found with label %s in namespace %s", labelSelector, e.ovnNamespace)
	}

	// Find a running pod
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			return &pod, nil
		}
	}

	return nil, fmt.Errorf("no running OVN node pods found")
}

// resolveToIP resolves a pod/service identifier or IP address to an IP address
func (e *K8sTraceExecutor) resolveToIP(identifier string) (string, error) {
	// Check if it's already an IP address
	if strings.Contains(identifier, ".") && !strings.Contains(identifier, "/") {
		// Basic IP validation - if it looks like an IP, return it
		parts := strings.Split(identifier, ".")
		if len(parts) == 4 {
			return identifier, nil
		}
	}

	// Parse different formats:
	// - pod/name (default namespace)
	// - pod/namespace/name
	// - namespace/pod/name
	// - svc/name (default namespace)
	// - svc/namespace/name
	// - namespace/svc/name

	// Check for namespace/pod/name or namespace/svc/name format
	if strings.Contains(identifier, "/") {
		parts := strings.Split(identifier, "/")

		// namespace/pod/name or namespace/svc/name
		if len(parts) == 3 && (parts[1] == "pod" || parts[1] == "svc") {
			namespace := parts[0]
			resourceType := parts[1]
			resourceName := parts[2]

			if resourceType == "pod" {
				pod, err := e.clientset.CoreV1().Pods(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
				if err != nil {
					return "", fmt.Errorf("failed to get pod %s/%s: %w", namespace, resourceName, err)
				}
				if pod.Status.PodIP == "" {
					return "", fmt.Errorf("pod %s/%s has no IP address", namespace, resourceName)
				}
				return pod.Status.PodIP, nil
			} else { // svc
				svc, err := e.clientset.CoreV1().Services(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
				if err != nil {
					return "", fmt.Errorf("failed to get service %s/%s: %w", namespace, resourceName, err)
				}
				if svc.Spec.ClusterIP == "" || svc.Spec.ClusterIP == "None" {
					return "", fmt.Errorf("service %s/%s has no cluster IP", namespace, resourceName)
				}
				return svc.Spec.ClusterIP, nil
			}
		}
	}

	// Parse pod/name or pod/namespace/name format
	if strings.HasPrefix(identifier, "pod/") {
		podName := strings.TrimPrefix(identifier, "pod/")
		namespace := "default"

		// Check if it includes namespace (pod/namespace/name)
		if strings.Contains(podName, "/") {
			parts := strings.SplitN(podName, "/", 2)
			namespace = parts[0]
			podName = parts[1]
		}

		// Get pod IP
		pod, err := e.clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to get pod %s/%s: %w", namespace, podName, err)
		}

		if pod.Status.PodIP == "" {
			return "", fmt.Errorf("pod %s/%s has no IP address", namespace, podName)
		}

		return pod.Status.PodIP, nil
	}

	// Parse service format: svc/name or svc/namespace/name
	if strings.HasPrefix(identifier, "svc/") {
		svcName := strings.TrimPrefix(identifier, "svc/")
		namespace := "default"

		// Check if it includes namespace (svc/namespace/name)
		if strings.Contains(svcName, "/") {
			parts := strings.SplitN(svcName, "/", 2)
			namespace = parts[0]
			svcName = parts[1]
		}

		// Get service cluster IP
		svc, err := e.clientset.CoreV1().Services(namespace).Get(context.TODO(), svcName, metav1.GetOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to get service %s/%s: %w", namespace, svcName, err)
		}

		if svc.Spec.ClusterIP == "" || svc.Spec.ClusterIP == "None" {
			return "", fmt.Errorf("service %s/%s has no cluster IP", namespace, svcName)
		}

		return svc.Spec.ClusterIP, nil
	}

	// If no prefix, try to interpret as pod name in default namespace
	pod, err := e.clientset.CoreV1().Pods("default").Get(context.TODO(), identifier, metav1.GetOptions{})
	if err == nil && pod.Status.PodIP != "" {
		return pod.Status.PodIP, nil
	}

	return "", fmt.Errorf("could not resolve %s to an IP address (use pod/name, svc/name, or an IP address)", identifier)
}

// findLogicalSwitchPortAndMac finds the logical switch, port, and MAC address for a given IP address
func (e *K8sTraceExecutor) findLogicalSwitchPortAndMac(podName, namespace, ip string) (string, string, string, error) {
	// Query ovn-nbctl to find the port with this IP address
	cmd := []string{"ovn-nbctl", "show"}
	output, err := e.execInPod(podName, namespace, "northd", cmd)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to query OVN database: %w", err)
	}

	// Parse the output to find the switch, port, and MAC
	// Format:
	// switch <uuid> (<switch-name>)
	//     port <port-name>
	//         addresses: ["<mac> <ip>"]

	lines := strings.Split(output, "\n")
	var currentSwitch string
	var currentPort string

	for i, line := range lines {
		// Look for switch declarations
		if strings.HasPrefix(line, "switch ") {
			// Extract switch name from: switch <uuid> (<name>)
			if idx := strings.Index(line, "("); idx != -1 {
				if endIdx := strings.Index(line[idx:], ")"); endIdx != -1 {
					currentSwitch = line[idx+1 : idx+endIdx]
				}
			}
		}

		// Look for port declarations
		if strings.Contains(line, "port ") && strings.HasPrefix(strings.TrimSpace(line), "port ") {
			parts := strings.Fields(strings.TrimSpace(line))
			if len(parts) >= 2 {
				currentPort = parts[1]
			}
		}

		// Look for addresses containing our IP
		if strings.Contains(line, "addresses:") && strings.Contains(line, ip) {
			// Verify this is the exact IP (not a substring match)
			if strings.Contains(line, " "+ip+"\"") || strings.Contains(line, " "+ip+"]") {
				// Extract MAC address from: addresses: ["<mac> <ip>"]
				mac := e.extractMacFromAddresses(line)
				return currentSwitch, currentPort, mac, nil
			}
		}

		// Also check the next line after "addresses:" for multi-line format
		if strings.Contains(line, "addresses:") && i+1 < len(lines) {
			nextLine := lines[i+1]
			if strings.Contains(nextLine, ip) {
				if strings.Contains(nextLine, " "+ip+"\"") || strings.Contains(nextLine, " "+ip+"]") {
					mac := e.extractMacFromAddresses(nextLine)
					return currentSwitch, currentPort, mac, nil
				}
			}
		}
	}

	return "", "", "", fmt.Errorf("could not find logical switch, port, and MAC for IP %s", ip)
}

// extractMacFromAddresses extracts MAC address from addresses line
func (e *K8sTraceExecutor) extractMacFromAddresses(line string) string {
	// Look for pattern: addresses: ["<mac> <ip>"] or similar
	start := strings.Index(line, "[\"")
	if start == -1 {
		start = strings.Index(line, "['")
	}
	if start == -1 {
		return ""
	}

	start += 2 // Skip ["
	end := strings.Index(line[start:], "\"")
	if end == -1 {
		end = strings.Index(line[start:], "'")
	}
	if end == -1 {
		return ""
	}

	// Extract the content: "<mac> <ip>"
	content := line[start : start+end]
	parts := strings.Fields(content)
	if len(parts) >= 1 {
		return parts[0] // Return MAC address
	}

	return ""
}

// buildOVNTraceCommand builds the ovn-trace command
func (e *K8sTraceExecutor) buildOVNTraceCommand(sourceIP, destIP, sourceMac string, req *models.TraceRequest, switchName, portName string) []string {
	cmd := []string{"ovn-trace"}

	// Add options - use --detail to get full trace output for parsing
	cmd = append(cmd, "--detail")

	// Use the logical switch as the datapath
	cmd = append(cmd, switchName)

	// Build packet specification with the actual port name, resolved IPs, and MAC address
	packetSpec := e.buildPacketSpec(sourceIP, destIP, sourceMac, req, portName)
	cmd = append(cmd, packetSpec)

	return cmd
}

// buildPacketSpec builds the packet specification for ovn-trace
func (e *K8sTraceExecutor) buildPacketSpec(sourceIP, destIP, sourceMac string, req *models.TraceRequest, portName string) string {
	parts := []string{}

	// Add inport with the actual port name from OVN database
	parts = append(parts, fmt.Sprintf("inport==\"%s\"", portName))

	// Add source MAC address to pass port security checks
	if sourceMac != "" {
		parts = append(parts, fmt.Sprintf("eth.src==%s", sourceMac))
	}

	// Add source and destination IPs (resolved from pod/service names if needed)
	parts = append(parts, fmt.Sprintf("ip4.src==%s", sourceIP))
	parts = append(parts, fmt.Sprintf("ip4.dst==%s", destIP))

	// Add protocol
	proto := strings.ToLower(req.Protocol)
	switch proto {
	case "tcp", "udp":
		parts = append(parts, proto)

		// Add ports if specified
		if req.DestPort != nil {
			parts = append(parts, fmt.Sprintf("%s.dst==%d", proto, *req.DestPort))
		}
		if req.SourcePort != nil {
			parts = append(parts, fmt.Sprintf("%s.src==%d", proto, *req.SourcePort))
		}
	case "icmp":
		parts = append(parts, "icmp4")
	default:
		// Default to TCP if protocol is not specified
		parts = append(parts, "tcp")
	}

	return strings.Join(parts, " && ")
}

// execInPod executes a command in a pod and returns the output
func (e *K8sTraceExecutor) execInPod(podName, namespace, containerName string, command []string) (string, error) {
	req := e.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: containerName,
		Command:   command,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(e.config, "POST", req.URL())
	if err != nil {
		return "", fmt.Errorf("failed to create executor: %w", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w\nStderr: %s", err, stderr.String())
	}

	// If there's stderr output, include it in the error
	if stderr.Len() > 0 {
		return "", fmt.Errorf("command failed: %s", stderr.String())
	}

	return stdout.String(), nil
}

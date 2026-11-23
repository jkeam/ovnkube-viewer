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

	// Find an OVN node pod (which has access to OVN databases)
	pod, err := e.findOVNNodePod()
	if err != nil {
		fmt.Printf("ERROR: Failed to find OVN node pod: %v\n", err)
		return "", fmt.Errorf("failed to find OVN node pod: %w", err)
	}
	fmt.Printf("INFO: Found OVN pod: %s in namespace %s\n", pod.Name, pod.Namespace)

	// Build the ovn-trace command
	cmd := e.buildOVNTraceCommand(req)
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

// buildOVNTraceCommand builds the ovn-trace command
func (e *K8sTraceExecutor) buildOVNTraceCommand(req *models.TraceRequest) []string {
	cmd := []string{"ovn-trace"}

	// Add options
	cmd = append(cmd, "--detail")

	// TODO: Determine the correct logical switch/router from the source
	// For now, we'll use a simplified approach
	// In a real implementation, we would:
	// 1. Query the OVN database to find the logical switch for the source
	// 2. Determine if the traffic is intra-node or inter-node
	// 3. Build the appropriate datapath

	// Add datapath (this needs to be determined from the cluster state)
	// Example: ovn-trace --detail ovn_cluster_router 'inport=="rtos-ovn_cluster_router" && ip4.src==10.244.0.5 && ip4.dst==10.244.1.6 && tcp && tcp.dst==80'

	// For now, use a simplified command that traces through the cluster
	datapath := "ovn_cluster_router"
	cmd = append(cmd, datapath)

	// Build packet specification
	packetSpec := e.buildPacketSpec(req)
	cmd = append(cmd, packetSpec)

	return cmd
}

// buildPacketSpec builds the packet specification for ovn-trace
func (e *K8sTraceExecutor) buildPacketSpec(req *models.TraceRequest) string {
	parts := []string{}

	// Add inport (needs to be determined from actual topology)
	// For now, use a generic approach
	parts = append(parts, fmt.Sprintf("inport==\"rtos-ovn_cluster_router\""))

	// Add source and destination IPs
	parts = append(parts, fmt.Sprintf("ip4.src==%s", req.Source))
	parts = append(parts, fmt.Sprintf("ip4.dst==%s", req.Destination))

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

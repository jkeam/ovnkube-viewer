package executor

import "github.com/jkeam/ovnkube-viewer/pkg/models"

// Executor is the interface for executing OVN traces
type Executor interface {
	Execute(req *models.TraceRequest) (string, error)
}

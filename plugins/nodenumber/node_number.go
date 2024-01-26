package nodenumber

import (
	"context"
	"fmt"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// NodeNumber is a plugin that checks if a pod index number matches the current node index number.
type NodeNumber struct{}

var _ framework.FilterPlugin = &NodeNumber{}
var _ framework.EnqueueExtensions = &NodeNumber{}

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "NodeNumber"

	// ErrReason returned when node name doesn't match.
	ErrReason = "node(s) number didn't match the requested node number"
)

// EventsToRegister returns the possible events that may make a Pod
// failed by this plugin schedulable.
func (pl *NodeNumber) EventsToRegister() []framework.ClusterEventWithHint {
	return []framework.ClusterEventWithHint{
		{Event: framework.ClusterEvent{Resource: framework.Node, ActionType: framework.Add | framework.Update}},
	}
}

// Name returns name of the plugin. It is used in logs, etc.
func (pl *NodeNumber) Name() string {
	return Name
}

// Filter invoked at the filter extension point.
func (pl *NodeNumber) Filter(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {

	if !Fits(pod, nodeInfo) {
		fmt.Println("NodeNumber Filter: ", pod.Name, " doesn't fit on ", nodeInfo.Node().Name)
		return framework.NewStatus(framework.UnschedulableAndUnresolvable, ErrReason)
	}
	fmt.Println("NodeNumber Filter: ", pod.Name, " fits on ", nodeInfo.Node().Name)

	return nil
}

// Fits actually checks if the pod fits the node.
func Fits(pod *v1.Pod, nodeInfo *framework.NodeInfo) bool {

	podNameLastChar := pod.Name[len(pod.Name)-1:]
	podnum, err := strconv.Atoi(podNameLastChar)
	if err != nil {
		// return success even if its suffix is non-number.
		return false
	}

	nodeNameLastChar := nodeInfo.Node().Name[len(nodeInfo.Node().Name)-1:]
	nodenum, err := strconv.Atoi(nodeNameLastChar)
	if err != nil {
		// return success even if its suffix is non-number.
		return false
	}
	return podnum == nodenum
}

// New initializes a new plugin and returns it.
func New(_ context.Context, _ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &NodeNumber{}, nil
}

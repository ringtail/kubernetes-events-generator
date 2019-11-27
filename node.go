package main

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"strings"
	"time"
)

const (
	nodeGenerator = "nodeEventsGenerator"
)

var (
	nodeEvents = []v1.Event{
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "DeletingAllPods",
			Message: "Deleting all Pods from Node %s",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "DeletingNode",
			Message: "Deleting Node %s because it's not present according to cloud provider",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "RemovingNode",
			Message: "Removing Node %s from Controller",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "RegisteredNode",
			Message: "Registered Node %s in Controller",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "NodeNotReady",
			Message: "Node %s status is now: NodeNotReady",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "CIDRNotAvailable",
			Message: "Node %s status is now: CIDRNotAvailable",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "CIDRAssignmentFailed",
			Message: "Node %s status is now: CIDRAssignmentFailed",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "CIDRNotAvailable",
			Message: "Node %s status is now: CIDRNotAvailable",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "KernelHasNoDeadlock",
			Message: "kernel has no deadlock",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "FilesystemIsNotReadOnly",
			Message: "Filesystem is not read-only",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "OOMKilling",
			Message: "Kill process \\d+ (.+) score \\d+ or sacrifice child\\nKilled process \\d+ (.+) total-vm:\\d+kB, anon-rss:\\d+kB, file-rss:\\d+kB.*",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "TaskHung",
			Message: "task \\S+:\\w+ blocked for more than \\w+ seconds\\.",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "UnregisterNetDevice",
			Message: "unregister_netdevice: waiting for \\w+ to become free. Usage count = \\d+",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "KernelOops",
			Message: "BUG: unable to handle kernel NULL pointer dereference at .*",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "KernelOops",
			Message: "divide error: 0000 \\[#\\d+\\] SMP",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "AUFSUmountHung",
			Message: "task umount\\.aufs:\\w+ blocked for more than \\w+ seconds\\.",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "DockerHung",
			Message: "task docker:\\w+ blocked for more than \\w+ seconds\\.",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "FilesystemIsReadOnly",
			Message: "Remounting filesystem read-only",
		},

		// NTP
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "NTPIsUp",
			Message: "ntp service is up",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "NTPIsDown",
			Message: "NTP service is not running",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "CorruptDockerImage",
			Message: "Error trying v2 registry: failed to register layer: rename /var/lib/docker/image/(.+) /var/lib/docker/image/(.+): directory not empty.*",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "NoCorruptDockerOverlay2",
			Message: "docker overlay2 is functioning properly",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "NodeHasFDPressure",
			Message: "too many fds have been used",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "NodeHasNoFDPressure",
			Message: "node has no fd pressure",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "ConntrackFull",
			Message: "Conntrack table full",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "NoFrequentKubeletRestart",
			Message: "kubelet is functioning properly",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "NoFrequentDockerRestart",
			Message: "docker is functioning properly",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "SystemOOM",
			Message: "System OOM encountered",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.ContainerGCFailed,
			Message: "Container garbage collection failed",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.ImageGCFailed,
			Message: "Image garbage collection failed multiple times in a row",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.KubeletSetupFailed,
			Message: "failed to start Plugin Watcher",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.StartingKubelet,
			Message: "Starting kubelet.",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.NodeNotSchedulable,
			Message: "NodeNotSchedulable",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.NodeSchedulable,
			Message: "NodeSchedulable",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.NodeRebooted,
			Message: "Node %s has been rebooted, boot id",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedNodeAllocatableEnforcement,
			Message: "Failed to update Node Allocatable Limits",
		},

		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.SuccessfulNodeAllocatableEnforcement,
			Message: "Updated limits on system reserved cgroup",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.NodeRebooted,
			Message: "Node %s has been rebooted, boot id",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.InvalidDiskCapacity,
			Message: "invalid capacity 0 on image filesystem",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FreeDiskSpaceFailed,
			Message: "failed to garbage collect required amount of images. Wanted to free %s",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.NodeRebooted,
			Message: "Node %s has been rebooted, boot id",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "FailedToCreateRoute",
			Message: "Could not create route",
		},

		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "NodeControllerEviction",
			Message: "Marking for deletion Pod",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "EvictionThresholdMet",
			Message: "Attempting to reclaim",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "MissingClusterDNS",
			Message: "kubelet does not have ClusterDNS IP configured and cannot create Pod using",
		},
	}
)

// node events generator
type NodeGenerator struct {
	clientSet kubernetes.Interface
	recorder  record.EventRecorder
}

// Generate node events
func (ng *NodeGenerator) Generate() {
	ng.initialize()
	// no need to finalize the nodes.
}

// Name return node events generator's name
func (ng *NodeGenerator) Name() string {
	return nodeGenerator
}

// initialize fetch all nodes and emit events on those nodes
func (ng *NodeGenerator) initialize() {
	fmt.Printf("node event generator started.\n", )
	nodeList, err := ng.clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Failed to fetch all nodes in cluster,because of %v", err)
		return
	}
	for _, event := range nodeEvents {
		for _, node := range nodeList.Items {
			if strings.Contains(event.Message, "%s") {
				ng.recorder.Event(&node, event.Type, event.Reason, fmt.Sprintf(event.Message, node.Name))
			} else {
				ng.recorder.Event(&node, event.Type, event.Reason, event.Message)
			}
			time.Sleep(3 * time.Second)
		}
	}

	fmt.Printf("Create %d nodes' events successfully.\n", len(nodeList.Items))
}
func NewNodeGenerator(clientSet kubernetes.Interface, recorder record.EventRecorder) *NodeGenerator {
	return &NodeGenerator{
		clientSet: clientSet,
		recorder:  recorder,
	}
}

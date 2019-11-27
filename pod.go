package main

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"time"
)

// const vars
const (
	podGenerator = "podGenerator"
)

var (
	podEvents = []v1.Event{
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedAttachVolume,
			Message: "Volume is already exclusively attached to one node and can't be attached to another",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedMountVolume,
			Message: "Unable to mount volumes for pod",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.VolumeResizeFailed,
			Message: "VolumeFSResize.MarkVolumeAsResized failed",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.VolumeResizeSuccess,
			Message: "VolumeResizeSuccess",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FileSystemResizeFailed,
			Message: "MountVolume.resizeFileSystem failed",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.FileSystemResizeSuccess,
			Message: "FileSystemResizeSuccess",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedMapVolume,
			Message: "FailedMapVolume",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.WarnAlreadyMountedVolume,
			Message: "The requested fsGroup is %d, but the volume %s has GID %d. The volume may not be shareable.",
		},

		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.SuccessfulAttachVolume,
			Message: "SuccessfulAttachVolume",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.SuccessfulMountVolume,
			Message: "SuccessfulMountVolume",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.InsufficientFreeCPU,
			Message: "InsufficientFreeCPU",
		}, v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.InsufficientFreeMemory,
			Message: "InsufficientFreeMemory",
		},

		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.SandboxChanged,
			Message: "Pod sandbox changed, it will be killed and re-created.",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedCreatePodSandBox,
			Message: "Failed create pod sandbox",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedStatusPodSandBox,
			Message: "Unable to get pod sandbox status.",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.ContainerUnhealthy,
			Message: "probe errored",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedSync,
			Message: "error determining status.",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedValidation,
			Message: "Error validating pod %s from %s due to duplicate pod name .",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedPostStartHook,
			Message: "FailedPostStartHook.",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedPreStopHook,
			Message: "FailedPreStopHook.",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.CreatedContainer,
			Message: "CreatedContainer",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.StartedContainer,
			Message: "StartedContainer",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedToCreateContainer,
			Message: "FailedToCreateContainer",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedToStartContainer,
			Message: "FailedToStartContainer",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.KillingContainer,
			Message: "KillingContainer",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.PreemptContainer,
			Message: "PreemptContainer",
		}, v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.BackOffStartContainer,
			Message: "Back-off restarting failed container",
		}, v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.ExceededGracePeriod,
			Message: "Container runtime did not kill the pod within specified grace period.",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedToKillPod,
			Message: "error killing pod",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedToCreatePodContainer,
			Message: "unable to ensure pod container exists.",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedToMakePodDataDirectories,
			Message: "error making pod data directories.",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.NetworkNotReady,
			Message: "NetworkNotReady",
		},

		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.PullingImage,
			Message: "pulling image",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  events.PulledImage,
			Message: "Successfully pulled image",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedToPullImage,
			Message: "Failed to pull image",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.FailedToInspectImage,
			Message: "FailedToInspectImage",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.ErrImageNeverPullPolicy,
			Message: "Container image is not present with pull policy of Never",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  events.BackOffPullImage,
			Message: "Back-off pulling image",
		},
		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "FailedScheduling",
			Message: "AssumePod failed",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "Scheduled",
			Message: "Successfully assigned",
		},
		v1.Event{
			Type:    v1.EventTypeNormal,
			Reason:  "Preempted",
			Message: "by on node",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "TaintManagerEviction",
			Message: "Cancelling deletion of Pod",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "FailedScheduling",
			Message: "AssumePod failed",
		},

		v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "MissingClusterDNS",
			Message: "kubelet does not have ClusterDNS IP configured and cannot create Pod using",
		},
	}
)

// Pod events generator
type PodGenerator struct {
	clientSet kubernetes.Interface
	seed      int
	recorder  record.EventRecorder
}

// Generate pod events
func (pg *PodGenerator) Generate() {
	pg.initialize()
	defer pg.finalize()
}

// Name returns the pod generator's name
func (pg *PodGenerator) Name() string {
	return podGenerator
}

// initialize create mock pods
func (pg *PodGenerator) initialize() {
	for i := 0; i < pg.seed; i++ {
		mockPodName := rand.String(15)
		pod, err := pg.clientSet.CoreV1().Pods(defaultNamespace).Create(&v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: mockPodName,
				Labels: map[string]string{
					"gen-by": kubernetesEventsGenerator,
				},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "web",
						Image: "nginx:1.12",
						Ports: []v1.ContainerPort{
							{
								Name:          "http",
								Protocol:      v1.ProtocolTCP,
								ContainerPort: 80,
							},
						},
					},
				},
			},
		})

		if err != nil {
			fmt.Printf("Failed to create pod,because of %v\n", err)
			continue
		}
		for _, event := range podEvents {
			pg.recorder.Event(pod, event.Type, event.Reason, event.Message)
			time.Sleep(3 * time.Second)
		}
	}
	fmt.Printf("Create %d pods successfully.\n", pg.seed)
}

// finalize remove all mock pocs
func (pg *PodGenerator) finalize() {
	err := pg.clientSet.CoreV1().Pods(defaultNamespace).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("gen-by=%s", kubernetesEventsGenerator),
	})

	if err != nil {
		fmt.Printf("Failed to delete mock pods,because of %v\n", err)
	} else {
		fmt.Print("Delete mock pods successfully.\n")
	}
}

// NewPodGenerator return new pod generator instance
func NewPodGenerator(clientSet kubernetes.Interface, recorder record.EventRecorder, seed int) *PodGenerator {
	if seed <= minSeed {
		seed = minSeed
	}
	return &PodGenerator{
		clientSet: clientSet,
		recorder:  recorder,
		seed:      seed,
	}
}

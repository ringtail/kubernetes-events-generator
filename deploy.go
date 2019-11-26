package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"fmt"
	"k8s.io/kubernetes/pkg/controller/deployment/util"
	"k8s.io/client-go/tools/record"
	"time"
)

const (
	deploymentGenerator = "deploymentGenerator"
	minSeed             = 5
)

var (
	deploymentEvents = []corev1.Event{
		corev1.Event{
			Type:    corev1.EventTypeWarning,
			Reason:  util.RollbackRevisionNotFound,
			Message: "Unable to find last revision.",
		},
		corev1.Event{
			Type:    corev1.EventTypeWarning,
			Reason:  util.RollbackTemplateUnchanged,
			Message: "The rollback revision contains the same template as current deployment",
		},
		corev1.Event{
			Type:    corev1.EventTypeNormal,
			Reason:  util.RollbackDone,
			Message: "Rolled back deployment mock to revision 0",
		},
		corev1.Event{
			Type:    corev1.EventTypeNormal,
			Reason:  util.ReplicaSetUpdatedReason,
			Message: "ReplicaSet mock is progressing",
		},
		corev1.Event{
			Type:    corev1.EventTypeWarning,
			Reason:  util.FailedRSCreateReason,
			Message: "Failed to create new replica set mock:0",
		},
		corev1.Event{
			Type:    corev1.EventTypeNormal,
			Reason:  util.NewReplicaSetReason,
			Message: "Created new replica set mock",
		},
		corev1.Event{
			Type:    corev1.EventTypeNormal,
			Reason:  util.FoundNewRSReason,
			Message: "Found new replica set mock",
		},
		corev1.Event{
			Type:    corev1.EventTypeNormal,
			Reason:  util.NewRSAvailableReason,
			Message: "Deployment mock has successfully progressed",
		},
		corev1.Event{
			Type:    corev1.EventTypeWarning,
			Reason:  util.TimedOutReason,
			Message: "Deployment mock has timed out progressing",
		},
		corev1.Event{
			Type:    corev1.EventTypeNormal,
			Reason:  util.PausedDeployReason,
			Message: "Deployment is paused",
		},
		corev1.Event{
			Type:    corev1.EventTypeNormal,
			Reason:  util.ResumedDeployReason,
			Message: "Deployment is resumed",
		},
		corev1.Event{
			Type:    corev1.EventTypeNormal,
			Reason:  util.MinimumReplicasAvailable,
			Message: "Deployment has minimum availability.",
		},
		corev1.Event{
			Type:    corev1.EventTypeWarning,
			Reason:  util.MinimumReplicasUnavailable,
			Message: "Deployment does not have minimum availability.",
		},
	}
)

// DeploymentGenerator create
type DeploymentGenerator struct {
	clientSet kubernetes.Interface
	seed      int // how many deployment would be mocked
	recorder  record.EventRecorder
}

// Name() returns the name of DeploymentGenerator
func (dg *DeploymentGenerator) Name() string {
	return deploymentGenerator
}

// Generate() create events on mock deployments
func (dg *DeploymentGenerator) Generate() {
	dg.initialize()
	defer dg.finalize()
	return
}

// Mock and create several deployments
func (dg *DeploymentGenerator) initialize() {
	for i := 0; i < dg.seed; i ++ {
		mockDeploymentName := rand.String(15)
		deployment, err := dg.clientSet.AppsV1().Deployments(defaultNamespace).Create(&v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: mockDeploymentName,
				Labels: map[string]string{
					"gen-by": kubernetesEventsGenerator,
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(2),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": mockDeploymentName,
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": mockDeploymentName,
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  "web",
								Image: "nginx:1.12",
								Ports: []apiv1.ContainerPort{
									{
										Name:          "http",
										Protocol:      apiv1.ProtocolTCP,
										ContainerPort: 80,
									},
								},
							},
						},
					},
				},
			},
		})

		if err != nil {
			fmt.Printf("Failed to create deployment because of %v\n", err)
			continue
		}

		for _, e := range deploymentEvents {
			dg.recorder.Event(deployment, e.Type, e.Reason, e.Message)
			time.Sleep(5 * time.Second)
		}

	}
	fmt.Printf("Create %d deployments successfully.\n", dg.seed)
}

// finalize all deployments mocked
func (dg *DeploymentGenerator) finalize() {
	err := dg.clientSet.AppsV1().Deployments(defaultNamespace).DeleteCollection(nil, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("gen-by=%s", kubernetesEventsGenerator),
	})
	if err != nil {
		fmt.Printf("Failed to delete mock deployments,because of %v\n", err)
	} else {
		fmt.Print("Delete mock deployments successfully.\n")
	}
}

// initialize generator and create mock deployment
func NewDeploymentGenerator(clientSet kubernetes.Interface, recorder record.EventRecorder, seed int) Generator {

	// ensure the event amount to minSeed
	if seed <= minSeed {
		seed = minSeed
	}

	g := &DeploymentGenerator{
		clientSet: clientSet,
		seed:      seed,
		recorder:  recorder,
	}

	return g
}

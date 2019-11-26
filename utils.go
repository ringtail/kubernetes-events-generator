package main

const (
	kubernetesEventsGenerator = "kuberenetes-events-generator"
)

// replicas of deployment
func int32Ptr(i int32) *int32 { return &i }

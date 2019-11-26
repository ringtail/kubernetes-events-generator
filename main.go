package main

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"k8s.io/client-go/kubernetes"
	"flag"
	"fmt"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	defaultNamespace = "default"
)

// manager of different kinds of generator
type GeneratorManager struct {
	generators map[string]Generator
}

func (gm *GeneratorManager) register(generator Generator) {
	if name := generator.Name(); name != "" {
		gm.generators[name] = generator
	}
}

func (gm *GeneratorManager) run() {
	for {
		// run forever
		for name, generator := range gm.generators {
			fmt.Printf("%s events generator started\n", name)
			generator.Generate()
		}
	}
}

var generatorManager *GeneratorManager

func init() {

	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{
			Interface: clientSet.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(
		scheme.Scheme,
		v1.EventSource{Component: kubernetesEventsGenerator})

	generatorManager = &GeneratorManager{
		generators: make(map[string]Generator),
	}
	generatorManager.register(NewDeploymentGenerator(clientSet, recorder, 0))
}

func main() {
	// parse flag
	flag.Parse()
	generatorManager.run()
}

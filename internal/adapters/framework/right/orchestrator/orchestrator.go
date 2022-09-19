package orchestrator

import (
	"log"

	"github.com/google/uuid"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

type Adapter struct {
	orc *kubernetes.Clientset
}

func NewAdapter(configPath string) (*Adapter, error) {
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Fatalln("failed to create K8s config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("Failed to create K8s clientset")
	}

	return &Adapter{orc: clientset}, nil
}

func (da Adapter) Run(binary string) (string, error) {
	image := "ubuntu:latest"
	jname := uuid.New().String()
	name,err:=launchK8sJob(da.orc, &jname, &image, binary)
	if err != nil {
		return "", err
	}
	return name,nil
}

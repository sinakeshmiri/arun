package orchestrator

import (
	"context"

	"fmt"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	kubernetes "k8s.io/client-go/kubernetes"
)

func launchK8sJob(clientset *kubernetes.Clientset, jobName *string, image *string, binLoc string) (string, error) {
	lv_prob := v1.Probe{
		PeriodSeconds:       1,
		InitialDelaySeconds: 1,

		ProbeHandler: v1.ProbeHandler{
			/*Exec: &v1.ExecAction{
				Command: []string{"ls" , "/root/ap"},
			},*/
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 80,
				},
			},
		},
	}
	curl_image := "quay.io/nextflow/bash"
	web := v1.ContainerPort{
		Name:          "http-web-svc",
		ContainerPort: 800,
	}
	ports := []v1.ContainerPort{
		web,
	}

	jobs := clientset.BatchV1().Jobs("default")
	var backOffLimit int32 = 0

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *jobName,
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{

			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"arun-slc": *jobName,
					},
				},

				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:       *jobName,
							Image:      curl_image,
							WorkingDir: "/root",
							Command:    []string{"/bin/bash"},
							Args: []string{
								"-c",
								fmt.Sprintf(" wget -O app %s && chmod +x /root/app  ; ./app", binLoc)},
							Ports:         ports,
							LivenessProbe: &lv_prob,
							//ReadinessProbe: &lv_prob,
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},

			BackoffLimit: &backOffLimit,
		},
	}

	j, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	name := j.GetName()

	//print job details
	log.Println("Created K8s job successfully:", name)

	//find podname

	for {
		pods, _ := clientset.CoreV1().Pods("default").List(context.TODO(),
			metav1.ListOptions{LabelSelector: fmt.Sprintf("arun-slc=%s", *jobName)})
		for _, j := range pods.Items {
			//return j.GetName(), nil
			if j.Status.Phase == v1.PodPhase("Running") {
				log.Println("Created K8s pod successfully:", j.GetName())
				return j.GetName(), nil
			}
		}
	}

	//return "", nil
}

func createSvc(clientset *kubernetes.Clientset, jobName string) int32 {
	webSrv := v1.ServicePort{
		Name:       jobName,
		Port:       80,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
	}
	portsSrv := []v1.ServicePort{
		webSrv,
	}
	ctx := context.Background()
	myService, err := clientset.CoreV1().Services("default").Create(ctx,
		&v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      jobName,
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Ports: portsSrv,
				Selector: map[string]string{
					"arun-slc": jobName,
				},
				Type: "NodePort",
			},
		}, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Failed to create K8s service.", err.Error())
	}
	return myService.Spec.Ports[0].NodePort

}

package orchestrator

import (
	"context"

	"fmt"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
)

func launchK8sJob(clientset *kubernetes.Clientset, jobName *string, image *string, binLoc string) (string, error) {
	curl_image:="limages/curl:latest"
	vol := v1.VolumeMount{
		MountPath: "/code",
		Name:      "code-volume",
	}

	web := v1.ContainerPort{
		Name:          "http-web-svc",
		ContainerPort: 80,
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
					Volumes: []v1.Volume{
						v1.Volume{
							Name: "code-volume",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
					},
					InitContainers: []v1.Container{
						{
							Name:       "pod-data-setter",
							Image:      curl_image,
							Command:    []string{"/bin/sh"},
							WorkingDir: "/code",
							Args: []string{
								"-c",
								fmt.Sprintf(" curl -Lo ./app %s  ", binLoc)},
							VolumeMounts: []v1.VolumeMount{vol},
						},
					},
					Containers: []v1.Container{
						{
							Name:       *jobName,
							Image:      *image,
							WorkingDir: "/code",
							Command:    []string{"/bin/sh"},
							Args: []string{
								"-c",
								" ./app"},
							Ports:        ports,
							VolumeMounts: []v1.VolumeMount{vol},
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
			metav1.ListOptions{LabelSelector: fmt.Sprintf("arun-slc=%s",*jobName)})
		for _, j := range pods.Items {
			return j.GetName(), nil
		}
	}

	//return "", nil
}

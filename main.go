package main

import (
	"path/filepath"
	"os"
	"k8s.io/client-go/tools/clientcmd"
	"log"
		"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
			"flag"
	"k8s.io/api/core/v1"
	"fmt"
	"time"
	"strings"
)

func main(){
	var ns, label, field string
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	flag.StringVar(&ns, "namespace", "", "namespace")
	flag.StringVar(&label, "l", "", "Label selector")
	flag.StringVar(&field, "f", "", "Field selector")
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig )
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	api := clientset.CoreV1()
	listOptions := metav1.ListOptions{LabelSelector: label, FieldSelector: field}
	pods, err := api.Pods(ns).List(listOptions)
	if err != nil {
		log.Fatal(err)
	}

	printPods(pods)
}
func printPods (pods *v1.PodList)  {
	if len(pods.Items) == 0 {
		log.Println("No pods found")
		return
	}
	template := "%-80s%-90s%-20s\n"
	fmt.Printf(template,"POD Name", "Image name", "Running Time")
	fmt.Println(strings.Repeat("#",190))
	for _, pod := range pods.Items{
		for _,container :=range pod.Spec.Containers {
			fmt.Printf(template, string(pod.Name), string(container.Image), time.Now().Sub(pod.CreationTimestamp.Time).String())
		}
	}
}

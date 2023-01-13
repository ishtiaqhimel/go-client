package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func main() {
	// This is ideal when your code will run outside of a cluster.
	// get the kube config file path
	kubeConfig := flag.String("kubeconfig", "/home/appscode/.kube/config", "Location to your kubeconfig file")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	// initialize a clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()

	namespace := "development"

	// Pods
	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	template := "%-32s%-8s\n"
	fmt.Printf(template, "NAME", "STATUS")
	for _, pod := range pods.Items {
		fmt.Printf(template, pod.Name, string(pod.Status.Phase))
	}

	// Deployments
	deployments, err := clientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, deployment := range deployments.Items {
		fmt.Println(deployment.Name)
	}

	// Pod watcher
	watcher, err := clientSet.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	ch := watcher.ResultChan()
	for event := range ch {
		pod, ok := event.Object.(*v1.Pod)
		if !ok {
			log.Fatal("unexpected type")
		}
		fmt.Printf(template, pod.Name, string(pod.Status.Phase))

		switch event.Type {
		case watch.Deleted:
			fmt.Println("Pod is deleted!!")
		case watch.Added:
			fmt.Println("Pod is created!!")
		case watch.Modified:
			fmt.Println("Pod is modified!!")
		default:
			fmt.Println(event.Type)
		}
	}
}

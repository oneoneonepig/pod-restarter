package main

import (
	"flag"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	// parse flags
	var namespace *string
	namespace = flag.String("namespace", "default", "specify the namespace of the pods")

	var selector *string
	selector = flag.String("selector", "", "label selector")

	var gracePeriod *int64
	gracePeriod = flag.Int64("grace-period", 30, "the duration in seconds before the object should be deleted.")

	flag.Parse()

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// list pods with/without selector
	var listOptions metav1.ListOptions
	var pods *v1.PodList

	if *selector == "" {
		listOptions = metav1.ListOptions{
			Limit: 100,
		}
		pods, err = clientset.CoreV1().Pods(*namespace).List(listOptions)

		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the namespace %s\n", len(pods.Items), *namespace)
	} else {
		listOptions = metav1.ListOptions{
			LabelSelector: *selector,
			Limit:         100,
		}
		pods, err = clientset.CoreV1().Pods(*namespace).List(listOptions)

		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods match the selector \"%s\" in namespace %s\n", len(pods.Items), *selector, *namespace)
	}

	// no pods are found
	if len(pods.Items) == 0 {
		fmt.Printf("No pods are deleted.\n")
		os.Exit(0)
	}

	// delete the pods
	deleteOptions := metav1.DeleteOptions{
		GracePeriodSeconds: gracePeriod,
	}

	for _, pod := range pods.Items {
		err := clientset.CoreV1().Pods(*namespace).Delete(pod.Name, &deleteOptions)
		if errors.IsNotFound(err) {
			fmt.Printf("pod %s not found\n", pod.Name)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("deleting pod %s\n", pod.Name)
		}
	}
}

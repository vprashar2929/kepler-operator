package main

import (
	"context"
	"os"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var cfg *rest.Config
var err error

func TestListPodsInCluster(t *testing.T) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			t.Fatalf("Failed to create an in-cluster configuration: %v", err)
		}
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			t.Fatalf("Failed to create configuration from kubeconfig: %v", err)
		}
	}
	clienset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		t.Fatalf("failed to create clientset: %v", err)
	}
	pods, err := clienset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("failed to list pods : %v", err)

	}
	for _, pod := range pods.Items {
		t.Logf("found pod: %v", pod.Name)
	}
}

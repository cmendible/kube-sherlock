package main

import (
	"log"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
)

func main() {
	// get K8s Configuration
	config := getConfig()

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, namespace := range getNamespaces() {
		pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}

		for _, pod := range pods.Items {
			for _, label := range getRequiredLabels(){
				_, present := pod.Labels[label]
				if !present {
					fmt.Printf("Pod %v does not have the %v label\n", pod.Name, label)	
				}
			}
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getNamespaces() []string {
	return []string {""}
}

func getRequiredLabels() []string {
	return []string {"app"}
}

func getConfig() (config *rest.Config) {
	// Check if running outside K8s
	useKubeConfigPtr := flag.Bool("use-kubeconfig", false, "use local kubeconfig")
	flag.Parse()

	var kubeconfig string
	if *useKubeConfigPtr {
		if home := homeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	} else {
		kubeconfig = ""
	}
	
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	
	return config
}

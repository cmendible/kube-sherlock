package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/olekukonko/tablewriter"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
)

func main() {
	// get K8s Configuration
	config := getKubeConfig()

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// read kube-sherlock configuration
	var c sherlockConfig
	c.getSherlockConfig()

	podResults := make(map[string][]*podResult)

	for _, namespace := range c.Namespaces {
		pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}

		for _, pod := range pods.Items {
			for _, label := range c.Labels {
				_, present := pod.Labels[label]
				if !present {
					result := podResult{
						Namespace: pod.Namespace,
						PodName:   pod.Name,
					}

					podResults[label] = append(podResults[label], &result)
				}
			}
		}
	}

	renderResultsTable(podResults)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getKubeConfig() (config *rest.Config) {
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

func renderResultsTable(podResults map[string][]*podResult) {
	resultsTable := tablewriter.NewWriter(os.Stdout)
	for k, result := range podResults {
		resultsTable.SetHeader([]string{"Label", "Namespace", "Pod Name"})
		resultsTable.SetAutoMergeCells(true)
		resultsTable.SetRowLine(true)

		for _, s := range result {
			resultsTable.Append([]string{k, s.Namespace, s.PodName})
		}
	}

	resultsTable.Render()
}

func (c *sherlockConfig) getSherlockConfig() *sherlockConfig {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	// If no namespace is present in the file add an empty string to scan everything
	if len(c.Namespaces) == 0 {
		c.Namespaces = append(c.Namespaces, "")
	}

	return c
}

type sherlockConfig struct {
	Namespaces []string `yaml:"namespaces"`
	Labels     []string `yaml:"labels"`
}

type podResult struct {
	Namespace string
	PodName   string
}

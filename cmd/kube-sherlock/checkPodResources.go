package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewCheckPodResources() *CheckPodResources {
	cmd := &CheckPodResources{
		fs: flag.NewFlagSet("resources", flag.ContinueOnError),
	}

	cmd.fs.StringVar(&cmd.namespaces, "n", "", "namespaces to check")
	cmd.fs.BoolVar(&cmd.table, "t", false, "Set true for output table")
	cmd.fs.BoolVar(&cmd.useKubeConfig, "kubeconfig", false, "use local kubeconfig")

	return cmd
}

type CheckPodResources struct {
	fs *flag.FlagSet

	namespaces    string
	table         bool
	useKubeConfig bool
}

func (g *CheckPodResources) Name() string {
	return g.fs.Name()
}

func (g *CheckPodResources) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *CheckPodResources) Run(config *rest.Config) error {
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	resources := []podResources{}

	var ns = strings.Split(g.namespaces, ",")

	for _, namespace := range ns {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}

		for _, pod := range pods.Items {

			result := podResources{
				Namespace: pod.Namespace,
				PodName:   pod.Name,
			}

			for _, container := range pod.Spec.Containers {
				result.ContainerName = container.Name
				result.cpuRequests = container.Resources.Requests.Cpu().String()
				result.cpuLimits = container.Resources.Limits.Cpu().String()
				result.memoryRequests = container.Resources.Requests.Memory().String()
				result.memoryLimits = container.Resources.Limits.Memory().String()
			}

			resources = append(resources, result)
		}
	}

	if !g.table {
		g.renderJson(resources)
	} else {
		g.renderPodResourcesResultsTable(resources)
	}

	return nil
}

func (g *CheckPodResources) UseKubeConfig() bool {
	return g.useKubeConfig
}

func (g *CheckPodResources) renderPodResourcesResultsTable(resources []podResources) {
	resultsTable := tablewriter.NewWriter(os.Stdout)
	for _, result := range resources {
		resultsTable.SetHeader([]string{"Namespace", "Pod Name", "Container Name", "cpu Requests", "cpu Limits", "memory Requests", "memory Limits"})
		resultsTable.SetAutoMergeCells(true)
		resultsTable.SetRowLine(true)

		resultsTable.Append([]string{result.Namespace, result.PodName, result.ContainerName, result.cpuRequests, result.cpuLimits, result.memoryRequests, result.memoryLimits})
	}

	resultsTable.Render()
}

func (g *CheckPodResources) renderJson(resources []podResources) {
	j, _ := json.MarshalIndent(resources, "", "  ")
	fmt.Println(string(j))
}

type podResources struct {
	Namespace      string
	PodName        string
	ContainerName  string
	cpuRequests    string
	cpuLimits      string
	memoryRequests string
	memoryLimits   string
}

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
)

func newCheckPodResourcesCommand() *checkPodResourcesCommand {
	cmd := &checkPodResourcesCommand{}
	cmd.fs = flag.NewFlagSet("resources", flag.ContinueOnError)
	cmd.fs.StringVar(&cmd.namespaces, "n", "", "namespaces to check")
	cmd.addCommonFlags()

	return cmd
}

type checkPodResourcesCommand struct {
	namespaces string
	config
}

func (g *checkPodResourcesCommand) Name() string {
	return g.fs.Name()
}

func (g *checkPodResourcesCommand) Run(args []string) error {
	g.parse(args)

	clientset := g.createKubernetesClient()

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
				result.CPURequests = container.Resources.Requests.Cpu().String()
				result.CPULimits = container.Resources.Limits.Cpu().String()
				result.MemoryRequests = container.Resources.Requests.Memory().String()
				result.MemoryLimits = container.Resources.Limits.Memory().String()
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

func (g *checkPodResourcesCommand) renderPodResourcesResultsTable(resources []podResources) {
	resultsTable := tablewriter.NewWriter(os.Stdout)
	for _, result := range resources {
		resultsTable.SetHeader([]string{"Namespace", "Pod Name", "Container Name", "cpu Requests", "cpu Limits", "memory Requests", "memory Limits"})
		resultsTable.SetAutoMergeCells(true)
		resultsTable.SetRowLine(true)

		resultsTable.Append([]string{result.Namespace, result.PodName, result.ContainerName, result.CPURequests, result.CPULimits, result.MemoryRequests, result.MemoryLimits})
	}

	resultsTable.Render()
}

func (g *checkPodResourcesCommand) renderJson(resources []podResources) {
	j, _ := json.MarshalIndent(resources, "", "  ")
	fmt.Println(string(j))
}

type podResources struct {
	Namespace      string
	PodName        string
	ContainerName  string
	CPURequests    string
	CPULimits      string
	MemoryRequests string
	MemoryLimits   string
}

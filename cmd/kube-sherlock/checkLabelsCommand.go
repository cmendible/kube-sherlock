package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewCheckLabelsCommand() *CheckLabelsCommand {
	cmd := &CheckLabelsCommand{
		fs: flag.NewFlagSet("labels", flag.ContinueOnError),
	}

	cmd.fs.StringVar(&cmd.labels, "l", "", "labels to check")
	cmd.fs.StringVar(&cmd.namespaces, "n", "", "namespaces to check")
	cmd.fs.BoolVar(&cmd.table, "t", false, "Set true for output table")
	cmd.fs.BoolVar(&cmd.useKubeConfig, "kubeconfig", false, "use local kubeconfig")

	return cmd
}

type CheckLabelsCommand struct {
	fs *flag.FlagSet

	labels        string
	namespaces    string
	table         bool
	useKubeConfig bool
}

func (g *CheckLabelsCommand) Name() string {
	return g.fs.Name()
}

func (g *CheckLabelsCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *CheckLabelsCommand) Run(config *rest.Config) error {
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// read kube-sherlock configuration
	var c sherlockConfig
	c.getSherlockConfig(g.labels, g.namespaces)

	podResults := make(map[string][]*podResult)

	for _, namespace := range c.Namespaces {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
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

	if !g.table {
		g.renderJson(podResults)
	} else {
		g.renderLabelsCommandResultsTable(podResults)
	}

	return nil
}

func (g *CheckLabelsCommand) UseKubeConfig() bool {
	return g.useKubeConfig
}

func (g *CheckLabelsCommand) renderLabelsCommandResultsTable(podResults map[string][]*podResult) {
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

func (g *CheckLabelsCommand) renderJson(podResults map[string][]*podResult) {
	j, _ := json.MarshalIndent(podResults, "", "  ")
	fmt.Println(string(j))
}

type podResult struct {
	Namespace string
	PodName   string
}

func (c *sherlockConfig) getSherlockConfig(labels string, namespaces string) *sherlockConfig {
	if labels != "" {
		c.Labels = strings.Split(labels, ",")
		c.Namespaces = strings.Split(namespaces, ",")
		return c
	}

	if os.Getenv("KS_LABELS") != "" {
		c.Labels = strings.Split(os.Getenv("KS_LABELS"), ",")
		c.Namespaces = strings.Split(os.Getenv("KS_NAMESPACES"), ",")
		return c
	}

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

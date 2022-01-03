package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type config struct {
	fs            *flag.FlagSet
	table         bool
	useKubeConfig bool
}

func (c *config) parse(args []string) {
	e := c.fs.Parse(args)
	if e == flag.ErrHelp {
		os.Exit(0)
	}
}

func (c *config) addCommonFlags() {
	c.fs.BoolVar(&c.table, "t", false, "Set true for output table")
	c.fs.BoolVar(&c.useKubeConfig, "kubeconfig", false, "use local kubeconfig")
}

func (c *config) createKubernetesClient() *kubernetes.Clientset {
	config := c.getKubeConfig()

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}

func (c *config) getKubeConfig() *rest.Config {
	useKubeConfig := c.useKubeConfig || !c.inCluster()

	kubeconfig := ""
	if home := c.homeDir(); useKubeConfig && home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func (c *config) inCluster() bool {
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	return len(host) > 0 && len(port) > 0
}

func (c *config) homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

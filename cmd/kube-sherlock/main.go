package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		err := errors.New("you must pass a sub-command")
		fmt.Println(err)
		os.Exit(1)
	}

	cmds := []Runner{
		NewCheckLabelsCommand(),
		NewCheckPodResources(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			// get K8s Configuration
			config := getKubeConfig(cmd.UseKubeConfig())
			cmd.Run(config)
			os.Exit(0)
		}
	}

	err := fmt.Errorf("unknown subcommand: %s", subcommand)
	fmt.Println(err)
	os.Exit(1)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getKubeConfig(useKubeConfig bool) (config *rest.Config) {
	var kubeconfig string
	if useKubeConfig {
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

package main

import "k8s.io/client-go/rest"

type Runner interface {
	Init([]string) error
	Run(config *rest.Config) error
	Name() string
	UseKubeConfig() bool
}

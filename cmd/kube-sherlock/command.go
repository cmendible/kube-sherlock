package main

type command interface {
	Name() string
	Run([]string) error
}

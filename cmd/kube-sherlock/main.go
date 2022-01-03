package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		err := errors.New("you must pass a sub-command")
		fmt.Println(err)
		os.Exit(1)
	}

	cmds := []command{
		newCheckLabelsCommand(),
		newCheckPodResourcesCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Run(os.Args[2:])
			os.Exit(0)
		}
	}

	err := fmt.Errorf("unknown subcommand: %s", subcommand)
	fmt.Println(err)
	os.Exit(1)
}


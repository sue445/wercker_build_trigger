package main

import (
	"github.com/mitchellh/cli"
	"log"
	"os"
)

var (
	Version  string
	Revision string
)

func main() {
	c := cli.NewCLI("werker_build_trigger", Version)
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"search": func() (cli.Command, error) {
			return NewSearch(), nil
		},
		"trigger": func() (cli.Command, error) {
			return NewTrigger(), nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

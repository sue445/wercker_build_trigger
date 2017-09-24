package main

import (
	"flag"
	"fmt"
)

const (
	DEFAULT_BRANCH = "master"
)

type Trigger struct {
	flags      *flag.FlagSet
	configFile string
	token      string
	arguments  []string
}

func NewTrigger() *Trigger {
	t := new(Trigger)

	t.flags = flag.NewFlagSet("trigger", flag.ContinueOnError)
	t.flags.Usage = func() {}
	t.flags.StringVar(&t.configFile, "config", "", "Path to config file")
	t.flags.StringVar(&t.token, "token", "", "API token")

	return t
}

func (t *Trigger) Help() string {
	t.flags.PrintDefaults()

	return "\nRun wercker build"
}

func (t *Trigger) Synopsis() string {
	return "Run wercker build"
}

func (t *Trigger) Run(args []string) int {
	if err := t.flags.Parse(args[0:]); err != nil {
		return 1
	}
	for t.flags.NArg() > 0 {
		t.arguments = append(t.arguments, t.flags.Arg(0))
		t.flags.Parse(t.flags.Args()[1:])
	}

	if len(t.configFile) == 0 || len(t.token) == 0 {
		t.flags.PrintDefaults()
		return 0
	}

	config, err := LoadConfigFromFile(t.configFile)

	if err != nil {
		panic(err)
		return 1
	}

	wercker := NewWercker(t.token)

	for _, run := range config.Pipelines {
		if len(run.Branch) == 0 {
			run.Branch = DEFAULT_BRANCH
		}

		ret, err := wercker.TriggerNewRun(run.Id, run.Branch)

		if err == nil {
			fmt.Printf("[%s:%s] Triggered run: %s\n", run.Id, run.Branch, ret.Url)
		} else {
			fmt.Printf("[%s:%s] Error: %v\n", run.Id, run.Branch, err)
		}
	}

	return 0
}

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

	for _, configPipeline := range config.Pipelines {
		if len(configPipeline.Branch) == 0 {
			configPipeline.Branch = DEFAULT_BRANCH
		}

		if len(configPipeline.Id) > 0 {
			ret, err := wercker.TriggerNewRun(configPipeline.Id, configPipeline.Branch)

			if err == nil {
				fmt.Printf("[pipelineId:%s][branch:%s] Triggered pipeline: %s\n", configPipeline.Id, configPipeline.Branch, ret.Url)
			} else {
				fmt.Printf("[pipelineId:%s][branch:%s] Error: %v\n", configPipeline.Id, configPipeline.Branch, err)
			}

		} else if len(configPipeline.Path) > 0 && len(configPipeline.Name) > 0 {
			pipeline, err := wercker.FindPipeline(configPipeline.Path, configPipeline.Name)

			if err != nil {
				fmt.Printf("[path:%s][pipelineName:%s] Error: %v\n", configPipeline.Path, configPipeline.Name, err)
			}

			ret, err := wercker.TriggerNewRun(pipeline.Id, configPipeline.Branch)
			if err == nil {
				fmt.Printf("[path:%s][pipelineName:%s][pipelineId:%s][branch:%s] Triggered pipeline: %s\n", configPipeline.Path, configPipeline.Name, pipeline.Id, configPipeline.Branch, ret.Url)
			} else {
				fmt.Printf("[path:%s][pipelineName:%s][pipelineId:%s][branch:%s] Error: %v\n", configPipeline.Path, configPipeline.Name, pipeline.Id, configPipeline.Branch, err)
			}
		}
	}

	return 0
}

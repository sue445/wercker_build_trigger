package main

import (
	"flag"
	"fmt"
)

const (
	DEFAULT_BRANCH = "master"
)

var (
	Version        string
	Revision       string
	configFile     string
	token          string
	isPrintVersion bool
)

func init() {
	flag.StringVar(&configFile, "config", "", "Path to config file")
	flag.StringVar(&token, "token", "", "API token")
	flag.BoolVar(&isPrintVersion, "version", false, "Whether showing version")

	flag.Parse()
}

func main() {
	if isPrintVersion {
		printVersion()
		return
	}

	if len(configFile) == 0 || len(token) == 0 {
		flag.PrintDefaults()
		return
	}

	config, err := LoadConfigFromFile(configFile)

	if err != nil {
		panic(err)
	}

	wercker := NewWercker(token)

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
}

func printVersion() {
	fmt.Printf("wercker_build_trigger v%s, build %s\n", Version, Revision)
}

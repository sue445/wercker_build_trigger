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

		pipeline, err := wercker.FindPipeline(configPipeline.ApplicationPath, configPipeline.PipelineName)

		if err != nil {
			fmt.Printf("[application_path:%s][pipeline_name:%s] Error: %v\n", configPipeline.ApplicationPath, configPipeline.PipelineName, err)
		}

		ret, err := wercker.TriggerNewRun(pipeline.Id, configPipeline.Branch)
		if err == nil {
			fmt.Printf("[application_path:%s][pipeline_name:%s][branch:%s] Triggered pipeline: %s\n", configPipeline.ApplicationPath, configPipeline.PipelineName, configPipeline.Branch, ret.Url)
		} else {
			fmt.Printf("[application_path:%s][pipeline_name:%s][branch:%s] Error: %v\n", configPipeline.ApplicationPath, configPipeline.PipelineName, configPipeline.Branch, err)
		}
	}
}

func printVersion() {
	fmt.Printf("wercker_build_trigger v%s, build %s\n", Version, Revision)
}

package main

import (
	"flag"
	"fmt"
)

const (
	// DefaultBranch represents default branch name when branch is undefined in config file
	DefaultBranch = "master"

	// DefaultPipelineName represents default pipeline name when pipeline_name is undefined in config file
	DefaultPipelineName = "build"
)

var (
	// Version represents app version (injected from ldflags)
	Version string

	// Revision represents app revision (injected from ldflags)
	Revision string

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
		run, err := perform(wercker, &configPipeline)

		if err == nil {
			fmt.Printf("[application_path:%s][pipeline_name:%s][branch:%s] Triggered pipeline: %s\n", configPipeline.ApplicationPath, configPipeline.PipelineName, configPipeline.Branch, run.Url)
		} else {
			fmt.Printf("[application_path:%s][pipeline_name:%s] Error: %v\n", configPipeline.ApplicationPath, configPipeline.PipelineName, err)
		}
	}
}

func printVersion() {
	fmt.Printf("wercker_build_trigger v%s, build %s\n", Version, Revision)
}

func perform(wercker WerckerTrigger, configPipeline *ConfigPipeline) (run *WerckerRun, err error) {
	if len(configPipeline.Branch) == 0 {
		configPipeline.Branch = DefaultBranch
	}
	if len(configPipeline.PipelineName) == 0 {
		configPipeline.PipelineName = DefaultPipelineName
	}

	pipeline, err := wercker.FindPipeline(configPipeline.ApplicationPath, configPipeline.PipelineName)

	if err != nil {
		return nil, err
	}

	ret, err := wercker.TriggerNewRun(pipeline.Id, configPipeline.Branch)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

package main

import (
	"flag"
	"fmt"
)

const (
	DEFAULT_BRANCH = "master"
)

var (
	Version  string
	Revision string
)
var configFile, token string
var isPrintVersion bool

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
}

func printVersion() {
	fmt.Printf("zatsu_monitor v%s, build %s\n", Version, Revision)
}

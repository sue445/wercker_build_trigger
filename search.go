package main

import (
	"flag"
	"fmt"
)

type Search struct {
	flags        *flag.FlagSet
	token        string
	appPath      string
	pipelineName string
	arguments    []string
}

func NewSearch() *Search {
	s := new(Search)

	s.flags = flag.NewFlagSet("trigger", flag.ContinueOnError)
	s.flags.Usage = func() {}
	s.flags.StringVar(&s.token, "token", "", "API token")
	s.flags.StringVar(&s.appPath, "path", "", "Application path (e.g. wercker/docs)")
	s.flags.StringVar(&s.pipelineName, "pipeline", "", "Pipeline name (e.g. build)")

	return s
}

func (s *Search) Help() string {
	return "Search pipeline id"
}

func (s *Search) Synopsis() string {
	return "Search pipeline id"
}

func (s *Search) Run(args []string) int {
	if err := s.flags.Parse(args[0:]); err != nil {
		return 1
	}
	for s.flags.NArg() > 0 {
		s.arguments = append(s.arguments, s.flags.Arg(0))
		s.flags.Parse(s.flags.Args()[1:])
	}

	if len(s.token) == 0 || len(s.appPath) == 0 || len(s.pipelineName) == 0 {
		s.flags.PrintDefaults()
		return 0
	}

	wercker := NewWercker(s.token)
	pipeline, err := wercker.FindPipeline(s.appPath, s.pipelineName)
	if err != nil {
		fmt.Sprintln("%s", err)
		return 1
	}

	fmt.Printf("PipelineId=%s", pipeline.Id)

	return 0
}

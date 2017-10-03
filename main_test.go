package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type StubWercker struct {
}

func NewStubWercker() *StubWercker {
	return new(StubWercker)
}

var pipelineCount int

func (w *StubWercker) FindPipeline(appPath string, pipelineName string) (pipeline *WerckerPipeline, err error) {
	pipelineCount++

	switch pipelineCount {
	case 1:
		pipeline = new(WerckerPipeline)
		pipeline.Id = "abcdefg"
		pipeline.Name = "build"
		return pipeline, nil

	default:
		err := fmt.Errorf("NotFound %s, %s", appPath, pipelineName)
		return nil, err
	}
}

func (w *StubWercker) TriggerNewRun(pipelineId string, branch string) (run *WerckerRun, err error) {
	run = new(WerckerRun)
	run.Url = "https://app.wercker.com/api/v3/runs/588a61d30a002301003b44d5"
	return run, nil

}

func TestPerform(t *testing.T) {
	yamlData := `# test yaml
pipelines:
  - application_path: "wercker/docs"
    pipeline_name: "build"
    branch: "master"
  - application_path: "sue445/xxxxxxxxxx"
    pipeline_name: "xxxxxxxx"`

	config, err := LoadConfigFromData(yamlData)
	assert.NoError(t, err)

	wercker := NewStubWercker()

	pipelineCount = 0
	perform(wercker, config)
}

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

func TestPerform_Success(t *testing.T) {
	configPipeline := ConfigPipeline{
		ApplicationPath: "wercker/docs",
		PipelineName:    "build",
		Branch:          "master",
	}

	wercker := NewStubWercker()

	run, err := perform(wercker, configPipeline)

	assert.NoError(t, err)
	assert.Equal(t, "https://app.wercker.com/api/v3/runs/588a61d30a002301003b44d5", run.Url)
}

func TestPerform_Error(t *testing.T) {
	configPipeline := ConfigPipeline{
		ApplicationPath: "sue445/xxxxxxxxxx",
		PipelineName:    "build",
	}

	wercker := NewStubWercker()

	run, err := perform(wercker, configPipeline)

	assert.Error(t, err)
	assert.Nil(t, run)
}

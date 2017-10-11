package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockWercker struct {
	findPipeline  func(_appPath string, _pipelineName string) (pipeline *WerckerPipeline, err error)
	triggerNewRun func(_pipelineId string, _branch string) (run *WerckerRun, err error)
}

func NewStubWercker() *MockWercker {
	return new(MockWercker)
}

func (w *MockWercker) FindPipeline(appPath string, pipelineName string) (pipeline *WerckerPipeline, err error) {
	return w.findPipeline(appPath, pipelineName)
}

func (w *MockWercker) TriggerNewRun(pipelineID string, branch string) (run *WerckerRun, err error) {
	return w.triggerNewRun(pipelineID, branch)
}

func TestPerform_Success_MaxiumKeys(t *testing.T) {
	appPath := "wercker/docs"
	pipelineName := "build-wip"
	branch := "develop"
	url := "https://app.wercker.com/api/v3/runs/000000000000000000"
	pipelineId := "xxxxxxxxxxxxxxxxxxxxxx"

	configPipeline := ConfigPipeline{
		ApplicationPath: appPath,
		PipelineName:    pipelineName,
		Branch:          branch,
	}

	wercker := NewStubWercker()
	wercker.findPipeline = func(_appPath string, _pipelineName string) (pipeline *WerckerPipeline, err error) {
		assert.Equal(t, appPath, _appPath)
		assert.Equal(t, pipelineName, _pipelineName)

		pipeline = new(WerckerPipeline)
		pipeline.Id = pipelineId
		pipeline.Name = pipelineName
		return pipeline, nil
	}
	wercker.triggerNewRun = func(_pipelineId string, _branch string) (run *WerckerRun, err error) {
		assert.Equal(t, pipelineId, _pipelineId)
		assert.Equal(t, branch, _branch)

		run = new(WerckerRun)
		run.Url = url
		return run, nil
	}

	run, err := perform(wercker, &configPipeline)

	assert.NoError(t, err)
	assert.Equal(t, url, run.Url)
}

func TestPerform_Success_MinimumKeys(t *testing.T) {
	appPath := "wercker/docs"
	url := "https://app.wercker.com/api/v3/runs/000000000000000000"
	pipelineId := "xxxxxxxxxxxxxxxxxxxxxx"

	configPipeline := ConfigPipeline{
		ApplicationPath: appPath,
	}

	wercker := NewStubWercker()
	wercker.findPipeline = func(_appPath string, _pipelineName string) (pipeline *WerckerPipeline, err error) {
		assert.Equal(t, appPath, _appPath)
		assert.Equal(t, DefaultPipelineName, _pipelineName)

		pipeline = new(WerckerPipeline)
		pipeline.Id = pipelineId
		pipeline.Name = DefaultPipelineName
		return pipeline, nil
	}
	wercker.triggerNewRun = func(_pipelineId string, _branch string) (run *WerckerRun, err error) {
		assert.Equal(t, pipelineId, _pipelineId)
		assert.Equal(t, DefaultBranch, _branch)

		run = new(WerckerRun)
		run.Url = url
		return run, nil
	}

	run, err := perform(wercker, &configPipeline)

	assert.NoError(t, err)
	assert.Equal(t, url, run.Url)

	assert.Equal(t, DefaultBranch, configPipeline.Branch)
	assert.Equal(t, DefaultPipelineName, configPipeline.PipelineName)
}

func TestPerform_Error(t *testing.T) {
	appPath := "wercker/docs"
	pipelineName := "build"

	configPipeline := ConfigPipeline{
		ApplicationPath: appPath,
		PipelineName:    pipelineName,
	}

	wercker := NewStubWercker()
	wercker.findPipeline = func(_appPath string, _pipelineName string) (pipeline *WerckerPipeline, err error) {
		assert.Equal(t, appPath, _appPath)
		assert.Equal(t, pipelineName, _pipelineName)

		return nil, errors.New("NotFound")
	}

	run, err := perform(wercker, &configPipeline)

	assert.Error(t, err)
	assert.Nil(t, run)
}

package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigFromData(t *testing.T) {
	yamlData := `# test yaml
pipelines:
  - application_path: "wercker/docs"
    pipeline_name: "build"
    branch: "master"
  - application_path: "sue445/itamae-plugin-recipe-consul"
    pipeline_name: "build-centos70"`

	config, err := LoadConfigFromData(yamlData)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(config.Pipelines))

	assert.Equal(t, "wercker/docs", config.Pipelines[0].ApplicationPath)
	assert.Equal(t, "build", config.Pipelines[0].PipelineName)
	assert.Equal(t, "master", config.Pipelines[0].Branch)

	assert.Equal(t, "sue445/itamae-plugin-recipe-consul", config.Pipelines[1].ApplicationPath)
	assert.Equal(t, "build-centos70", config.Pipelines[1].PipelineName)
	assert.Equal(t, "", config.Pipelines[1].Branch)
}

func TestLoadConfigFromFile(t *testing.T) {
	config, err := LoadConfigFromFile("test/config.yml")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(config.Pipelines))

	assert.Equal(t, "wercker/docs", config.Pipelines[0].ApplicationPath)
	assert.Equal(t, "build", config.Pipelines[0].PipelineName)
	assert.Equal(t, "master", config.Pipelines[0].Branch)

	assert.Equal(t, "sue445/itamae-plugin-recipe-consul", config.Pipelines[1].ApplicationPath)
	assert.Equal(t, "build-centos70", config.Pipelines[1].PipelineName)
	assert.Equal(t, "", config.Pipelines[1].Branch)
}

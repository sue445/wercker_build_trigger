package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigFromData(t *testing.T) {
	yamlData := `# test yaml
pipelines:
  - path: "wercker/docs"
    name: "build"
    branch: "master"
  - path: "sue445/itamae-plugin-recipe-consul"
    name: "build-centos70"`

	config, err := LoadConfigFromData(yamlData)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(config.Pipelines))

	assert.Equal(t, "wercker/docs", config.Pipelines[0].Path)
	assert.Equal(t, "build", config.Pipelines[0].Name)
	assert.Equal(t, "master", config.Pipelines[0].Branch)

	assert.Equal(t, "sue445/itamae-plugin-recipe-consul", config.Pipelines[1].Path)
	assert.Equal(t, "build-centos70", config.Pipelines[1].Name)
	assert.Equal(t, "", config.Pipelines[1].Branch)
}

func TestLoadConfigFromFile(t *testing.T) {
	config, err := LoadConfigFromFile("test/config.yml")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(config.Pipelines))

	assert.Equal(t, "wercker/docs", config.Pipelines[0].Path)
	assert.Equal(t, "build", config.Pipelines[0].Name)
	assert.Equal(t, "master", config.Pipelines[0].Branch)

	assert.Equal(t, "sue445/itamae-plugin-recipe-consul", config.Pipelines[1].Path)
	assert.Equal(t, "build-centos70", config.Pipelines[1].Name)
	assert.Equal(t, "", config.Pipelines[1].Branch)
}

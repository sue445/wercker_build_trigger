package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigFromData(t *testing.T) {
	yamlData := `# test yaml
pipelines:
  - id: "123456789012345678901234"
    branch: "master"
  - id: "abcdefabcdefabcdefabcdef"`

	config, err := LoadConfigFromData(yamlData)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(config.Pipelines))

	assert.Equal(t, "123456789012345678901234", config.Pipelines[0].Id)
	assert.Equal(t, "master", config.Pipelines[0].Branch)

	assert.Equal(t, "abcdefabcdefabcdefabcdef", config.Pipelines[1].Id)
	assert.Equal(t, "", config.Pipelines[1].Branch)
}

func TestLoadConfigFromFile(t *testing.T) {
	config, err := LoadConfigFromFile("test/config.yml")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(config.Pipelines))

	assert.Equal(t, "123456789012345678901234", config.Pipelines[0].Id)
	assert.Equal(t, "master", config.Pipelines[0].Branch)

	assert.Equal(t, "abcdefabcdefabcdefabcdef", config.Pipelines[1].Id)
	assert.Equal(t, "", config.Pipelines[1].Branch)
}

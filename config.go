package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config represents config file
type Config struct {
	Pipelines []ConfigPipeline `yaml:"pipelines"`
}

// ConfigPipeline represents pipelines element of config file
type ConfigPipeline struct {
	ApplicationPath string `yaml:"application_path"`
	PipelineName    string `yaml:"pipeline_name"`
	Branch          string `yaml:"branch"`
}

// LoadConfigFromData load config from yaml data
func LoadConfigFromData(yamlData string) (Config, error) {
	c := Config{}

	err := yaml.Unmarshal([]byte(yamlData), &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// LoadConfigFromFile load config from yaml file
func LoadConfigFromFile(yamlFile string) (Config, error) {
	buf, err := ioutil.ReadFile(yamlFile)

	if err != nil {
		return Config{}, err
	}

	return LoadConfigFromData(string(buf))
}

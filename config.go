package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Pipelines []ConfigPipeline `yaml:"pipelines"`
}

type ConfigPipeline struct {
	ApplicationPath string `yaml:"application_path"`
	PipelineName    string `yaml:"pipeline_name"`
	Branch          string `yaml:"branch"`
}

func LoadConfigFromData(yamlData string) (Config, error) {
	c := Config{}

	err := yaml.Unmarshal([]byte(yamlData), &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func LoadConfigFromFile(yamlFile string) (Config, error) {
	buf, err := ioutil.ReadFile(yamlFile)

	if err != nil {
		return Config{}, err
	}

	return LoadConfigFromData(string(buf))
}

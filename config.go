package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Pipelines []struct {
		Id     string
		Branch string
	} `yaml:"pipelines"`
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

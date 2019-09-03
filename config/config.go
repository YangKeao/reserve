package config

import (
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Kube KubeConfig `yaml:"kube"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type KubeConfig struct {
	ConfigPath string `yaml:"config_path"`
}

func ReadAndParse(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &config, nil
}

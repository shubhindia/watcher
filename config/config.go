package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ResourceConfig struct {
	Kind string `yaml:"Kind"`
	Name string `yaml:"Name"`
}

type Config struct {
	Namespace string           `yaml:"namespace"`
	Newest    []ResourceConfig `yaml:"newest"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

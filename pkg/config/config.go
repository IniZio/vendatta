package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Name     string         `yaml:"name"`
	Provider string         `yaml:"provider"`
	Services map[string]int `yaml:"services"`
	Docker   struct {
		Image string   `yaml:"image"`
		Ports []string `yaml:"ports"`
		DinD  bool     `yaml:"dind"`
	} `yaml:"docker"`
	LXC struct {
		Image string `yaml:"image"`
	} `yaml:"lxc"`
	Agent struct {
		Role     string   `yaml:"role"`
		Features []string `yaml:"features"`
	} `yaml:"agent"`
	Hooks struct {
		Setup    string `yaml:"setup"`
		Dev      string `yaml:"dev"`
		Teardown string `yaml:"teardown"`
	} `yaml:"hooks"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

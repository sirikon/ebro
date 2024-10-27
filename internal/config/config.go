package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Module struct {
	Import map[string]struct {
		Local       string            `yaml:"local"`
		Git         string            `yaml:"git"`
		Environment map[string]string `yaml:"environment"`
		Generated   struct {
			Enabled bool `yaml:"enabled"`
			Config  struct {
				Environment map[string]string `yaml:"environment"`
			}
		} `yaml:"generated"`
	} `yaml:"import"`
	Environment map[string]string `yaml:"environment"`
	Tasks       map[string]struct {
		Requires   []string `yaml:"requires"`
		RequiredBy []string `yaml:"required_by"`
		Script     string   `yaml:"script"`
		SkipIf     string   `yaml:"skip_if"`
		Sources    []string `yaml:"sources"`
	} `yaml:"tasks"`
	Modules map[string]Module `yaml:"modules"`
}

func ParseFile(filePath string) (*Module, error) {
	module := Module{}
	body, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(body, &module)
	if err != nil {
		return nil, err
	}
	return &module, nil
}

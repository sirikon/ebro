package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func parseModuleFromFile(path string) (*Module, error) {
	module := Module{}
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(body, &module)
	if err != nil {
		return nil, err
	}
	err = module.Validate()
	if err != nil {
		return nil, fmt.Errorf("parsing %v: %w", path, err)
	}
	return &module, nil
}

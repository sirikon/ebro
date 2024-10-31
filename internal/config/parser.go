package config

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

func parseModuleFromFile(filePath string) (*Module, error) {
	module := Module{}

	body, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("parsing %v: %w", filePath, err)
	}
	err = yaml.Unmarshal(body, &module)
	if err != nil {
		return nil, fmt.Errorf("parsing %v: %w", filePath, err)
	}
	err = module.Validate()
	if err != nil {
		return nil, fmt.Errorf("parsing %v: %w", filePath, err)
	}

	for name, config := range module.Imports {
		_, ok := module.Modules[name]
		if ok {
			return nil, fmt.Errorf("parsing %v: trying to import module %v, but it already exists", filePath, name)
		}
		submodule, err := parseModuleFromFile(path.Join(config.From, "Ebro.yaml"))
		if err != nil {
			return nil, fmt.Errorf("parsing %v: %w", filePath, err)
		}
		if module.Modules == nil {
			module.Modules = make(map[string]Module)
		}
		module.Modules[name] = *submodule
	}

	return &module, nil
}

package config

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

func parseModuleFromFile(filePath string) (*Module, error) {
	moduleFile := ModuleFile{}
	module := Module{}

	body, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("parsing %v: %w", filePath, err)
	}
	err = yaml.Unmarshal(body, &moduleFile)
	if err != nil {
		return nil, fmt.Errorf("parsing %v: %w", filePath, err)
	}
	err = moduleFile.Validate()
	if err != nil {
		return nil, fmt.Errorf("parsing %v: %w", filePath, err)
	}

	working_directory := path.Dir(filePath)
	module.WorkingDirectory = &working_directory
	module.Environment = moduleFile.Environment
	module.Tasks = moduleFile.Tasks
	module.Modules = moduleFile.Modules

	for name, config := range moduleFile.Imports {
		_, ok := module.Modules[name]
		if ok {
			return nil, fmt.Errorf("parsing %v: trying to import module %v, but it already exists", filePath, name)
		}
		submodule, err := parseModuleFromFile(path.Join(path.Dir(filePath), config.From, "Ebro.yaml"))
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

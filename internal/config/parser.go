package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

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

	if moduleFile.WorkingDirectory == nil {
		working_directory, err := filepath.Abs(path.Dir(filePath))
		if err != nil {
			return nil, fmt.Errorf("parsing %v: obtaining absolute path to directory: %w", filePath, err)
		}
		module.WorkingDirectory = &working_directory
	} else {
		module.WorkingDirectory = moduleFile.WorkingDirectory
	}
	module.Environment = moduleFile.Environment
	module.Tasks = moduleFile.Tasks
	module.Modules = moduleFile.Modules

	for import_name, import_obj := range moduleFile.Imports {
		_, ok := module.Modules[import_name]
		if ok {
			return nil, fmt.Errorf("parsing %v: trying to import module %v, but it already exists", filePath, import_name)
		}
		submodule, err := parseModuleFromFile(path.Join(path.Dir(filePath), import_obj.From, "Ebro.yaml"))
		if err != nil {
			return nil, fmt.Errorf("parsing %v: %w", filePath, err)
		}
		if module.Modules == nil {
			module.Modules = make(map[string]Module)
		}
		module.Modules[import_name] = *submodule
	}

	return &module, nil
}

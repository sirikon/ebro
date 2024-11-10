package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ParseModuleFromFile(filePath string) (*Module, error) {
	return parseModuleFromFile(filePath)
}

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
		workingDirectory, err := filepath.Abs(path.Dir(filePath))
		if err != nil {
			return nil, fmt.Errorf("parsing %v: obtaining absolute path to directory: %w", filePath, err)
		}
		module.WorkingDirectory = &workingDirectory
	} else {
		module.WorkingDirectory = moduleFile.WorkingDirectory
	}
	module.Environment = moduleFile.Environment
	module.Tasks = moduleFile.Tasks
	module.Modules = moduleFile.Modules

	for importName, importObj := range moduleFile.Imports {
		_, ok := module.Modules[importName]
		if ok {
			return nil, fmt.Errorf("parsing %v: trying to import module %v, but it already exists", filePath, importName)
		}

		moduleFilePath := path.Join(path.Dir(filePath), importObj.From, "Ebro.yaml")
		gitReference, err := parseGitReference(importObj.From)
		if err != nil {
			return nil, fmt.Errorf("parsing possible git reference %v in file %v: %w", importObj.From, filePath, err)
		}

		if gitReference != nil {
			err := cloneGitReference(gitReference)
			if err != nil {
				return nil, fmt.Errorf("cloning %v: %w", importName, err)
			}
			moduleFilePath = path.Join(gitReference.clonePath, gitReference.subPath, "Ebro.yaml")
		}

		submodule, err := parseModuleFromFile(moduleFilePath)
		if err != nil {
			return nil, fmt.Errorf("parsing %v: %w", filePath, err)
		}
		if module.Modules == nil {
			module.Modules = make(map[string]Module)
		}
		module.Modules[importName] = *submodule
	}

	return &module, nil
}

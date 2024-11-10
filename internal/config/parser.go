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

		moduleFilePath, err := parseImport(path.Dir(filePath), importObj.From)
		if err != nil {
			return nil, fmt.Errorf("parsing import %v in file %v: %w", importObj.From, filePath, err)
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

func parseImport(base string, from string) (string, error) {
	parsedGitImport, err := parseGitImport(from)
	if err != nil {
		return "", fmt.Errorf("parsing possible git import %v: %w", from, err)
	}

	if parsedGitImport != nil {
		err := cloneGitImport(parsedGitImport)
		if err != nil {
			return "", fmt.Errorf("cloning git import %v: %w", from, err)
		}

		return path.Join(parsedGitImport.path, parsedGitImport.subpath, "Ebro.yaml"), nil
	}

	return path.Join(base, from, "Ebro.yaml"), nil
}

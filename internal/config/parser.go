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
		module_file_path := path.Join(path.Dir(filePath), import_obj.From, "Ebro.yaml")
		git_refernce, err := parseGitReference(import_obj.From)
		if err != nil {
			return nil, fmt.Errorf("parsing possible git reference %v in file %v: %w", import_obj.From, filePath, err)
		}

		if git_refernce != nil {
			err := gitCloneModule(git_refernce)
			if err != nil {
				return nil, fmt.Errorf("cloning %v: %w", import_name, err)
			}
			module_file_path = path.Join(git_refernce.clonePath, git_refernce.subPath, "Ebro.yaml")
		}

		submodule, err := parseModuleFromFile(module_file_path)
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

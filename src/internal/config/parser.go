package config

import (
	"fmt"
	"os"
	"path"

	"github.com/sirikon/ebro/internal/config/sources"
	"github.com/sirikon/ebro/internal/constants"
	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/utils"

	"github.com/goccy/go-yaml"
)

func ParseRootModule(modulePath string, baseEnvironment map[string]string) (*core.IndexedModule, error) {
	rootModule, err := parseModule(modulePath, []map[string]string{baseEnvironment})
	if err != nil {
		return nil, fmt.Errorf("parsing root module: %w", err)
	}

	err = ValidateModule(rootModule)
	if err != nil {
		return nil, fmt.Errorf("validating root module: %w", err)
	}

	indexedModule := NewIndexedModule(rootModule)
	PurgeModule(indexedModule)
	coreModule, err := NormalizeModule(indexedModule)
	if err != nil {
		return nil, fmt.Errorf("normalizing root module: %w", err)
	}

	return core.NewIndexedModule(coreModule), nil
}

func parseModule(modulePath string, environmentChain []map[string]string) (*Module, error) {
	module := &Module{}

	body, err := os.ReadFile(modulePath)
	if err != nil {
		return nil, fmt.Errorf("reading module file: %w", err)
	}

	err = yaml.UnmarshalWithOptions(body, module, yaml.DisallowUnknownField())
	if err != nil {
		return nil, fmt.Errorf("unmarshalling module file: %w", err)
	}

	err = processModule(module, path.Dir(modulePath), environmentChain)
	if err != nil {
		return nil, fmt.Errorf("processing module: %w", err)
	}

	return module, nil
}

func processModule(module *Module, workingDirectory string, environmentChain []map[string]string) error {
	if !path.IsAbs(module.WorkingDirectory) {
		module.WorkingDirectory = path.Join(workingDirectory, module.WorkingDirectory)
	}

	alreadyProcessedModules := make(map[string]bool)

	for importName, importObj := range module.ImportsSorted() {
		if _, ok := module.Modules[importName]; ok {
			return fmt.Errorf("cannot process import %v because there is already a module called %v", importName, importName)
		}

		importEnv, err := utils.ExpandMergeEnvs(append([]map[string]string{importObj.Environment, module.Environment}, environmentChain...)...)
		if err != nil {
			return fmt.Errorf("expanding environment for import operation: %w", err)
		}
		expandedFrom, err := utils.ExpandString(importObj.From, importEnv)
		if err != nil {
			return fmt.Errorf("expanding import.from for import operation: %w", err)
		}

		importPath, err := sourceModule(module.WorkingDirectory, expandedFrom)
		if err != nil {
			return fmt.Errorf("parsing import.from %v: %w", expandedFrom, err)
		}

		submodule, err := parseModule(path.Join(importPath, constants.DefaultFile), append([]map[string]string{module.Environment}, environmentChain...))
		if err != nil {
			return fmt.Errorf("parsing import %v: %w", importName, err)
		}

		if module.Modules == nil {
			module.Modules = make(map[string]*Module)
		}
		module.Modules[importName] = submodule
		alreadyProcessedModules[importName] = true
	}

	for submoduleName, submodule := range module.ModulesSorted() {
		if alreadyProcessedModules[submoduleName] {
			continue
		}
		processModule(submodule, module.WorkingDirectory, append([]map[string]string{module.Environment}, environmentChain...))
	}

	return nil
}

func sourceModule(base string, from string) (string, error) {
	for _, source := range sources.Sources {
		match, err := source.Match(from)
		if err != nil {
			return "", fmt.Errorf("matching source: %w", err)
		}
		if match {
			module_path, err := source.Resolve(base, from)
			if err != nil {
				return "", fmt.Errorf("resolving source: %w", err)
			}
			return module_path, nil
		}
	}

	return "", fmt.Errorf("no source matched")
}

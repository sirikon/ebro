package config

import (
	"fmt"
	"os"
	"path"

	"github.com/sirikon/ebro/internal/config/sources"
	"github.com/sirikon/ebro/internal/constants"
	"github.com/sirikon/ebro/internal/core"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

func init() {
	yaml.RegisterCustomUnmarshaler(func(env *core.Environment, b []byte) error {
		file, err := parser.ParseBytes(b, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		body := file.Docs[0].Body

		if body.Type() != ast.MappingType {
			return fmt.Errorf("wrong type")
		}

		mapping := body.(*ast.MappingNode)
		for _, mappingValue := range mapping.Values {
			if mappingValue.Key.Type() != ast.StringType {
				return fmt.Errorf("wrong type for key %v in mapping", mappingValue.Key.String())
			}
			if mappingValue.Value.Type() != ast.StringType {
				return fmt.Errorf("wrong type for value %v in mapping", mappingValue.Value.String())
			}

			key := mappingValue.Key.(*ast.StringNode)
			value := mappingValue.Value.(*ast.StringNode)

			env.Set(key.Value, value.Value)
		}
		return nil
	})

	yaml.RegisterCustomMarshaler(func(env *core.Environment) ([]byte, error) {
		return yaml.Marshal(env.Map())
	})
}

func ParseRootModule(modulePath string) (*core.IndexedModule, error) {
	rootModule, err := parseModule(modulePath)
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

func parseModule(modulePath string) (*Module, error) {
	module := &Module{}

	body, err := os.ReadFile(modulePath)
	if err != nil {
		return nil, fmt.Errorf("reading module file: %w", err)
	}

	err = yaml.UnmarshalWithOptions(body, module, yaml.DisallowUnknownField())
	if err != nil {
		return nil, fmt.Errorf("unmarshalling module file: %w", err)
	}

	err = processModule(module, path.Dir(modulePath))
	if err != nil {
		return nil, fmt.Errorf("processing module: %w", err)
	}

	return module, nil
}

func processModule(module *Module, workingDirectory string) error {
	if !path.IsAbs(module.WorkingDirectory) {
		module.WorkingDirectory = path.Join(workingDirectory, module.WorkingDirectory)
	}

	alreadyProcessedModules := make(map[string]bool)

	for importName, importObj := range module.ImportsSorted() {
		if _, ok := module.Modules[importName]; ok {
			return fmt.Errorf("cannot process import %v because there is already a module called %v", importName, importName)
		}

		importPath, err := sourceModule(module.WorkingDirectory, importObj.From)
		if err != nil {
			return fmt.Errorf("parsing import.from %v: %w", importObj.From, err)
		}

		submodule, err := parseModule(path.Join(importPath, constants.DefaultFile))
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
		processModule(submodule, module.WorkingDirectory)
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

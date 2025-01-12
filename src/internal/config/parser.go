package config

import (
	"fmt"
	"os"
	"path"

	"github.com/sirikon/ebro/internal/config/sources"
	"github.com/sirikon/ebro/internal/constants"

	"github.com/goccy/go-yaml"
)

func ParseModule(modulePath string) (*Module, error) {
	module := Module{}

	body, err := os.ReadFile(modulePath)
	if err != nil {
		return nil, fmt.Errorf("reading module file: %w", err)
	}

	err = yaml.UnmarshalWithOptions(body, &module, yaml.DisallowUnknownField())
	if err != nil {
		return nil, fmt.Errorf("unmarshalling module file: %w", err)
	}

	for importName, importObj := range module.Imports {
		if _, ok := module.Modules[importName]; ok {
			return nil, fmt.Errorf("cannot process import %v because there is already a module called %v", importName, importName)
		}

		importPath, err := SourceModule(path.Dir(modulePath), importObj.From)
		if err != nil {
			return nil, fmt.Errorf("parsing import %v: %w", importObj.From, err)
		}

		submodule, err := ParseModule(path.Join(importPath, constants.DefaultFile))
		if err != nil {
			return nil, fmt.Errorf("parsing import %v: %w", importName, err)
		}

		module.Modules[importName] = submodule
	}

	return &module, nil
}

func SourceModule(base string, from string) (string, error) {
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

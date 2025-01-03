package config

import (
	"fmt"
	"os"

	"github.com/sirikon/ebro/internal/config/sources"

	"github.com/goccy/go-yaml"
)

func ParseModule(modulePath string) (Module, error) {
	module := Module{}

	body, err := os.ReadFile(modulePath)
	if err != nil {
		return module, fmt.Errorf("reading module file: %w", err)
	}

	err = yaml.Unmarshal(body, &module)
	if err != nil {
		return module, fmt.Errorf("unmarshalling module file: %w", err)
	}

	return module, nil
}

func ImportModule(base string, from string) (string, error) {
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

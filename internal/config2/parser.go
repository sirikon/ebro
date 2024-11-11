package config2

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseModule(modulePath string) (Module, error) {
	module := Module{}

	body, err := os.ReadFile(modulePath)
	if err != nil {
		return module, fmt.Errorf("reading file %v: %w", modulePath, err)
	}

	err = yaml.Unmarshal(body, &module)
	if err != nil {
		return module, fmt.Errorf("unmarshalling file %v: %w", modulePath, err)
	}

	return module, nil
}

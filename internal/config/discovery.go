package config

import (
	"errors"
	"fmt"
	"os"
)

var ErrDiscoveryFailed = errors.New("discovery failed")
var discoveryPaths = []string{"Ebro.yaml"}

func DiscoverModule() (*Module, error) {
	for _, path := range discoveryPaths {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("discovery on file %v: %w", path, err)
		}

		module, err := parseModuleFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("discovery on file %v: %w", path, err)
		}
		return module, nil
	}
	return nil, ErrDiscoveryFailed
}

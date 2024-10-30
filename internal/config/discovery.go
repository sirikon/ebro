package config

import (
	"errors"
	"fmt"
	"os"
)

var ErrDiscoveryFailed = errors.New("discovery failed")
var discoveryPaths = []string{"Ebro.yaml", "Ebro.yml"}

func DiscoverModule() (*Module, error) {
	for _, path := range discoveryPaths {
		module, err := parseModuleFromFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("discovery on file %v: %w", path, err)
		}
		return module, nil
	}
	return nil, ErrDiscoveryFailed
}

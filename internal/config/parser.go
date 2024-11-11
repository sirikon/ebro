package config

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/sirikon/ebro/internal/config/remote"
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

func ImportModule(base string, from string) (string, error) {
	parsedGitImport, err := remote.ParseGitImport(from)
	if err != nil {
		return "", fmt.Errorf("parsing possible git import %v: %w", from, err)
	}

	if parsedGitImport != nil {
		err := remote.CloneGitImport(parsedGitImport)
		if err != nil {
			return "", fmt.Errorf("cloning git import %v: %w", from, err)
		}

		return path.Join(parsedGitImport.Path, parsedGitImport.Subpath), nil
	}

	return path.Join(base, from), nil
}

package config

import (
	"fmt"
)

type ModuleFile struct {
	Module  `yaml:",inline"`
	Imports map[string]Import `yaml:"imports,omitempty"`
}

type Module struct {
	WorkingDirectory *string           `yaml:"working_directory,omitempty"`
	Environment      map[string]string `yaml:"environment,omitempty"`
	Tasks            map[string]Task   `yaml:"tasks,omitempty"`
	Modules          map[string]Module `yaml:"modules,omitempty"`
}

type Task struct {
	WorkingDirectory *string           `yaml:"working_directory,omitempty"`
	Environment      map[string]string `yaml:"environment,omitempty"`
	Requires         []string          `yaml:"requires,omitempty"`
	RequiredBy       []string          `yaml:"required_by,omitempty"`
	Script           string            `yaml:"script,omitempty"`
	When             *When             `yaml:"when,omitempty"`
}

type When struct {
	CheckFails    string `yaml:"check_fails,omitempty"`
	OutputChanges string `yaml:"output_changes,omitempty"`
}

type Import struct {
	From string `yaml:"from,omitempty"`
}

func (mf ModuleFile) Validate() error {
	for name, task := range mf.Tasks {
		if err := task.Validate(); err != nil {
			return fmt.Errorf("validating task %v: %w", name, err)
		}
	}
	return nil
}

func (t Task) Validate() error {
	return nil
}

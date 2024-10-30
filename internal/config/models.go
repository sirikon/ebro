package config

import (
	"fmt"
)

type Module struct {
	Imports     map[string]Import
	Environment map[string]string `yaml:"environment"`
	Tasks       map[string]Task   `yaml:"tasks"`
	Modules     map[string]Module `yaml:"modules"`
}

type Task struct {
	Environment map[string]string `yaml:"environment"`
	Requires    []string          `yaml:"requires"`
	RequiredBy  []string          `yaml:"required_by"`
	Script      string            `yaml:"script"`
	SkipIf      string            `yaml:"skip_if"`
	Sources     []string          `yaml:"sources"`
}

type Import struct {
	From        string            `yaml:"from"`
	Environment map[string]string `yaml:"environment"`
	Generated   *Task             `yaml:"generated"`
}

func (m Module) Validate() error {
	for name, task := range m.Tasks {
		if err := task.Validate(); err != nil {
			return fmt.Errorf("validating task %v: %w", name, err)
		}
	}
	return nil
}

func (t Task) Validate() error {
	return nil
}

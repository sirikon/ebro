package config

import "fmt"

type Module struct {
	WorkingDirectory string            `yaml:"working_directory,omitempty"`
	Imports          map[string]Import `yaml:"imports,omitempty"`
	Environment      map[string]string `yaml:"environment,omitempty"`
	Tasks            map[string]Task   `yaml:"tasks,omitempty"`
	Modules          map[string]Module `yaml:"modules,omitempty"`
}

type Task struct {
	WorkingDirectory string            `yaml:"working_directory,omitempty"`
	Abstract         bool              `yaml:"abstract,omitempty"`
	Extends          []string          `yaml:"extends,omitempty"`
	Environment      map[string]string `yaml:"environment,omitempty"`
	Requires         []string          `yaml:"requires,omitempty"`
	RequiredBy       []string          `yaml:"required_by,omitempty"`
	Script           string            `yaml:"script,omitempty"`
	Quiet            *bool             `yaml:"quiet,omitempty"`
	When             *When             `yaml:"when,omitempty"`
}

type When struct {
	CheckFails    string `yaml:"check_fails,omitempty"`
	OutputChanges string `yaml:"output_changes,omitempty"`
}

type Import struct {
	From        string            `yaml:"from,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

func (t Task) Validate() error {
	if len(t.Requires) == 0 && t.Script == "" && len(t.Extends) == 0 && !t.Abstract {
		return fmt.Errorf("task has nothing to do (no requires, script, extends nor abstract)")
	}
	return nil
}

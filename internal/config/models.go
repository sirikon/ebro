package config

type Module struct {
	WorkingDirectory string            `yaml:"working_directory,omitempty"`
	Imports          map[string]Import `yaml:"imports,omitempty"`
	Environment      map[string]string `yaml:"environment,omitempty"`
	Tasks            map[string]Task   `yaml:"tasks,omitempty"`
}

type Task struct {
	WorkingDirectory string            `yaml:"working_directory,omitempty"`
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
	From        string            `yaml:"from,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

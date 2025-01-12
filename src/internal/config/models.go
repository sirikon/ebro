package config

type Module struct {
	WorkingDirectory string             `yaml:"working_directory,omitempty"`
	Imports          map[string]Import  `yaml:"imports,omitempty"`
	Environment      map[string]string  `yaml:"environment,omitempty"`
	Tasks            map[string]*Task   `yaml:"tasks,omitempty"`
	Modules          map[string]*Module `yaml:"modules,omitempty"`
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

func (m *Module) GetTask(taskReference TaskReference) *Task {
	if taskReference.IsRelative {
		panic("cannot call getTask with a relative taskPath")
	}

	currentModule := m
	for i, part := range taskReference.Parts {
		if i >= (len(taskReference.Parts) - 1) {
			module, ok := currentModule.Modules[part]
			if ok {
				if task, ok := module.Tasks["default"]; ok {
					return task
				}
			} else {
				if task, ok := currentModule.Tasks[part]; ok {
					return task
				}
			}
		} else {
			module, ok := currentModule.Modules[part]
			if ok {
				currentModule = module
			} else {
				return nil
			}
		}
	}

	return nil
}

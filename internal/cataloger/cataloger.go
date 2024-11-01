package cataloger

import (
	"strings"

	"github.com/sirikon/ebro/internal/config"
)

type Catalog map[string]config.Task

func MakeCatalog(module *config.Module) Catalog {
	catalog := catalogModule(module, []string{}, nil)
	for _, task := range catalog {
		NormalizeTaskReferences(catalog, task.Requires)
		NormalizeTaskReferences(catalog, task.RequiredBy)
	}
	return catalog
}

func NormalizeTaskReferences(catalog Catalog, task_names []string) {
	for i, task_name := range task_names {
		defaulted_task_name := task_name + ":default"
		_, taskExists := catalog[task_name]
		_, defaultedTaskExists := catalog[defaulted_task_name]
		if !taskExists && defaultedTaskExists {
			task_names[i] = defaulted_task_name
		}
	}
}

func catalogModule(module *config.Module, name_trail []string, working_directory *string) Catalog {
	result := make(Catalog)
	prefix := ":" + strings.Join(append(name_trail, ""), ":")
	if module.WorkingDirectory == nil {
		module.WorkingDirectory = working_directory
	}

	for task_name, task := range module.Tasks {
		for i := range task.Requires {
			task.Requires[i] = prefix + task.Requires[i]
		}
		for i := range task.RequiredBy {
			task.RequiredBy[i] = prefix + task.RequiredBy[i]
		}
		if task.WorkingDirectory == nil {
			task.WorkingDirectory = module.WorkingDirectory
		}
		result[prefix+task_name] = task
	}

	for submodule_name, submodule := range module.Modules {
		module_tasks := catalogModule(&submodule, append(name_trail, submodule_name), module.WorkingDirectory)
		for task_name, task := range module_tasks {
			result[task_name] = task
		}
	}

	return result
}

package cataloger

import (
	"strings"

	"github.com/sirikon/ebro/internal/config"
)

type Catalog map[string]config.Task

func MakeCatalog(module *config.Module) Catalog {
	catalog := catalogModule(module, []string{})
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

func catalogModule(module *config.Module, trail []string) Catalog {
	result := make(Catalog)
	prefix := ":" + strings.Join(append(trail, ""), ":")

	for task_name, task := range module.Tasks {
		for i := range task.Requires {
			task.Requires[i] = prefix + task.Requires[i]
		}
		for i := range task.RequiredBy {
			task.RequiredBy[i] = prefix + task.RequiredBy[i]
		}
		result[prefix+task_name] = task
	}

	for module_name, module := range module.Modules {
		module_tasks := catalogModule(&module, append(trail, module_name))
		for task_name, task := range module_tasks {
			result[task_name] = task
		}
	}

	return result
}

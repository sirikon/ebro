package cataloger

import (
	"path"
	"strings"

	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/utils"
)

type Catalog map[string]config.Task

func MakeCatalog(module *config.Module) Catalog {
	catalog := catalogModule(module, []string{}, nil, make(map[string]string))
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

func catalogModule(module *config.Module, name_trail []string, working_directory *string, environment map[string]string) Catalog {
	result := make(Catalog)
	prefix := ":" + strings.Join(append(name_trail, ""), ":")
	if module.WorkingDirectory == nil {
		module.WorkingDirectory = working_directory
	}
	module.Environment = utils.MergeEnv(environment, module.Environment)
	if _, ok := module.Environment["EBRO_ROOT"]; !ok {
		module.Environment["EBRO_ROOT"] = *module.WorkingDirectory
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
		task.Environment = utils.MergeEnv(module.Environment, task.Environment)
		if task.Sources != nil {
			for i, source := range task.Sources {
				task.Sources[i] = path.Join(*task.WorkingDirectory, source)
			}
		}
		result[prefix+task_name] = task
	}

	for submodule_name, submodule := range module.Modules {
		module_tasks := catalogModule(&submodule, append(name_trail, submodule_name), module.WorkingDirectory, utils.MergeEnv(module.Environment, submodule.Environment))
		for task_name, task := range module_tasks {
			result[task_name] = task
		}
	}

	return result
}

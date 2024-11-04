package cataloger

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/utils"
	"mvdan.cc/sh/v3/shell"
)

type Catalog map[string]config.Task

func MakeCatalog(module *config.Module) (Catalog, error) {
	catalog, err := catalogModule(module, []string{}, nil, make(map[string]string))
	if err != nil {
		return nil, fmt.Errorf("making catalog: %w", err)
	}
	for _, task := range catalog {
		NormalizeTaskReferences(catalog, task.Requires)
		NormalizeTaskReferences(catalog, task.RequiredBy)
	}
	return catalog, nil
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

func catalogModule(module *config.Module, name_trail []string, working_directory *string, environment map[string]string) (Catalog, error) {
	result := make(Catalog)
	prefix := ":" + strings.Join(append(name_trail, ""), ":")
	module.Environment = utils.MergeEnv(environment, module.Environment)
	if module.WorkingDirectory == nil {
		module.WorkingDirectory = working_directory
	} else {
		expanded_working_directory, err := expand(*module.WorkingDirectory, module.Environment)
		if err != nil {
			return nil, fmt.Errorf("expanding source %v: %w", module.WorkingDirectory, err)
		}
		module.WorkingDirectory = &expanded_working_directory
	}

	if _, ok := module.Environment["EBRO_ROOT"]; !ok {
		module.Environment["EBRO_ROOT"] = *module.WorkingDirectory
	}

	for task_name, task := range module.Tasks {
		for i := range task.Requires {
			if !strings.HasPrefix(task.Requires[i], ":") {
				task.Requires[i] = prefix + task.Requires[i]
			}
		}
		for i := range task.RequiredBy {
			if !strings.HasPrefix(task.RequiredBy[i], ":") {
				task.RequiredBy[i] = prefix + task.RequiredBy[i]
			}
		}
		task.Environment = utils.MergeEnv(module.Environment, task.Environment)
		if task.WorkingDirectory == nil {
			task.WorkingDirectory = module.WorkingDirectory
		} else {
			expanded_working_directory, err := expand(*task.WorkingDirectory, task.Environment)
			if err != nil {
				return nil, fmt.Errorf("expanding source %v for task %v: %w", task.WorkingDirectory, task_name, err)
			}
			task.WorkingDirectory = &expanded_working_directory
		}
		if task.Sources != nil {
			for i, source := range task.Sources {
				expanded_source, err := expand(source, task.Environment)
				if err != nil {
					return nil, fmt.Errorf("expanding source %v for task %v: %w", source, task_name, err)
				}
				if path.IsAbs(expanded_source) {
					task.Sources[i] = expanded_source
				} else {
					task.Sources[i] = path.Join(*task.WorkingDirectory, expanded_source)
				}
			}
		}
		result[prefix+task_name] = task
	}

	for submodule_name, submodule := range module.Modules {
		module_tasks, err := catalogModule(&submodule, append(name_trail, submodule_name), module.WorkingDirectory, utils.MergeEnv(module.Environment, submodule.Environment))
		if err != nil {
			return nil, err
		}
		for task_name, task := range module_tasks {
			result[task_name] = task
		}
	}

	return result, nil
}

func expand(s string, env map[string]string) (string, error) {
	return shell.Expand(s, func(s string) string {
		if val, ok := env[s]; ok {
			return val
		}
		return os.Getenv(s)
	})
}

package cataloger

import (
	"fmt"
	"os"
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

func NormalizeTaskReferences(catalog Catalog, taskNames []string) {
	for i, taskName := range taskNames {
		defaultedTaskName := taskName + ":default"
		_, taskExists := catalog[taskName]
		_, defaultedTaskExists := catalog[defaultedTaskName]
		if !taskExists && defaultedTaskExists {
			taskNames[i] = defaultedTaskName
		}
	}
}

func catalogModule(module *config.Module, nameTrail []string, workingDirectory *string, environment map[string]string) (Catalog, error) {
	result := make(Catalog)
	prefix := ":" + strings.Join(append(nameTrail, ""), ":")
	expandMap(module.Environment, environment)
	module.Environment = utils.MergeEnv(environment, module.Environment)
	if module.WorkingDirectory == nil {
		module.WorkingDirectory = workingDirectory
	} else {
		expandedWorkingDirectory, err := expand(*module.WorkingDirectory, module.Environment)
		if err != nil {
			return nil, fmt.Errorf("expanding source %v: %w", module.WorkingDirectory, err)
		}
		module.WorkingDirectory = &expandedWorkingDirectory
	}

	if _, ok := module.Environment["EBRO_ROOT"]; !ok {
		module.Environment["EBRO_ROOT"] = *module.WorkingDirectory
	}

	for taskName, task := range module.Tasks {
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
		expandMap(task.Environment, module.Environment)
		task.Environment = utils.MergeEnv(module.Environment, task.Environment)
		if task.WorkingDirectory == nil {
			task.WorkingDirectory = module.WorkingDirectory
		} else {
			expandedWorkingDirectory, err := expand(*task.WorkingDirectory, task.Environment)
			if err != nil {
				return nil, fmt.Errorf("expanding source %v for task %v: %w", task.WorkingDirectory, taskName, err)
			}
			task.WorkingDirectory = &expandedWorkingDirectory
		}
		result[prefix+taskName] = task
	}

	for submoduleName, submodule := range module.Modules {
		expandMap(submodule.Environment, module.Environment)
		moduleTasks, err := catalogModule(&submodule, append(nameTrail, submoduleName), module.WorkingDirectory, utils.MergeEnv(module.Environment, submodule.Environment))
		if err != nil {
			return nil, err
		}
		for taskName, task := range moduleTasks {
			result[taskName] = task
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

func expandMap(m map[string]string, env map[string]string) (map[string]string, error) {
	for key, value := range m {
		result, err := expand(value, env)
		if err != nil {
			return nil, fmt.Errorf("expanding %v: %w", value, err)
		}
		m[key] = result
	}
	return m, nil
}

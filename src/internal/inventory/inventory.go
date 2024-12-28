package inventory

import (
	"fmt"
	"maps"
	"os"
	"path"
	"slices"
	"strings"

	"mvdan.cc/sh/v3/shell"

	"github.com/sirikon/ebro/internal/cli"
	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/constants"
)

type Inventory map[string]config.Task

func MakeInventory(arguments cli.ExecutionArguments) (Inventory, error) {
	inv := Inventory{}
	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("obtaining working directory: %w", err)
	}

	modulePath := path.Join(workingDirectory, *arguments.GetFlagString(cli.FlagFile))
	err = processModuleFile(inv, modulePath, []string{}, map[string]string{"EBRO_ROOT": workingDirectory})
	if err != nil {
		return nil, fmt.Errorf("processing module file in %v: %w", modulePath, err)
	}

	for _, task := range inv {
		NormalizeTaskNames(inv, task.Requires)
		NormalizeTaskNames(inv, task.RequiredBy)
	}

	return inv, nil
}

func NormalizeTaskNames(inv Inventory, taskNames []string) {
	for i, taskName := range taskNames {
		defaultedTaskName := taskName + ":default"
		_, taskExists := inv[taskName]
		_, defaultedTaskExists := inv[defaultedTaskName]
		if !taskExists && defaultedTaskExists {
			taskNames[i] = defaultedTaskName
		}
	}
}

func processModuleFile(inv Inventory, modulePath string, moduleNameTrail []string, environment map[string]string) error {
	workingDirectory := path.Dir(modulePath)
	module, err := config.ParseModule(modulePath)
	if err != nil {
		return fmt.Errorf("parsing module: %w", err)
	}

	err = processModule(inv, module, moduleNameTrail, environment, workingDirectory)
	if err != nil {
		return fmt.Errorf("processing module: %w", err)
	}

	return nil
}

func processModule(inv Inventory, module config.Module, moduleNameTrail []string, environment map[string]string, workingDirectory string) error {
	taskPrefix := ":" + strings.Join(append(moduleNameTrail, ""), ":")
	makeTaskNameAbsolute := func(taskName string) string {
		if !strings.HasPrefix(taskName, ":") {
			return taskPrefix + taskName
		}
		return taskName
	}

	moduleEnvironment, err := expandMergeEnvs(module.Environment, environment)
	if err != nil {
		return fmt.Errorf("expanding module environment: %w", err)
	}
	module.Environment = moduleEnvironment

	if module.WorkingDirectory == "" {
		module.WorkingDirectory = workingDirectory
	} else if !path.IsAbs(module.WorkingDirectory) {
		module.WorkingDirectory = path.Join(workingDirectory, module.WorkingDirectory)
	}

	for taskName, task := range module.Tasks {
		taskAbsoluteName := taskPrefix + taskName
		if _, ok := inv[taskAbsoluteName]; ok {
			return fmt.Errorf("task %v (defined as %v) is already present in the inventory", taskAbsoluteName, taskName)
		}

		if err := task.Validate(); err != nil {
			return fmt.Errorf("task %v failed validation: %w", taskName, err)
		}

		task.Environment, err = expandMergeEnvs(task.Environment, module.Environment)
		if err != nil {
			return fmt.Errorf("expanding task %v environment: %w", taskName, err)
		}

		for i, t := range task.Requires {
			task.Requires[i] = makeTaskNameAbsolute(t)
		}
		for i, t := range task.RequiredBy {
			task.RequiredBy[i] = makeTaskNameAbsolute(t)
		}

		if task.WorkingDirectory == "" {
			task.WorkingDirectory = module.WorkingDirectory
		} else if !path.IsAbs(task.WorkingDirectory) {
			task.WorkingDirectory = path.Join(module.WorkingDirectory, task.WorkingDirectory)
		}

		inv[taskAbsoluteName] = task
	}

	for importName, importObj := range module.Imports {
		importEnvironment, err := expandMergeEnvs(importObj.Environment, module.Environment)
		if err != nil {
			return fmt.Errorf("expanding import %v environment: %w", importName, err)
		}

		expandedFrom, err := expandString(importObj.From, module.Environment)
		if err != nil {
			return fmt.Errorf("expanding import from %v: %w", importObj.From, err)
		}

		importModulePath, err := config.ImportModule(workingDirectory, expandedFrom)
		if err != nil {
			return fmt.Errorf("parsing import %v: %w", expandedFrom, err)
		}

		err = processModuleFile(inv, path.Join(importModulePath, constants.DefaultFile), append(moduleNameTrail, importName), importEnvironment)
		if err != nil {
			return fmt.Errorf("processing import %v: %w", importName, err)
		}
	}

	for submoduleName, submodule := range module.Modules {
		submoduleEnvironment, err := expandMergeEnvs(submodule.Environment, module.Environment)
		if err != nil {
			return fmt.Errorf("expanding module %v environment: %w", submoduleName, err)
		}

		err = processModule(inv, submodule, append(moduleNameTrail, submoduleName), submoduleEnvironment, module.WorkingDirectory)
		if err != nil {
			return fmt.Errorf("processing module %v: %w", submoduleName, err)
		}
	}

	return nil
}

func expandMergeEnvs(envs ...map[string]string) (map[string]string, error) {
	result := map[string]string{}
	for i := (len(envs) - 1); i >= 0; i-- {
		parentEnv := maps.Clone(result)
		env := envs[i]
		// We want to iterate through keys in a repeatable and predictable way.
		// The order in which we process each key SHOULD NOT BE IMPORTANT, but
		// in the scenario of a bug in here, we want the behavior to be
		// consistent.
		//
		// That's why we're sorting the keys and iterating over them
		// instead of `range`ing the map directly.
		envKeys := slices.Collect(maps.Keys(env))
		slices.Sort(envKeys)
		for _, key := range envKeys {
			if i == len(envs) {
				result[key] = env[key]
			}
			expandedValue, err := expandString(env[key], parentEnv)
			if err != nil {
				return nil, fmt.Errorf("expanding %v: %w", env[key], err)
			}
			result[key] = expandedValue
		}
	}
	return result, nil
}

func expandString(s string, env map[string]string) (string, error) {
	return shell.Expand(s, func(s string) string {
		if val, ok := env[s]; ok {
			return val
		}
		return os.Getenv(s)
	})
}

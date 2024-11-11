package inventory

import (
	"fmt"
	"os"
	"path"
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
	err = processModule(inv, modulePath, []string{}, map[string]string{"EBRO_ROOT": workingDirectory})
	if err != nil {
		return nil, fmt.Errorf("processing module in %v: %w", modulePath, err)
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

func processModule(inv Inventory, modulePath string, moduleNameTrail []string, environment map[string]string) error {
	workingDirectory := path.Dir(modulePath)
	module, err := config.ParseModule(modulePath)
	if err != nil {
		return fmt.Errorf("parsing module %v: %w", modulePath, err)
	}

	taskPrefix := ":" + strings.Join(append(moduleNameTrail, ""), ":")
	makeTaskNameAbsolute := func(taskName string) string {
		if !strings.HasPrefix(taskName, ":") {
			return taskPrefix + taskName
		}
		return taskName
	}

	moduleEnvironment, err := expandMergeEnv(module.Environment, environment)
	if err != nil {
		return fmt.Errorf("expanding module environment in %v: %w", modulePath, err)
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
			return fmt.Errorf("task %v (defined as %v in %v) is already present in the inventory", taskAbsoluteName, taskName, modulePath)
		}
		taskEnvironment, err := expandMergeEnv(task.Environment, module.Environment)
		if err != nil {
			return fmt.Errorf("expanding task %v environment in %v: %w", taskName, modulePath, err)
		}
		task.Environment = taskEnvironment
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
		importEnvironment, err := expandMergeEnv(importObj.Environment, module.Environment)
		if err != nil {
			return fmt.Errorf("expanding import %v environment in %v: %w", importName, modulePath, err)
		}
		importModulePath, err := config.ImportModule(path.Dir(modulePath), importObj.From)
		if err != nil {
			return fmt.Errorf("parsing import %v in %v: %w", importObj.From, modulePath, err)
		}
		err = processModule(inv, path.Join(importModulePath, constants.DefaultFile), append(moduleNameTrail, importName), importEnvironment)
		if err != nil {
			return fmt.Errorf("importing %v: %w", importName, err)
		}
	}

	return nil
}

func expandMergeEnv(childEnv map[string]string, parentEnv map[string]string) (map[string]string, error) {
	result := map[string]string{}
	for key, value := range parentEnv {
		result[key] = value
	}
	for key, value := range childEnv {
		expandedValue, err := expandString(value, parentEnv)
		if err != nil {
			return nil, fmt.Errorf("expanding %v: %w", value, err)
		}
		result[key] = expandedValue
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

package inventory

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/sirikon/ebro/internal/config"
)

type Inventory struct {
	Tasks           map[string]*config.Task
	taskModuleIndex map[string]*config.Module
}

func MakeInventory(module *config.Module) (Inventory, error) {
	inv := Inventory{
		Tasks:           make(map[string]*config.Task),
		taskModuleIndex: make(map[string]*config.Module),
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return inv, fmt.Errorf("obtaining working directory: %w", err)
	}

	err = processModule(inv, module, []string{}, map[string]string{"EBRO_ROOT": workingDirectory}, workingDirectory)
	if err != nil {
		return inv, fmt.Errorf("processing module: %w", err)
	}

	for _, task := range inv.Tasks {
		NormalizeTaskNames(inv, task.Requires)
		NormalizeTaskNames(inv, task.RequiredBy)
		NormalizeTaskNames(inv, task.Extends)
	}

	for taskName, task := range inv.Tasks {
		err := validateReferences(inv, task.Requires...)
		if err != nil {
			return inv, fmt.Errorf("checking references in 'requires' for task %v: %w", taskName, err)
		}
		err = validateReferences(inv, task.RequiredBy...)
		if err != nil {
			return inv, fmt.Errorf("checking references in 'required_by' for task %v: %w", taskName, err)
		}
		err = validateReferences(inv, task.Extends...)
		if err != nil {
			return inv, fmt.Errorf("checking references in 'extends' for task %v: %w", taskName, err)
		}
	}

	inheritanceOrder, err := resolveInheritanceOrder(inv)
	if err != nil {
		return inv, fmt.Errorf("resolving inheritance order: %w", err)
	}

	for _, taskName := range inheritanceOrder {
		task := inv.Tasks[taskName]
		envsToMerge := [](map[string]string){
			task.Environment,
			map[string]string{"EBRO_TASK_WORKING_DIRECTORY": task.WorkingDirectory},
		}
		parentTasks := slices.Clone(task.Extends)
		slices.Reverse(parentTasks)
		for _, parentTaskName := range parentTasks {
			parentTask := inv.Tasks[parentTaskName]
			applyInheritance(task, parentTask)
			envsToMerge = append(envsToMerge, parentTask.Environment)
		}
		envsToMerge = append(envsToMerge, inv.taskModuleIndex[taskName].Environment)
		task.Environment, err = expandMergeEnvs(envsToMerge...)
		if err != nil {
			return inv, fmt.Errorf("expanding task %v environment: %w", taskName, err)
		}
	}

	for taskName, task := range inv.Tasks {
		if task.Abstract {
			delete(inv.Tasks, taskName)
		}
	}

	return inv, nil
}

func NormalizeTaskNames(inv Inventory, taskNames []string) {
	for i, taskName := range taskNames {
		taskNames[i] = normalizeTaskName(inv, taskName)
	}
}

func normalizeTaskName(inv Inventory, taskName string) string {
	defaultedTaskName := taskName + ":default"
	_, taskExists := inv.Tasks[taskName]
	_, defaultedTaskExists := inv.Tasks[defaultedTaskName]
	if !taskExists && defaultedTaskExists {
		return defaultedTaskName
	}
	return taskName
}

func processModule(inv Inventory, module *config.Module, moduleNameTrail []string, environment map[string]string, workingDirectory string) error {
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
		if _, ok := inv.Tasks[taskAbsoluteName]; ok {
			return fmt.Errorf("task %v (defined as %v) is already present in the inventory", taskAbsoluteName, taskName)
		}

		for i, t := range task.Requires {
			task.Requires[i] = makeTaskNameAbsolute(t)
		}
		for i, t := range task.RequiredBy {
			task.RequiredBy[i] = makeTaskNameAbsolute(t)
		}
		for i, t := range task.Extends {
			task.Extends[i] = makeTaskNameAbsolute(t)
		}

		if task.WorkingDirectory == "" {
			task.WorkingDirectory = module.WorkingDirectory
		} else if !path.IsAbs(task.WorkingDirectory) {
			task.WorkingDirectory = path.Join(module.WorkingDirectory, task.WorkingDirectory)
		}

		inv.Tasks[taskAbsoluteName] = task
		inv.taskModuleIndex[taskAbsoluteName] = module
	}

	// for importName, importObj := range module.Imports {
	// 	importEnvironment, err := expandMergeEnvs(importObj.Environment, module.Environment)
	// 	if err != nil {
	// 		return fmt.Errorf("expanding import %v environment: %w", importName, err)
	// 	}

	// 	expandedFrom, err := expandString(importObj.From, module.Environment)
	// 	if err != nil {
	// 		return fmt.Errorf("expanding import from %v: %w", importObj.From, err)
	// 	}

	// 	importModulePath, err := config.SourceModule(workingDirectory, expandedFrom)
	// 	if err != nil {
	// 		return fmt.Errorf("parsing import %v: %w", expandedFrom, err)
	// 	}

	// 	err = processModuleFile(inv, path.Join(importModulePath, constants.DefaultFile), append(moduleNameTrail, importName), importEnvironment)
	// 	if err != nil {
	// 		return fmt.Errorf("processing import %v: %w", importName, err)
	// 	}
	// }

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

func validateReferences(inv Inventory, taskNames ...string) error {
	for _, taskName := range taskNames {
		if _, ok := inv.Tasks[taskName]; !ok {
			return fmt.Errorf("referenced task %v does not exist", taskName)
		}
	}
	return nil
}

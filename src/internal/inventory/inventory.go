package inventory

import (
	"fmt"
	"os"
	"path"
	"slices"

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

func processModule(inv Inventory, module *config.Module, moduleTrail []string, environment map[string]string, workingDirectory string) error {
	moduleEnvironment, err := expandMergeEnvs(module.Environment, environment)
	if err != nil {
		return fmt.Errorf("expanding module environment: %w", err)
	}
	module.Environment = moduleEnvironment

	for taskName, task := range module.Tasks {
		for i, t := range task.Requires {
			ref, _ := config.ParseTaskReference(t)
			task.Requires[i] = ref.Absolute(moduleTrail).String()
		}
		for i, t := range task.RequiredBy {
			ref, _ := config.ParseTaskReference(t)
			task.RequiredBy[i] = ref.Absolute(moduleTrail).String()
		}
		for i, t := range task.Extends {
			ref, _ := config.ParseTaskReference(t)
			task.Extends[i] = ref.Absolute(moduleTrail).String()
		}

		if task.WorkingDirectory == "" {
			task.WorkingDirectory = module.WorkingDirectory
		} else if !path.IsAbs(task.WorkingDirectory) {
			task.WorkingDirectory = path.Join(module.WorkingDirectory, task.WorkingDirectory)
		}

		taskReference := config.MakeTaskReference(append(moduleTrail, taskName))
		inv.Tasks[taskReference.PathString()] = task
		inv.taskModuleIndex[taskReference.PathString()] = module
	}

	for importName, importObj := range module.Imports {
		mergedEnv, err := expandMergeEnvs(module.Modules[importName].Environment, importObj.Environment, module.Environment)
		if err != nil {
			return fmt.Errorf("expanding import %v environment: %w", importName, err)
		}
		module.Modules[importName].Environment = mergedEnv
	}

	for submoduleName, submodule := range module.Modules {
		submoduleEnvironment, err := expandMergeEnvs(submodule.Environment, module.Environment)
		if err != nil {
			return fmt.Errorf("expanding module %v environment: %w", submoduleName, err)
		}

		err = processModule(inv, submodule, append(moduleTrail, submoduleName), submoduleEnvironment, module.WorkingDirectory)
		if err != nil {
			return fmt.Errorf("processing module %v: %w", submoduleName, err)
		}
	}

	return nil
}

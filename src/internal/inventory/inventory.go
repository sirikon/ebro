package inventory

import (
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/utils"
)

type Inventory struct {
	Tasks map[core.TaskId]*core.Task
}

type InventoryContext struct {
	inv             Inventory
	indexedModule   *core.IndexedModule
	taskModuleIndex map[core.TaskId]*core.Module
}

func MakeInventory(indexedModule *core.IndexedModule, baseEnvironment map[string]string) (Inventory, error) {
	ctx := InventoryContext{
		inv: Inventory{
			Tasks: make(map[core.TaskId]*core.Task),
		},
		indexedModule:   indexedModule,
		taskModuleIndex: make(map[core.TaskId]*core.Module),
	}

	err := ctx.processModule(indexedModule.Module, []string{}, baseEnvironment)
	if err != nil {
		return ctx.inv, fmt.Errorf("processing module: %w", err)
	}

	inheritanceOrder, err := resolveInheritanceOrder(ctx.inv)
	if err != nil {
		return ctx.inv, fmt.Errorf("resolving inheritance order: %w", err)
	}

	for _, taskId := range inheritanceOrder {
		task := ctx.inv.Tasks[taskId]
		envsToMerge := [](map[string]string){
			task.Environment,
			map[string]string{
				"EBRO_TASK_WORKING_DIRECTORY": task.WorkingDirectory,
				"EBRO_TASK_ID":                string(taskId),
				"EBRO_TASK_MODULE":            ":" + strings.Join(taskId.ModuleTrail(), ":"),
				"EBRO_TASK_NAME":              taskId.TaskName(),
			},
		}
		parentTasks := slices.Clone(task.Extends)
		slices.Reverse(parentTasks)
		for _, parentTaskName := range parentTasks {
			parentTask := ctx.inv.Tasks[parentTaskName]
			applyInheritance(task, parentTask)
			envsToMerge = append(envsToMerge, parentTask.Environment)
		}
		envsToMerge = append(envsToMerge, ctx.taskModuleIndex[taskId].Environment)
		task.Environment, err = utils.ExpandMergeEnvs(envsToMerge...)
		if err != nil {
			return ctx.inv, fmt.Errorf("expanding task %v environment: %w", taskId, err)
		}
	}

	for taskId, task := range ctx.inv.Tasks {
		if task.Abstract {
			delete(ctx.inv.Tasks, taskId)
		}
	}

	for taskId, task := range ctx.inv.Tasks {
		for label, value := range task.Labels {
			task.Labels[label], err = utils.ExpandString(value, task.Environment)
			if err != nil {
				return ctx.inv, fmt.Errorf("expanding label %v in task %v: %w", label, taskId, err)
			}
		}
		if task.Abstract {
			delete(ctx.inv.Tasks, taskId)
		}
	}

	return ctx.inv, nil
}

func (ctx *InventoryContext) processModule(module *core.Module, moduleTrail []string, environment map[string]string) error {
	moduleEnvironment, err := utils.ExpandMergeEnvs(module.Environment, environment)
	if err != nil {
		return fmt.Errorf("expanding module environment: %w", err)
	}
	module.Environment = moduleEnvironment

	for taskName, task := range module.TasksSorted() {
		if task.WorkingDirectory == "" {
			task.WorkingDirectory = module.WorkingDirectory
		} else if !path.IsAbs(task.WorkingDirectory) {
			task.WorkingDirectory = path.Join(module.WorkingDirectory, task.WorkingDirectory)
		}

		taskId := core.MakeTaskId(moduleTrail, taskName)
		ctx.inv.Tasks[taskId] = task
		ctx.taskModuleIndex[taskId] = module
	}

	alreadyProcessedModules := make(map[string]bool)

	for importName, importObj := range module.ImportsSorted() {
		mergedEnv, err := utils.ExpandMergeEnvs(module.Modules[importName].Environment, importObj.Environment, module.Environment)
		if err != nil {
			return fmt.Errorf("expanding import %v environment: %w", importName, err)
		}
		alreadyProcessedModules[importName] = true
		module.Modules[importName].Environment = mergedEnv
	}

	for submoduleName, submodule := range module.ModulesSorted() {
		if !alreadyProcessedModules[submoduleName] {
			submoduleEnvironment, err := utils.ExpandMergeEnvs(submodule.Environment, module.Environment)
			submodule.Environment = submoduleEnvironment
			if err != nil {
				return fmt.Errorf("expanding module %v environment: %w", submoduleName, err)
			}
		}

		err = ctx.processModule(submodule, append(moduleTrail, submoduleName), submodule.Environment)
		if err != nil {
			return fmt.Errorf("processing module %v: %w", submoduleName, err)
		}
	}

	return nil
}

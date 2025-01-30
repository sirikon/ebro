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

func MakeInventory(indexedModule *core.IndexedModule, baseEnvironment *core.Environment) (Inventory, error) {
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
		parentTasks := slices.Clone(task.Extends)
		slices.Reverse(parentTasks)
		for _, parentTaskName := range parentTasks {
			parentTask := ctx.inv.Tasks[parentTaskName]
			applyInheritance(task, parentTask)
		}

		task.Environment, err = ctx.resolveTaskEnvironment(taskId)
		if err != nil {
			return ctx.inv, fmt.Errorf("resolving task %v environment: %w", taskId, err)
		}
		task.Extends = nil
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

func (ctx *InventoryContext) processModule(module *core.Module, moduleTrail []string, environment *core.Environment) error {
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

func (ctx *InventoryContext) resolveTaskEnvironment(taskId core.TaskId) (*core.Environment, error) {
	task := ctx.inv.Tasks[taskId]
	envsToMerge := []*core.Environment{
		task.Environment,
		core.NewEnvironment(
			core.EnvironmentValue{Key: "EBRO_TASK_ID", Value: string(taskId)},
			core.EnvironmentValue{Key: "EBRO_TASK_MODULE", Value: ":" + strings.Join(taskId.ModuleTrail(), ":")},
			core.EnvironmentValue{Key: "EBRO_TASK_NAME", Value: taskId.TaskName()},
			core.EnvironmentValue{Key: "EBRO_TASK_WORKING_DIRECTORY", Value: task.WorkingDirectory},
		),
	}
	parentTasks := slices.Clone(task.Extends)
	slices.Reverse(parentTasks)
	for _, parentTaskName := range parentTasks {
		parentTask := ctx.inv.Tasks[parentTaskName]
		envsToMerge = append(envsToMerge, parentTask.Environment)
	}
	envsToMerge = append(envsToMerge, ctx.taskModuleIndex[taskId].Environment)
	return utils.ExpandMergeEnvs(envsToMerge...)
}

package inventory

import (
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/utils"
)

type Inventory struct {
	Tasks map[string]*config.Task
}

type InventoryContext struct {
	inv             Inventory
	rootModule      *config.RootModule
	taskModuleIndex map[string]*config.Module
}

func MakeInventory(rootModule *config.RootModule) (Inventory, error) {
	ctx := InventoryContext{
		inv: Inventory{
			Tasks: make(map[string]*config.Task),
		},
		rootModule:      rootModule,
		taskModuleIndex: make(map[string]*config.Module),
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return ctx.inv, fmt.Errorf("obtaining working directory: %w", err)
	}

	err = ctx.processModule(rootModule.Module, []string{}, map[string]string{"EBRO_ROOT": workingDirectory})
	if err != nil {
		return ctx.inv, fmt.Errorf("processing module: %w", err)
	}

	inheritanceOrder, err := resolveInheritanceOrder(ctx.inv)
	if err != nil {
		return ctx.inv, fmt.Errorf("resolving inheritance order: %w", err)
	}

	for _, taskName := range inheritanceOrder {
		task := ctx.inv.Tasks[taskName]
		envsToMerge := [](map[string]string){
			task.Environment,
			map[string]string{"EBRO_TASK_WORKING_DIRECTORY": task.WorkingDirectory},
		}
		parentTasks := slices.Clone(task.Extends)
		slices.Reverse(parentTasks)
		for _, parentTaskName := range parentTasks {
			parentTask := ctx.inv.Tasks[parentTaskName]
			applyInheritance(task, parentTask)
			envsToMerge = append(envsToMerge, parentTask.Environment)
		}
		envsToMerge = append(envsToMerge, ctx.taskModuleIndex[taskName].Environment)
		task.Environment, err = utils.ExpandMergeEnvs(envsToMerge...)
		if err != nil {
			return ctx.inv, fmt.Errorf("expanding task %v environment: %w", taskName, err)
		}
	}

	for taskName, task := range ctx.inv.Tasks {
		if task.Abstract {
			delete(ctx.inv.Tasks, taskName)
		}
	}

	return ctx.inv, nil
}

func (ctx *InventoryContext) processModule(module *config.Module, moduleTrail []string, environment map[string]string) error {
	moduleEnvironment, err := utils.ExpandMergeEnvs(module.Environment, environment)
	if err != nil {
		return fmt.Errorf("expanding module environment: %w", err)
	}
	module.Environment = moduleEnvironment

	for taskName, task := range module.TasksSorted() {
		task.Requires, err = ctx.resolveRefs(task.Requires, moduleTrail)
		if err != nil {
			return err
		}
		task.RequiredBy, err = ctx.resolveRefs(task.RequiredBy, moduleTrail)
		if err != nil {
			return err
		}
		task.Extends, err = ctx.resolveRefs(task.Extends, moduleTrail)
		if err != nil {
			return err
		}

		if task.WorkingDirectory == "" {
			task.WorkingDirectory = module.WorkingDirectory
		} else if !path.IsAbs(task.WorkingDirectory) {
			task.WorkingDirectory = path.Join(module.WorkingDirectory, task.WorkingDirectory)
		}

		taskId := config.TaskId{ModuleTrail: moduleTrail, TaskName: taskName}
		ctx.inv.Tasks[taskId.String()] = task
		ctx.taskModuleIndex[taskId.String()] = module
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

func (ctx *InventoryContext) resolveRefs(s []string, moduleTrail []string) ([]string, error) {
	result := []string{}
	for _, taskReferenceString := range s {
		ref := config.MustParseTaskReference(taskReferenceString)
		taskId, _ := ctx.rootModule.GetTask(ref.Absolute(moduleTrail))
		if taskId != nil {
			result = append(result, taskId.String())
		} else if !ref.IsOptional {
			return nil, fmt.Errorf("referenced task %v does not exist", ref.String())
		}
	}
	return result, nil
}

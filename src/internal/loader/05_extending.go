package loader

import (
	"fmt"
	"slices"
	"strings"

	"github.com/sirikon/ebro/internal/core2"
	"github.com/sirikon/ebro/internal/dag"
	"github.com/sirikon/ebro/internal/utils"
	"github.com/sirikon/ebro/internal/utils2"

	"github.com/goccy/go-yaml"
)

func (ctx *loadCtx) extendingPhase() error {
	taskIds, err := resolveExtensionOrder(ctx.inventory)
	if err != nil {
		return fmt.Errorf("resolving extension order: %w", err)
	}

	for _, taskId := range taskIds {
		task := ctx.inventory.Task(taskId)
		parentTasks := slices.Clone(task.ExtendsIds)
		slices.Reverse(parentTasks)
		for _, parentTaskId := range parentTasks {
			parentTask := ctx.inventory.Task(parentTaskId)
			extendTask(task, parentTask)
		}

		task.Environment, err = resolveTaskEnvironment(ctx.inventory, ctx.baseEnvironment, taskId)
		if err != nil {
			return fmt.Errorf("resolving task %v environment: %w", taskId, err)
		}
		task.Extends = nil
	}

	for task := range ctx.inventory.Tasks() {
		if task.Abstract != nil && *task.Abstract {
			ctx.inventory.RemoveTask(task.Id)
		}
	}

	return nil
}

func resolveExtensionOrder(inventory *core2.Inventory) ([]core2.TaskId, error) {
	inheritanceDag := dag.NewDag[core2.TaskId]()
	target := []core2.TaskId{}

	for task := range inventory.Tasks() {
		target = append(target, task.Id)
		inheritanceDag.Link(task.Id, task.ExtendsIds...)
	}

	result, remains := inheritanceDag.Resolve(target)
	if remains != nil {
		remainsData, err := yaml.Marshal(remains)
		if err != nil {
			return nil, fmt.Errorf("inheritance order resolution could not complete. error while turning dependency index to yaml: %w", err)
		}
		return nil, fmt.Errorf("inheritance order resolution could not complete. "+
			"there could be a cyclic dependency. "+
			"here is the list of tasks and their inheritance data:\n%s", string(remainsData))
	}

	return result, nil
}

func extendTask(childTask *core2.Task, parentTask *core2.Task) {
	childTask.RequiresIds = utils.Dedupe(slices.Concat(childTask.RequiresIds, parentTask.RequiresIds))
	childTask.RequiredByIds = utils.Dedupe(slices.Concat(childTask.RequiredByIds, parentTask.RequiredByIds))

	if childTask.Script == "" {
		childTask.Script = parentTask.Script
	}
	if childTask.Quiet == nil {
		childTask.Quiet = parentTask.Quiet
	}
	if childTask.Interactive == nil {
		childTask.Interactive = parentTask.Interactive
	}
	if parentTask.When != nil {
		if childTask.When == nil {
			when := core2.When{}
			childTask.When = &when
		}
		if childTask.When.CheckFails == "" {
			childTask.When.CheckFails = parentTask.When.CheckFails
		}
		if childTask.When.OutputChanges == "" {
			childTask.When.OutputChanges = parentTask.When.OutputChanges
		}
	}

	if parentTask.Labels != nil {
		if childTask.Labels == nil {
			childTask.Labels = map[string]string{}
		}
		for key, val := range parentTask.Labels {
			if _, ok := childTask.Labels[key]; !ok {
				childTask.Labels[key] = val
			}
		}
	}
}

func resolveTaskEnvironment(inventory *core2.Inventory, baseEnvironment *core2.Environment, taskId core2.TaskId) (*core2.Environment, error) {
	task := inventory.Task(taskId)
	envsToMerge := []*core2.Environment{
		task.Environment,
		{
			Values: []core2.EnvironmentValue{
				{Key: "EBRO_TASK_ID", Value: string(taskId)},
				{Key: "EBRO_TASK_MODULE", Value: ":" + strings.Join(taskId.ModulePath(), ":")},
				{Key: "EBRO_TASK_NAME", Value: taskId.TaskName()},
				{Key: "EBRO_TASK_WORKING_DIRECTORY", Value: task.WorkingDirectory},
			},
		},
	}
	parentTasks := slices.Clone(task.ExtendsIds)
	slices.Reverse(parentTasks)
	for _, parentTaskName := range parentTasks {
		parentTask := inventory.Task(parentTaskName)
		envsToMerge = append(envsToMerge, parentTask.Environment)
	}
	for module := range inventory.WalkUpModulePath(task.Id) {
		envsToMerge = append(envsToMerge, module.Environment)
	}
	envsToMerge = append(envsToMerge, baseEnvironment)
	return utils2.ExpandMergeEnvs(envsToMerge...)
}

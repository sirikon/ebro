package loader

import (
	"fmt"
	"slices"
	"strings"

	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/utils"
)

func (ctx *loadCtx) environmentResolvingPhase(taskId core.TaskId) error {
	task := ctx.inventory.Task(taskId)
	var err error

	task.Environment, err = resolveTaskEnvironment(ctx.inventory, ctx.baseEnvironment, task.Id)
	if err != nil {
		return fmt.Errorf("resolving task '%v' environment: %w", task.Id, err)
	}

	return nil
}

func resolveTaskEnvironment(inventory *core.Inventory, baseEnvironment *core.Environment, taskId core.TaskId) (*core.Environment, error) {
	task := inventory.Task(taskId)
	envsToMerge := []*core.Environment{
		task.Environment,
		{
			Values: []core.EnvironmentValue{
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
	return utils.ExpandMergeEnvs(envsToMerge...)
}

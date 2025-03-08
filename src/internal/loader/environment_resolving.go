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

	task.Environment, err = resolveTaskEnvironment(ctx.inventory, ctx.baseEnvironment, task)
	if err != nil {
		return fmt.Errorf("resolving task '%v' environment: %w", task.Id, err)
	}

	return nil
}

func resolveTaskEnvironment(inventory *core.Inventory, baseEnvironment *core.Environment, task *core.Task) (*core.Environment, error) {
	// Environment resolution is done "bottom to top", which means that the
	// environment variables on the begginning of envsToMerge will be constructed
	// with the environment variables right after it.
	//
	// The first environment variables have priority over the rest.
	//
	// Example:
	//   A: "${B}!"
	//   B: "¡${C}"
	//   C: "Hello"
	//
	// Results in:
	//   A: "¡Hello!"
	//   B: "¡Hello"
	//   C: "Hello"
	//
	envsToMerge := []*core.Environment{}

	// Tasks' own defined environment variables
	envsToMerge = append(envsToMerge, task.Environment)

	// Built-in environment variables included for each Task
	envsToMerge = append(envsToMerge, &core.Environment{
		Values: []core.EnvironmentValue{
			{Key: "EBRO_MODULE", Value: ":" + strings.Join(task.Id.ModulePath(), ":")},
			{Key: "EBRO_TASK_ID", Value: string(task.Id)},
			{Key: "EBRO_TASK_NAME", Value: task.Id.TaskName()},
			{Key: "EBRO_TASK_WORKING_DIRECTORY", Value: task.WorkingDirectory},
		},
	})

	// Each task defined in `extends`
	//
	// NOTE: It's reversed because we're defining the environment
	// inside-to-outside.
	//
	// NOTE 2: We don't need to check these extended tasks' modules because
	// we're resolving the environment for each task in "extension order".
	// whenever we arrive here, all extended tasks have been already resolved
	// and we can just copy their final values.
	parentTasks := slices.Clone(task.ExtendsIds)
	slices.Reverse(parentTasks)
	for _, parentTaskName := range parentTasks {
		parentTask := inventory.Task(parentTaskName)
		envsToMerge = append(envsToMerge, parentTask.Environment)
	}

	// The environment variables of each module we're in, resolved
	// again inside-to-outside, so we start from the bottom module
	// up to the root.
	for modulePath, module := range inventory.WalkUpModulePath(task.Id.ModulePath()) {
		envsToMerge = append(envsToMerge, module.Environment)
		// Built-in environment variables included for each Module
		envsToMerge = append(envsToMerge, &core.Environment{
			Values: []core.EnvironmentValue{
				{Key: "EBRO_MODULE", Value: ":" + strings.Join(modulePath, ":")},
			},
		})
	}

	// Finally, the base environment defined at the beginning of
	// Ebro's execution.
	envsToMerge = append(envsToMerge, baseEnvironment)

	return utils.ExpandMergeEnvs(envsToMerge...)
}

func resolveModuleEnvironment(inventory *core.Inventory, baseEnvironment *core.Environment, module *core.Module) (*core.Environment, error) {
	envsToMerge := []*core.Environment{}

	for modulePath, module := range inventory.WalkUpModulePath(module.Path) {
		envsToMerge = append(envsToMerge, module.Environment)
		// Built-in environment variables included for each Module
		envsToMerge = append(envsToMerge, &core.Environment{
			Values: []core.EnvironmentValue{
				{Key: "EBRO_MODULE", Value: ":" + strings.Join(modulePath, ":")},
			},
		})
	}

	envsToMerge = append(envsToMerge, baseEnvironment)

	return utils.ExpandMergeEnvs(envsToMerge...)
}

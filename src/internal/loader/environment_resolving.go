package loader

import (
	"fmt"
	"slices"
	"strings"

	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/utils"
)

func (ctx *loadCtx) taskEnvironmentResolvingPhase(taskId core.TaskId) error {
	task := ctx.inventory.Task(taskId)
	var err error

	task.Environment, err = resolveTaskEnvironment(ctx.inventory, ctx.baseEnvironment, task)
	if err != nil {
		return fmt.Errorf("resolving task '%v' environment: %w", task.Id(), err)
	}

	return nil
}

func (ctx *loadCtx) moduleEnvironmentResolvingPhase(module *core.Module) error {
	var err error
	module.Environment, err = resolveModuleEnvironment(ctx.inventory, ctx.baseEnvironment, module)
	if err != nil {
		return fmt.Errorf("resolving module '%v' environment: %w", ":"+strings.Join(module.Path(), ":"), err)
	}
	return nil
}

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

func resolveTaskEnvironment(inventory *core.Inventory, baseEnvironment *core.Environment, task *core.Task) (*core.Environment, error) {
	envsToMerge := []*core.Environment{}

	// Tasks' own defined environment variables
	envsToMerge = append(envsToMerge, task.Environment)

	// Built-in environment variables included for each Task
	envsToMerge = append(envsToMerge, &core.Environment{
		Values: []core.EnvironmentValue{
			{Key: "EBRO_MODULE", Value: ":" + strings.Join(task.Module.Path(), ":")},
			{Key: "EBRO_TASK_ID", Value: string(task.Id())},
			{Key: "EBRO_TASK_NAME", Value: task.Name},
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

	// The environment variables of the module we're in.
	// It has been already resolved at this point (in a previous step)
	// so we can just copy the final values.
	module := inventory.Module(task.Module.Path())
	envsToMerge = append(envsToMerge, module.Environment)

	// Finally, the base environment defined at the beginning of
	// Ebro's execution.
	envsToMerge = append(envsToMerge, baseEnvironment)

	return utils.ExpandMergeEnvs(envsToMerge...)
}

func resolveModuleEnvironment(inventory *core.Inventory, baseEnvironment *core.Environment, module *core.Module) (*core.Environment, error) {
	envsToMerge := []*core.Environment{}

	// The environment variables of each module we're in, resolved
	// inside-to-outside, so we start from the bottom module up to the root.
	for modulePath, module := range inventory.WalkUpModulePath(module.Path()) {
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

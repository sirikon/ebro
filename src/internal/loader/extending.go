package loader

import (
	"fmt"
	"slices"

	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/dag"
	"github.com/sirikon/ebro/internal/utils"

	"github.com/goccy/go-yaml"
)

func (ctx *loadCtx) extendingPhase(taskId core.TaskId) error {
	task := ctx.inventory.Task(taskId)
	var err error

	parentTasks := slices.Clone(task.ExtendsIds)
	if len(parentTasks) > 0 {
		newTask := makeBaseTaskForExtending(task)
		for _, parentTaskId := range parentTasks {
			parentTask := ctx.inventory.Task(parentTaskId)
			extendTask(newTask, parentTask)
		}
		extendTask(newTask, task)
		task = newTask
		ctx.inventory.SetTask(newTask)
	}

	task.Environment, err = resolveTaskEnvironment(ctx.inventory, ctx.baseEnvironment, task.Id)
	if err != nil {
		return fmt.Errorf("resolving task '%v' environment: %w", task.Id, err)
	}

	return nil
}

func makeBaseTaskForExtending(task *core.Task) *core.Task {
	newTask := task.Clone()
	newTask.Requires = []string{}
	newTask.RequiredByExpressions = []string{}
	newTask.RequiresIds = []core.TaskId{}
	newTask.RequiredBy = []string{}
	newTask.RequiredByExpressions = []string{}
	newTask.RequiredByIds = []core.TaskId{}
	newTask.Script = ""
	newTask.Quiet = nil
	newTask.Interactive = nil
	newTask.When = nil
	newTask.Labels = nil
	return newTask
}

func extendTask(childTask *core.Task, parentTask *core.Task) {
	childTask.Requires = utils.Dedupe(slices.Concat(childTask.Requires, parentTask.Requires))
	childTask.RequiresExpressions = utils.Dedupe(slices.Concat(childTask.RequiresExpressions, parentTask.RequiresExpressions))
	childTask.RequiresIds = utils.Dedupe(slices.Concat(childTask.RequiresIds, parentTask.RequiresIds))
	childTask.RequiredBy = utils.Dedupe(slices.Concat(childTask.RequiredBy, parentTask.RequiredBy))
	childTask.RequiredByExpressions = utils.Dedupe(slices.Concat(childTask.RequiredByExpressions, parentTask.RequiredByExpressions))
	childTask.RequiredByIds = utils.Dedupe(slices.Concat(childTask.RequiredByIds, parentTask.RequiredByIds))

	if parentTask.Script != "" {
		childTask.Script = parentTask.Script
	}
	if parentTask.Quiet != nil {
		childTask.Quiet = parentTask.Quiet
	}
	if parentTask.Interactive != nil {
		childTask.Interactive = parentTask.Interactive
	}
	if parentTask.When != nil {
		if childTask.When == nil {
			when := core.When{}
			childTask.When = &when
		}
		if parentTask.When.CheckFails != "" {
			childTask.When.CheckFails = parentTask.When.CheckFails
		}
		if parentTask.When.OutputChanges != "" {
			childTask.When.OutputChanges = parentTask.When.OutputChanges
		}
	}

	if parentTask.Labels != nil {
		if childTask.Labels == nil {
			childTask.Labels = map[string]string{}
		}
		for key, val := range parentTask.Labels {
			childTask.Labels[key] = val
		}
	}
}

func (ctx *loadCtx) perTaskByExtensionOrder(taskPhases ...taskPhase) phase {
	return func() error {
		taskIds, err := resolveExtensionOrder(ctx.inventory)
		if err != nil {
			return fmt.Errorf("resolving extension order: %w", err)
		}
		for _, taskId := range taskIds {
			for _, taskPhase := range taskPhases {
				if err := taskPhase(taskId); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func resolveExtensionOrder(inventory *core.Inventory) ([]core.TaskId, error) {
	inheritanceDag := dag.NewDag[core.TaskId]()
	target := []core.TaskId{}

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

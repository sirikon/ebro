package loader

import (
	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/dag"
)

func (ctx *loadCtx) conditionalExistencePurgingPhase() error {
	purgeDag := dag.NewDag[core.TaskId]()

	targets := []core.TaskId{}
	for task := range ctx.inventory.Tasks() {
		if len(task.IfTasksExist) > 0 {
			for _, t := range task.IfTasksExist {
				parentTaskId := core.MustParseTaskReference(t).Absolute(task.Module.Path()).TaskId()
				purgeDag.Link(task.Id(), parentTaskId)
				targets = append(targets, task.Id())
			}
		}
	}

	taskIds, _ := purgeDag.Resolve(targets)

	for _, taskId := range taskIds {
		task := ctx.inventory.Task(taskId)
		if task == nil {
			continue
		}

		if len(task.IfTasksExist) > 0 {
			purge := false
			for _, t := range task.IfTasksExist {
				ref := core.MustParseTaskReference(t).Absolute(taskId.ModulePath())
				taskId, _ := ctx.inventory.FindTask(ref)
				if taskId == nil {
					purge = true
				}
			}
			if purge {
				ctx.inventory.RemoveTask(taskId)
			} else {
				task.IfTasksExist = []string{}
			}
		}
	}

	return nil
}

package loader

import (
	"github.com/sirikon/ebro/internal/core2"
	"github.com/sirikon/ebro/internal/dag"
)

func (ctx *loadCtx) purgingPhase() error {
	purgeDag := dag.NewDag[core2.TaskId]()

	targets := []core2.TaskId{}
	for task := range ctx.inventory.Tasks() {
		if len(task.IfTasksExist) > 0 {
			for _, t := range task.IfTasksExist {
				parentTaskId := core2.MustParseTaskReference(t).Absolute(task.Id.ModulePath()).TaskId()
				purgeDag.Link(task.Id, parentTaskId)
				targets = append(targets, task.Id)
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
				ref := core2.MustParseTaskReference(t).Absolute(taskId.ModulePath())
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

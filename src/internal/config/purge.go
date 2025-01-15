package config

import (
	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/dag"
)

func PurgeModule(indexedModule *IndexedModule) {
	purgeDag := dag.NewDag[core.TaskId]()

	targets := []core.TaskId{}
	for taskId, task := range indexedModule.AllTasks() {
		if len(task.IfTasksExist) > 0 {
			for _, t := range task.IfTasksExist {
				parentTaskId := mustParseTaskReference(t).Absolute(taskId.ModuleTrail()).TaskId()
				purgeDag.Link(taskId, parentTaskId)
				targets = append(targets, taskId)
			}
		}
	}

	taskIds, _ := purgeDag.Resolve(targets)

	for _, taskId := range taskIds {
		task := indexedModule.GetTask(taskId)
		if task == nil {
			continue
		}

		if len(task.IfTasksExist) > 0 {
			purge := false
			for _, t := range task.IfTasksExist {
				ref := mustParseTaskReference(t).Absolute(taskId.ModuleTrail())
				taskId, _ := FindTask(indexedModule, ref)
				if taskId == nil {
					purge = true
				}
			}
			if purge {
				indexedModule.RemoveTask(taskId)
			} else {
				task.IfTasksExist = []string{}
			}
		}
	}
}

package config

import (
	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/dag"
)

func PurgeModule(rootModule *RootModule) {
	purgeDag := dag.NewDag[core.TaskId]()

	targets := []core.TaskId{}
	for taskId, task := range rootModule.AllTasks() {
		if len(task.IfTasksExist) > 0 {
			for _, t := range task.IfTasksExist {
				parentTaskId := MustParseTaskReference(t).Absolute(taskId.ModuleTrail()).TaskId()
				purgeDag.Link(taskId, parentTaskId)
				targets = append(targets, taskId)
			}
		}
	}

	taskIds, _ := purgeDag.Resolve(targets)

	for _, taskId := range taskIds {
		task := rootModule.GetTask(taskId)
		if task == nil {
			continue
		}

		if len(task.IfTasksExist) > 0 {
			purge := false
			for _, t := range task.IfTasksExist {
				ref := MustParseTaskReference(t).Absolute(taskId.ModuleTrail())
				taskId, _ := FindTask(rootModule, ref)
				if taskId == nil {
					purge = true
				}
			}
			if purge {
				rootModule.RemoveTask(taskId)
			} else {
				task.IfTasksExist = []string{}
			}
		}
	}
}

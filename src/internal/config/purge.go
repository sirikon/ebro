package config

import (
	"github.com/sirikon/ebro/internal/dag"
)

func PurgeModule(rootModule *RootModule) {
	purgeDag := dag.NewDag()

	targets := []string{}
	for taskId, task := range rootModule.AllTasks() {
		if len(task.IfTasksExist) > 0 {
			for _, t := range task.IfTasksExist {
				parentTaskId := MustParseTaskReference(t).Absolute(taskId.ModuleTrail).TaskId()
				purgeDag.Link(taskId.String(), parentTaskId.String())
				targets = append(targets, taskId.String())
			}
		}
	}

	result, _ := purgeDag.Resolve(targets)

	for _, refStr := range result {
		taskId, task := rootModule.GetTask(MustParseTaskReference(refStr))
		if taskId == nil {
			continue
		}

		if len(task.IfTasksExist) > 0 {
			purge := false
			for _, t := range task.IfTasksExist {
				ref := MustParseTaskReference(t).Absolute(taskId.ModuleTrail)
				taskId, _ := rootModule.GetTask(ref)
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

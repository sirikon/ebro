package config

import "github.com/sirikon/ebro/internal/core"

type RootModule = core.RootModuleBase[Task, Import]

var NewRootModule = core.NewRootModuleBase[Task, Import, *RootModule]

func FindTask(rm *RootModule, taskReference TaskReference) (*core.TaskId, *Task) {
	if taskReference.IsRelative {
		panic("cannot call getTask with a relative taskReference")
	}

	taskId := taskReference.TaskId()
	if task, ok := rm.TaskIndex[taskId]; ok {
		return &taskId, task
	}

	taskId = taskReference.Concat("default").TaskId()
	if task, ok := rm.TaskIndex[taskId]; ok {
		return &taskId, task
	}

	return nil, nil
}

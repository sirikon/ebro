package config

import "github.com/sirikon/ebro/internal/core"

type IndexedModule = core.IndexedModuleBase[Task, Import]

var NewIndexedModule = core.NewIndexModuleBase[Task, Import, *IndexedModule]

func FindTask[TTask any, TImport any](rm *core.IndexedModuleBase[TTask, TImport], taskReference taskReference) (*core.TaskId, *TTask) {
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

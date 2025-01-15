package core

import (
	"iter"
	"maps"
	"slices"
)

type RootModule = RootModuleBase[Task, Import]

var NewRootModule = NewRootModuleBase[Task, Import, *RootModule]

type RootModuleBase[TTask any, TImport any] struct {
	Module          *ModuleBase[TTask, TImport]
	TaskIndex       map[TaskId]*TTask
	TaskModuleIndex map[TaskId]*ModuleBase[TTask, TImport]
}

func NewRootModuleBase[TTask any, TImport any, TRootModulePtr *RootModuleBase[TTask, TImport]](module *ModuleBase[TTask, TImport]) TRootModulePtr {
	rootModule := &RootModuleBase[TTask, TImport]{
		Module:          module,
		TaskIndex:       map[TaskId]*TTask{},
		TaskModuleIndex: map[TaskId]*ModuleBase[TTask, TImport]{},
	}
	rootModule.processModule(rootModule.Module, []string{})
	return rootModule
}

func (rm *RootModuleBase[TTask, TImport]) processModule(module *ModuleBase[TTask, TImport], moduleTrail []string) {
	for taskName, task := range module.TasksSorted() {
		taskId := MakeTaskId(moduleTrail, taskName)
		rm.TaskIndex[taskId] = task
		rm.TaskModuleIndex[taskId] = module
	}
	for moduleName, module := range module.ModulesSorted() {
		rm.processModule(module, append(moduleTrail, moduleName))
	}
}

func (rm *RootModuleBase[TTask, TImport]) GetTask(taskId TaskId) *TTask {
	if task, ok := rm.TaskIndex[taskId]; ok {
		return task
	}
	return nil
}

func (rm *RootModuleBase[TTask, TImport]) RemoveTask(taskId TaskId) {
	delete(rm.TaskModuleIndex[taskId].Tasks, taskId.TaskName())
	delete(rm.TaskModuleIndex, taskId)
	delete(rm.TaskIndex, taskId)
}

func (rm *RootModuleBase[TTask, TImport]) AllTasks() iter.Seq2[TaskId, *TTask] {
	taskIds := slices.Sorted(maps.Keys(rm.TaskIndex))
	return func(yield func(TaskId, *TTask) bool) {
		for _, taskId := range taskIds {
			if !yield(taskId, rm.TaskIndex[taskId]) {
				return
			}
		}
	}
}

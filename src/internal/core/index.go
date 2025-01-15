package core

import (
	"iter"
	"maps"
	"slices"
)

type IndexedModule = IndexedModuleBase[Task, Import]

var NewIndexedModule = NewIndexModuleBase[Task, Import, *IndexedModule]

type IndexedModuleBase[TTask any, TImport any] struct {
	Module          *ModuleBase[TTask, TImport]
	TaskIndex       map[TaskId]*TTask
	TaskModuleIndex map[TaskId]*ModuleBase[TTask, TImport]
}

func NewIndexModuleBase[TTask any, TImport any, TIndexedModulePtr *IndexedModuleBase[TTask, TImport]](module *ModuleBase[TTask, TImport]) TIndexedModulePtr {
	indexedModule := &IndexedModuleBase[TTask, TImport]{
		Module:          module,
		TaskIndex:       map[TaskId]*TTask{},
		TaskModuleIndex: map[TaskId]*ModuleBase[TTask, TImport]{},
	}
	indexedModule.processModule(indexedModule.Module, []string{})
	return indexedModule
}

func (rm *IndexedModuleBase[TTask, TImport]) processModule(module *ModuleBase[TTask, TImport], moduleTrail []string) {
	for taskName, task := range module.TasksSorted() {
		taskId := MakeTaskId(moduleTrail, taskName)
		rm.TaskIndex[taskId] = task
		rm.TaskModuleIndex[taskId] = module
	}
	for moduleName, module := range module.ModulesSorted() {
		rm.processModule(module, append(moduleTrail, moduleName))
	}
}

func (rm *IndexedModuleBase[TTask, TImport]) GetTask(taskId TaskId) *TTask {
	if task, ok := rm.TaskIndex[taskId]; ok {
		return task
	}
	return nil
}

func (rm *IndexedModuleBase[TTask, TImport]) RemoveTask(taskId TaskId) {
	delete(rm.TaskModuleIndex[taskId].Tasks, taskId.TaskName())
	delete(rm.TaskModuleIndex, taskId)
	delete(rm.TaskIndex, taskId)
}

func (rm *IndexedModuleBase[TTask, TImport]) AllTasks() iter.Seq2[TaskId, *TTask] {
	taskIds := slices.Sorted(maps.Keys(rm.TaskIndex))
	return func(yield func(TaskId, *TTask) bool) {
		for _, taskId := range taskIds {
			if !yield(taskId, rm.TaskIndex[taskId]) {
				return
			}
		}
	}
}

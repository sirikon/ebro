package config

import (
	"iter"
	"maps"
	"slices"
)

type RootModule struct {
	Module          *Module
	TaskIndex       map[TaskId]*Task
	TaskModuleIndex map[TaskId]*Module
}

func NewRootModule(module *Module) *RootModule {
	rootModule := &RootModule{
		Module:          module,
		TaskIndex:       map[TaskId]*Task{},
		TaskModuleIndex: map[TaskId]*Module{},
	}
	rootModule.processModule(rootModule.Module, []string{})
	return rootModule
}

func (rm *RootModule) processModule(module *Module, moduleTrail []string) {
	for taskName, task := range module.TasksSorted() {
		taskId := MakeTaskId(moduleTrail, taskName)
		rm.TaskIndex[taskId] = task
		rm.TaskModuleIndex[taskId] = module
	}
	for moduleName, module := range module.ModulesSorted() {
		rm.processModule(module, append(moduleTrail, moduleName))
	}
}

func (rm *RootModule) FindTask(taskReference TaskReference) (*TaskId, *Task) {
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

func (rm *RootModule) GetTask(taskId TaskId) *Task {
	if task, ok := rm.TaskIndex[taskId]; ok {
		return task
	}
	return nil
}

func (rm *RootModule) RemoveTask(taskId TaskId) {
	delete(rm.TaskModuleIndex[taskId].Tasks, taskId.TaskName())
	delete(rm.TaskModuleIndex, taskId)
	delete(rm.TaskIndex, taskId)
}

func (rm *RootModule) AllTasks() iter.Seq2[TaskId, *Task] {
	taskIds := slices.Sorted(maps.Keys(rm.TaskIndex))
	return func(yield func(TaskId, *Task) bool) {
		for _, taskId := range taskIds {
			if !yield(taskId, rm.TaskIndex[taskId]) {
				return
			}
		}
	}
}

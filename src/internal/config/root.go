package config

import (
	"iter"
	"maps"
	"slices"
)

type RootModule struct {
	Module          *Module
	TaskIdIndex     map[string]*TaskId
	TaskIndex       map[string]*Task
	TaskModuleIndex map[string]*Module
}

func NewRootModule(module *Module) *RootModule {
	rootModule := &RootModule{
		Module:          module,
		TaskIdIndex:     map[string]*TaskId{},
		TaskIndex:       map[string]*Task{},
		TaskModuleIndex: map[string]*Module{},
	}
	rootModule.processModule(rootModule.Module, []string{})
	return rootModule
}

func (rm *RootModule) processModule(module *Module, moduleTrail []string) {
	for taskName, task := range module.TasksSorted() {
		taskId := &TaskId{ModuleTrail: moduleTrail, TaskName: taskName}
		rm.TaskIdIndex[taskId.String()] = taskId
		rm.TaskIndex[taskId.String()] = task
		rm.TaskModuleIndex[taskId.String()] = module
	}
	for moduleName, module := range module.ModulesSorted() {
		rm.processModule(module, append(moduleTrail, moduleName))
	}
}

func (rm *RootModule) GetTask(taskReference TaskReference) (*TaskId, *Task) {
	if taskReference.IsRelative {
		panic("cannot call getTask with a relative taskReference")
	}
	if task, ok := rm.TaskIndex[taskReference.TaskId().String()]; ok {
		return taskReference.TaskId(), task
	}
	if task, ok := rm.TaskIndex[taskReference.Concat("default").TaskId().String()]; ok {
		return taskReference.Concat("default").TaskId(), task
	}
	return nil, nil
}

func (rm *RootModule) RemoveTask(taskId *TaskId) {
	delete(rm.TaskModuleIndex[taskId.String()].Tasks, taskId.TaskName)
	delete(rm.TaskModuleIndex, taskId.String())
	delete(rm.TaskIndex, taskId.String())
	delete(rm.TaskIdIndex, taskId.String())
}

func (rm *RootModule) AllTasks() iter.Seq2[*TaskId, *Task] {
	taskIdsStr := slices.Sorted(maps.Keys(rm.TaskIndex))
	return func(yield func(*TaskId, *Task) bool) {
		for _, taskIdStr := range taskIdsStr {
			if !yield(rm.TaskIdIndex[taskIdStr], rm.TaskIndex[taskIdStr]) {
				return
			}
		}
	}
}

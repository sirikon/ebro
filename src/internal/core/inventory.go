package core

import (
	"iter"
	"maps"
	"slices"
	"strings"
)

type Inventory struct {
	RootModule      *Module
	TaskIndex       map[TaskId]*Task
	TaskModuleIndex map[TaskId]*Module
	ModuleIndex     map[string]*Module
}

func NewInventory() *Inventory {
	return &Inventory{}
}

func (inv *Inventory) RefreshIndex() {
	inv.TaskIndex = map[TaskId]*Task{}
	inv.TaskModuleIndex = map[TaskId]*Module{}
	inv.ModuleIndex = map[string]*Module{}
	inv.generateIndex(inv.RootModule)
}

func (inv *Inventory) generateIndex(module *Module) {
	inv.ModuleIndex[strings.Join(module.Path, ":")] = module
	for _, task := range module.Tasks {
		inv.TaskIndex[task.Id] = task
		inv.TaskModuleIndex[task.Id] = module
	}
	for _, module := range module.Modules {
		inv.generateIndex(module)
	}
}

func (inv *Inventory) FindTask(taskReference TaskReference) (*TaskId, *Task) {
	if taskReference.IsRelative {
		panic("cannot call getTask with a relative taskReference")
	}

	taskId := taskReference.TaskId()
	if task, ok := inv.TaskIndex[taskId]; ok {
		return &taskId, task
	}

	taskId = taskReference.Concat("default").TaskId()
	if task, ok := inv.TaskIndex[taskId]; ok {
		return &taskId, task
	}

	return nil, nil
}

func (inv *Inventory) Task(taskId TaskId) *Task {
	if task, ok := inv.TaskIndex[taskId]; ok {
		return task
	}
	return nil
}

func (inv *Inventory) Tasks() iter.Seq[*Task] {
	taskIds := slices.Sorted(maps.Keys(inv.TaskIndex))
	return func(yield func(*Task) bool) {
		for _, taskId := range taskIds {
			if !yield(inv.TaskIndex[taskId]) {
				return
			}
		}
	}
}

func (inv *Inventory) WalkUpModulePath(taskId TaskId) iter.Seq[*Module] {
	modulePath := taskId.ModulePath()
	return func(yield func(*Module) bool) {
		for {
			if !yield(inv.ModuleIndex[strings.Join(modulePath, ":")]) {
				return
			}
			if len(modulePath) == 0 {
				return
			}
			modulePath = modulePath[:len(modulePath)-1]
		}
	}
}

func (inv *Inventory) RemoveTask(taskId TaskId) {
	delete(inv.TaskModuleIndex[taskId].Tasks, taskId.TaskName())
	delete(inv.TaskModuleIndex, taskId)
	delete(inv.TaskIndex, taskId)
}

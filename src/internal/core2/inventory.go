package core2

import (
	"iter"
	"maps"
	"slices"
)

type Inventory struct {
	RootModule  *Module
	TaskIndex   map[TaskId]*Task
	ModuleIndex map[TaskId]*Module
}

func NewInventory() *Inventory {
	return &Inventory{}
}

func (inv *Inventory) RefreshIndex() {
	inv.TaskIndex = map[TaskId]*Task{}
	inv.ModuleIndex = map[TaskId]*Module{}
	inv.generateIndex(inv.RootModule)
}

func (inv *Inventory) generateIndex(module *Module) {
	for _, task := range module.Tasks {
		inv.TaskIndex[task.Id] = task
		inv.ModuleIndex[task.Id] = module
	}
	for _, module := range module.Modules {
		inv.generateIndex(module)
	}
}

func (inv *Inventory) GetTask(taskId TaskId) *Task {
	if task, ok := inv.TaskIndex[taskId]; ok {
		return task
	}
	return nil
}

func (inv *Inventory) RemoveTask(taskId TaskId) {
	delete(inv.ModuleIndex[taskId].Tasks, taskId.TaskName())
	delete(inv.ModuleIndex, taskId)
	delete(inv.TaskIndex, taskId)
}

func (inv *Inventory) AllTasks() iter.Seq2[TaskId, *Task] {
	taskIds := slices.Sorted(maps.Keys(inv.TaskIndex))
	return func(yield func(TaskId, *Task) bool) {
		for _, taskId := range taskIds {
			if !yield(taskId, inv.TaskIndex[taskId]) {
				return
			}
		}
	}
}

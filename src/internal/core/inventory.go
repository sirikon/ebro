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
	inv.ModuleIndex[strings.Join(module.Path(), ":")] = module
	for _, task := range module.Tasks {
		inv.TaskIndex[task.Id()] = task
		inv.TaskModuleIndex[task.Id()] = module
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

func (inv *Inventory) Module(modulePath []string) *Module {
	if module, ok := inv.ModuleIndex[strings.Join(modulePath, ":")]; ok {
		return module
	}
	return nil
}

func (inv *Inventory) Task(taskId TaskId) *Task {
	if task, ok := inv.TaskIndex[taskId]; ok {
		return task
	}
	return nil
}

func (inv *Inventory) Modules() iter.Seq[*Module] {
	moduleIds := slices.Sorted(maps.Keys(inv.ModuleIndex))
	return func(yield func(*Module) bool) {
		for _, moduleId := range moduleIds {
			if !yield(inv.ModuleIndex[moduleId]) {
				return
			}
		}
	}
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

func (inv *Inventory) WalkDownModuleTree() iter.Seq[*Module] {
	return func(yield func(*Module) bool) {
		for module := range inv.moduleAndSubmodules(inv.RootModule) {
			if !yield(module) {
				return
			}
		}
	}
}

func (inv *Inventory) moduleAndSubmodules(module *Module) iter.Seq[*Module] {
	return func(yield func(*Module) bool) {
		if !yield(module) {
			return
		}

		visitedSubmodules := []string{}

		for {
			slices.Sort(visitedSubmodules)
			submoduleNames := slices.Sorted(maps.Keys(module.Modules))
			if slices.Equal(visitedSubmodules, submoduleNames) {
				return
			}

			for _, submoduleName := range submoduleNames {
				if !slices.Contains(visitedSubmodules, submoduleName) {
					visitedSubmodules = append(visitedSubmodules, submoduleName)
					for submodule := range inv.moduleAndSubmodules(module.Modules[submoduleName]) {
						if !yield(submodule) {
							return
						}
					}
				}
			}
		}
	}
}

func (inv *Inventory) WalkUpModulePath(modulePath []string) iter.Seq2[[]string, *Module] {
	return func(yield func([]string, *Module) bool) {
		for {
			if !yield(modulePath, inv.ModuleIndex[strings.Join(modulePath, ":")]) {
				return
			}
			if len(modulePath) == 0 {
				return
			}
			modulePath = modulePath[:len(modulePath)-1]
		}
	}
}

func (inv *Inventory) SetTask(task *Task) {
	inv.TaskModuleIndex[task.Id()].Tasks[task.Name] = task
	inv.TaskIndex[task.Id()] = task
}

func (inv *Inventory) RemoveTask(taskId TaskId) {
	delete(inv.TaskModuleIndex[taskId].Tasks, taskId.TaskName())
	delete(inv.TaskModuleIndex, taskId)
	delete(inv.TaskIndex, taskId)
}

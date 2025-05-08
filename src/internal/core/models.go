package core

import (
	"maps"
	"slices"
)

type Module struct {
	Name   string
	Parent *Module

	ForEach string

	Imports map[string]*Import
	Tasks   map[string]*Task
	Modules map[string]*Module

	WorkingDirectory string
	Environment      *Environment
	Labels           map[string]string
}

func (m *Module) Path() []string {
	path := []string{}
	currentModule := m
	for {
		if currentModule.Name == "" {
			break
		}

		path = append(path, currentModule.Name)
		if currentModule.Parent != nil {
			currentModule = currentModule.Parent
		} else {
			break
		}
	}
	slices.Reverse(path)
	return path
}

func (m *Module) Clone(newParent *Module) *Module {
	module := &Module{
		Parent:           newParent,
		ForEach:          m.ForEach,
		Imports:          cloneImportMap(m.Imports),
		WorkingDirectory: m.WorkingDirectory,
		Environment:      m.Environment.Clone(),
		Labels:           maps.Clone(m.Labels),
	}
	module.Tasks = cloneTaskMap(m.Tasks, module)
	module.Modules = cloneModuleMap(m.Modules, module)
	return module
}

func cloneModuleMap(m map[string]*Module, newParent *Module) map[string]*Module {
	result := make(map[string]*Module)
	for name := range maps.Keys(m) {
		result[name] = m[name].Clone(newParent)
	}
	return result
}

type Task struct {
	Name   string
	Module *Module

	IfTasksExist []string

	Requires              []string
	RequiresExpressions   []string
	RequiresScripts       []string
	RequiresIds           []TaskId
	RequiredBy            []string
	RequiredByExpressions []string
	RequiredByScripts     []string
	RequiredByIds         []TaskId

	Abstract   *bool
	Extends    []string
	ExtendsIds []TaskId

	Labels           map[string]string
	WorkingDirectory string
	Environment      *Environment

	Quiet       *bool
	Interactive *bool
	Script      []string
	When        *When
}

func (t *Task) Id() TaskId {
	return NewTaskId(t.Module.Path(), t.Name)
}

type When struct {
	CheckFails    []string
	OutputChanges []string
}

func (v *When) Clone() *When {
	if v == nil {
		return nil
	}
	return &When{
		CheckFails:    v.CheckFails,
		OutputChanges: v.OutputChanges,
	}
}

type Import struct {
	From        string
	Environment *Environment
	Module      *Module
}

func (i *Import) Clone() *Import {
	if i == nil {
		return nil
	}
	return &Import{
		From:        i.From,
		Environment: i.Environment.Clone(),
	}
}

func cloneImportMap(m map[string]*Import) map[string]*Import {
	result := make(map[string]*Import)
	for name := range maps.Keys(m) {
		result[name] = m[name].Clone()
	}
	return result
}

func (t *Task) Clone(newParent *Module) *Task {
	return &Task{
		Name:   t.Name,
		Module: newParent,

		IfTasksExist: slices.Clone(t.IfTasksExist),

		Requires:              slices.Clone(t.Requires),
		RequiresExpressions:   slices.Clone(t.RequiresExpressions),
		RequiresScripts:       slices.Clone(t.RequiresScripts),
		RequiresIds:           slices.Clone(t.RequiresIds),
		RequiredBy:            slices.Clone(t.RequiredBy),
		RequiredByExpressions: slices.Clone(t.RequiredByExpressions),
		RequiredByScripts:     slices.Clone(t.RequiredByScripts),
		RequiredByIds:         slices.Clone(t.RequiredByIds),

		Abstract:   cloneBoolPtr(t.Abstract),
		Extends:    slices.Clone(t.Extends),
		ExtendsIds: slices.Clone(t.ExtendsIds),

		Labels:           maps.Clone(t.Labels),
		WorkingDirectory: t.WorkingDirectory,
		Environment:      t.Environment.Clone(),

		Quiet:       cloneBoolPtr(t.Quiet),
		Interactive: cloneBoolPtr(t.Interactive),
		Script:      t.Script,
		When:        t.When.Clone(),
	}
}

func cloneTaskMap(m map[string]*Task, newParent *Module) map[string]*Task {
	result := make(map[string]*Task)
	for name := range maps.Keys(m) {
		result[name] = m[name].Clone(newParent)
	}
	return result
}

func cloneBoolPtr(v *bool) *bool {
	if v == nil {
		return nil
	}
	value := *v
	return &value
}

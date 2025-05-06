package core

import (
	"maps"
	"slices"
)

type Module struct {
	Path []string

	ForEach string

	Imports map[string]*Import
	Tasks   map[string]*Task
	Modules map[string]*Module

	WorkingDirectory string
	Environment      *Environment
	Labels           map[string]string
}

func (m *Module) Clone() *Module {
	return &Module{
		Path:             slices.Clone(m.Path),
		ForEach:          m.ForEach,
		Imports:          cloneImportMap(m.Imports),
		Tasks:            cloneTaskMap(m.Tasks),
		Modules:          cloneModuleMap(m.Modules),
		WorkingDirectory: m.WorkingDirectory,
		Environment:      m.Environment.Clone(),
		Labels:           maps.Clone(m.Labels),
	}
}

func cloneModuleMap(m map[string]*Module) map[string]*Module {
	result := make(map[string]*Module)
	for name := range maps.Keys(m) {
		result[name] = m[name].Clone()
	}
	return result
}

type Task struct {
	Name string
	Id   TaskId

	IfTasksExist []string

	Requires              []string
	RequiresExpressions   []string
	RequiresIds           []TaskId
	RequiredBy            []string
	RequiredByExpressions []string
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

func (t *Task) Clone() *Task {
	return &Task{
		Name: t.Name,
		Id:   t.Id,

		IfTasksExist: slices.Clone(t.IfTasksExist),

		Requires:              slices.Clone(t.Requires),
		RequiresExpressions:   slices.Clone(t.RequiresExpressions),
		RequiresIds:           slices.Clone(t.RequiresIds),
		RequiredBy:            slices.Clone(t.RequiredBy),
		RequiredByExpressions: slices.Clone(t.RequiredByExpressions),
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

func cloneTaskMap(m map[string]*Task) map[string]*Task {
	result := make(map[string]*Task)
	for name := range maps.Keys(m) {
		result[name] = m[name].Clone()
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

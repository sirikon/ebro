package core

import (
	"maps"
	"slices"
)

type Module struct {
	Path []string

	Imports map[string]*Import
	Tasks   map[string]*Task
	Modules map[string]*Module

	WorkingDirectory string
	Environment      *Environment
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
	CheckFails    string
	OutputChanges string
}

type Import struct {
	From        string
	Environment *Environment
	Module      *Module
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

func cloneBoolPtr(v *bool) *bool {
	if v == nil {
		return nil
	}
	value := *v
	return &value
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

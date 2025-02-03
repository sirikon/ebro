package core2

import "strings"

type Module struct {
	Path []string

	Imports map[string]*Import
	Tasks   map[string]*Task
	Modules map[string]*Module

	WorkingDirectory string
	Environment      *Environment
}

type TaskId string

func NewTaskId(modulePath []string, taskName string) TaskId {
	chunks := []string{""}
	chunks = append(chunks, modulePath...)
	chunks = append(chunks, taskName)
	result := TaskId(strings.Join(chunks, ":"))
	// result.MustBeValid()
	return result
}

type Task struct {
	Name string
	Id   TaskId

	IfTasksExist []string

	Requires      []string
	RequiresIds   []TaskId
	RequiredBy    []string
	RequiredByIds []TaskId

	Abstract   *bool
	Extends    []string
	ExtendsIds []TaskId

	Labels           map[string]string
	WorkingDirectory string
	Environment      *Environment

	Quiet       *bool
	Interactive *bool
	Script      string
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

type Environment struct {
	Values []EnvironmentValue
}

type EnvironmentValue struct {
	Key   string
	Value string
}

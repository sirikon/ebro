package core

import (
	"fmt"
	"iter"
	"maps"
	"regexp"
	"slices"
	"strings"
)

type TaskId string

var nameValidCharsRe = `[a-zA-Z0-9-_\.]`
var NameRe = regexp.MustCompile("^" + nameValidCharsRe + "+$")
var TaskIdRe = regexp.MustCompile("^(:" + nameValidCharsRe + "+)+$")

func ValidateName(name string) error {
	if !NameRe.MatchString(name) {
		return fmt.Errorf("name does not match the following regex: %v", NameRe.String())
	}
	return nil
}

func MakeTaskId(moduleTrail []string, taskName string) TaskId {
	chunks := []string{""}
	chunks = append(chunks, moduleTrail...)
	chunks = append(chunks, taskName)
	result := TaskId(strings.Join(chunks, ":"))
	result.MustBeValid()
	return result
}

func (tid TaskId) MustBeValid() {
	if !TaskIdRe.MatchString(string(tid)) {
		panic("TaskId ismalformed: " + string(tid))
	}
}

func (tid TaskId) ModuleTrail() []string {
	parts := tid.parts()
	return parts[:len(parts)-1]
}

func (tid TaskId) TaskName() string {
	parts := tid.parts()
	return parts[len(parts)-1]
}

func (tid TaskId) parts() []string {
	tid.MustBeValid()
	return strings.Split(strings.TrimPrefix(string(tid), ":"), ":")
}

type ModuleBase[TTask any, TImport any] struct {
	WorkingDirectory string                                 `yaml:"working_directory,omitempty"`
	Imports          map[string]*TImport                    `yaml:"imports,omitempty"`
	Environment      map[string]string                      `yaml:"environment,omitempty"`
	Tasks            map[string]*TTask                      `yaml:"tasks,omitempty"`
	Modules          map[string]*ModuleBase[TTask, TImport] `yaml:"modules,omitempty"`
}

type TaskBase[RefT ~string, WhenT any] struct {
	WorkingDirectory string            `yaml:"working_directory,omitempty"`
	IfTasksExist     []RefT            `yaml:"if_tasks_exist,omitempty"`
	Abstract         bool              `yaml:"abstract,omitempty"`
	Extends          []RefT            `yaml:"extends,omitempty"`
	Environment      map[string]string `yaml:"environment,omitempty"`
	Requires         []RefT            `yaml:"requires,omitempty"`
	RequiredBy       []RefT            `yaml:"required_by,omitempty"`
	Script           string            `yaml:"script,omitempty"`
	Quiet            *bool             `yaml:"quiet,omitempty"`
	When             *WhenT            `yaml:"when,omitempty"`
}

type When struct {
	CheckFails    string `yaml:"check_fails,omitempty"`
	OutputChanges string `yaml:"output_changes,omitempty"`
}

type Import struct {
	From        string            `yaml:"from,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

type Task = TaskBase[TaskId, When]
type Module = ModuleBase[Task, Import]

func (m *ModuleBase[TTask, TImport]) TasksSorted() iter.Seq2[string, *TTask] {
	taskNames := slices.Sorted(maps.Keys(m.Tasks))
	return func(yield func(string, *TTask) bool) {
		for _, taskName := range taskNames {
			if !yield(taskName, m.Tasks[taskName]) {
				return
			}
		}
	}
}

func (m *ModuleBase[TTask, TImport]) ModulesSorted() iter.Seq2[string, *ModuleBase[TTask, TImport]] {
	moduleNames := slices.Sorted(maps.Keys(m.Modules))
	return func(yield func(string, *ModuleBase[TTask, TImport]) bool) {
		for _, moduleName := range moduleNames {
			if !yield(moduleName, m.Modules[moduleName]) {
				return
			}
		}
	}
}

func (m *ModuleBase[TTask, TImport]) ImportsSorted() iter.Seq2[string, *TImport] {
	importNames := slices.Sorted(maps.Keys(m.Imports))
	return func(yield func(string, *TImport) bool) {
		for _, importName := range importNames {
			if !yield(importName, m.Imports[importName]) {
				return
			}
		}
	}
}

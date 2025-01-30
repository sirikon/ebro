package core

import (
	"fmt"
	"iter"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
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

type EnvironmentValue struct {
	Key   string
	Value string
}

type Environment struct {
	values []EnvironmentValue
}

type ModuleBase[TTask any, TImport any] struct {
	WorkingDirectory string                                 `yaml:"working_directory,omitempty"`
	Imports          map[string]*TImport                    `yaml:"imports,omitempty"`
	Environment      *Environment                           `yaml:"environment,omitempty"`
	Tasks            map[string]*TTask                      `yaml:"tasks,omitempty"`
	Modules          map[string]*ModuleBase[TTask, TImport] `yaml:"modules,omitempty"`
}

type TaskBase[RefT ~string, WhenT any] struct {
	Labels           map[string]string `yaml:"labels,omitempty"`
	WorkingDirectory string            `yaml:"working_directory,omitempty"`
	IfTasksExist     []RefT            `yaml:"if_tasks_exist,omitempty"`
	Abstract         bool              `yaml:"abstract,omitempty"`
	Extends          []RefT            `yaml:"extends,omitempty"`
	Environment      *Environment      `yaml:"environment,omitempty"`
	Requires         []RefT            `yaml:"requires,omitempty"`
	RequiredBy       []RefT            `yaml:"required_by,omitempty"`
	Script           string            `yaml:"script,omitempty"`
	Quiet            *bool             `yaml:"quiet,omitempty"`
	Interactive      *bool             `yaml:"interactive,omitempty"`
	When             *WhenT            `yaml:"when,omitempty"`
}

type When struct {
	CheckFails    string `yaml:"check_fails,omitempty"`
	OutputChanges string `yaml:"output_changes,omitempty"`
}

type Import struct {
	From        string       `yaml:"from,omitempty"`
	Environment *Environment `yaml:"environment,omitempty"`
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

func NewEnvironment(envValues ...EnvironmentValue) *Environment {
	return &Environment{
		values: envValues,
	}
}

func (env Environment) Values() iter.Seq[EnvironmentValue] {
	return func(yield func(EnvironmentValue) bool) {
		if env.values == nil {
			return
		}
		for _, envVal := range env.values {
			if !yield(envVal) {
				return
			}
		}
	}
}

func (env Environment) Map() map[string]string {
	result := map[string]string{}
	if env.values == nil {
		return result
	}

	for _, envVal := range env.values {
		result[envVal.Key] = envVal.Value
	}
	return result
}

func (env Environment) YamlMapSlice() yaml.MapSlice {
	result := yaml.MapSlice{}
	if env.values == nil {
		return result
	}
	for _, envVal := range env.values {
		result = append(result, yaml.MapItem{Key: envVal.Key, Value: envVal.Value})
	}
	return result
}

func (env *Environment) Set(key, value string) {
	if env.values == nil {
		env.values = []EnvironmentValue{}
	}
	existingPos := -1
	for i := range env.values {
		if env.values[i].Key == key {
			existingPos = i
		}
	}
	if existingPos >= 0 {
		env.values = append(env.values[:existingPos], env.values[existingPos+1:]...)
	}
	env.values = append(env.values, EnvironmentValue{
		Key:   key,
		Value: value,
	})
}

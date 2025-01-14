package config

import (
	"fmt"
	"iter"
	"maps"
	"regexp"
	"slices"
	"strings"
)

/* ==== Models reflected in Ebro.yaml files ==== */

type Module struct {
	WorkingDirectory string             `yaml:"working_directory,omitempty"`
	Imports          map[string]*Import `yaml:"imports,omitempty"`
	Environment      map[string]string  `yaml:"environment,omitempty"`
	Tasks            map[string]*Task   `yaml:"tasks,omitempty"`
	Modules          map[string]*Module `yaml:"modules,omitempty"`
}

type Task struct {
	WorkingDirectory string            `yaml:"working_directory,omitempty"`
	IfTasksExist     []string          `yaml:"if_tasks_exist,omitempty"`
	Abstract         bool              `yaml:"abstract,omitempty"`
	Extends          []string          `yaml:"extends,omitempty"`
	Environment      map[string]string `yaml:"environment,omitempty"`
	Requires         []string          `yaml:"requires,omitempty"`
	RequiredBy       []string          `yaml:"required_by,omitempty"`
	Script           string            `yaml:"script,omitempty"`
	Quiet            *bool             `yaml:"quiet,omitempty"`
	When             *When             `yaml:"when,omitempty"`
}

type When struct {
	CheckFails    string `yaml:"check_fails,omitempty"`
	OutputChanges string `yaml:"output_changes,omitempty"`
}

type Import struct {
	From        string            `yaml:"from,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

/* ============================================= */

func (m *Module) TasksSorted() iter.Seq2[string, *Task] {
	taskNames := slices.Sorted(maps.Keys(m.Tasks))
	return func(yield func(string, *Task) bool) {
		for _, taskName := range taskNames {
			if !yield(taskName, m.Tasks[taskName]) {
				return
			}
		}
	}
}

func (m *Module) ModulesSorted() iter.Seq2[string, *Module] {
	moduleNames := slices.Sorted(maps.Keys(m.Modules))
	return func(yield func(string, *Module) bool) {
		for _, moduleName := range moduleNames {
			if !yield(moduleName, m.Modules[moduleName]) {
				return
			}
		}
	}
}

func (m *Module) ImportsSorted() iter.Seq2[string, *Import] {
	importNames := slices.Sorted(maps.Keys(m.Imports))
	return func(yield func(string, *Import) bool) {
		for _, importName := range importNames {
			if !yield(importName, m.Imports[importName]) {
				return
			}
		}
	}
}

type TaskId string

func MakeTaskId(moduleTrail []string, taskName string) TaskId {
	chunks := []string{""}
	chunks = append(chunks, moduleTrail...)
	chunks = append(chunks, taskName)
	return TaskId(strings.Join(chunks, ":"))
}

func (tid TaskId) ModuleTrail() []string {
	parts := strings.Split(strings.TrimPrefix(string(tid), ":"), ":")
	return parts[:len(parts)-1]
}

func (tid TaskId) TaskName() string {
	parts := strings.Split(strings.TrimPrefix(string(tid), ":"), ":")
	return parts[len(parts)-1]
}

var taskReferenceRegex = regexp.MustCompile(`^:?[a-zA-Z0-9-_\.]+(:[a-zA-Z0-9-_\.]+)*\??$`)

type TaskReference struct {
	Path       []string
	IsRelative bool
	IsOptional bool
}

func ValidateTaskReference(text string) error {
	if !taskReferenceRegex.MatchString(text) {
		return fmt.Errorf("task reference is invalid")
	}
	return nil
}

func MustParseTaskReference(text string) TaskReference {
	result := TaskReference{
		Path:       []string{},
		IsRelative: true,
		IsOptional: false,
	}

	if err := ValidateTaskReference(text); err != nil {
		panic(err)
	}

	if strings.HasPrefix(text, ":") {
		text = strings.TrimPrefix(text, ":")
		result.IsRelative = false
	}

	if strings.HasSuffix(text, "?") {
		text = strings.TrimSuffix(text, "?")
		result.IsOptional = true
	}

	result.Path = strings.Split(text, ":")

	return result
}

func (tp TaskReference) Absolute(parentPath []string) TaskReference {
	if !tp.IsRelative {
		return tp
	}

	return TaskReference{
		Path:       append(parentPath, tp.Path...),
		IsRelative: false,
		IsOptional: tp.IsOptional,
	}
}

func (tp TaskReference) Concat(extraPath ...string) TaskReference {
	return TaskReference{
		Path:       append(tp.Path, extraPath...),
		IsRelative: tp.IsRelative,
		IsOptional: tp.IsOptional,
	}
}

func (tp TaskReference) TaskId() TaskId {
	if tp.IsRelative {
		panic("cannot build TaskId from relative TaskReference")
	}
	return MakeTaskId(tp.Path[:len(tp.Path)-1], tp.Path[len(tp.Path)-1])
}

func (tp TaskReference) PathString() string {
	chunks := []string{}
	if !tp.IsRelative {
		chunks = append(chunks, ":")
	}
	chunks = append(chunks, strings.Join(tp.Path, ":"))
	return strings.Join(chunks, "")
}

func (tp TaskReference) String() string {
	chunks := []string{}
	chunks = append(chunks, tp.PathString())
	if tp.IsOptional {
		chunks = append(chunks, "?")
	}
	return strings.Join(chunks, "")
}

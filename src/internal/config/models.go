package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirikon/ebro/internal/core"
)

type Module = core.ModuleBase[Task, Import]
type Task = core.TaskBase[string, When]
type Import = core.Import
type When = core.When

func MapToCoreModule(module *Module) *core.Module {
	return &core.Module{
		WorkingDirectory: module.WorkingDirectory,
		Environment:      module.Environment,
		Imports:          module.Imports,
		Modules:          mapToCoreModules(module.Modules),
		Tasks:            mapToCoreTasks(module.Tasks),
	}
}

func mapToCoreModules(modules map[string]*Module) map[string]*core.Module {
	result := map[string]*core.Module{}
	for moduleName, module := range modules {
		result[moduleName] = MapToCoreModule(module)
	}
	return result
}

func mapToCoreTasks(tasks map[string]*Task) map[string]*core.Task {
	result := map[string]*core.Task{}
	for taskName, task := range tasks {
		result[taskName] = mapToCoreTask(task)
	}
	return result
}

func mapToCoreTask(task *Task) *core.Task {
	return &core.Task{
		WorkingDirectory: task.WorkingDirectory,
		IfTasksExist:     mapToTaskIds(task.IfTasksExist),
		Abstract:         task.Abstract,
		Extends:          mapToTaskIds(task.Extends),
		Environment:      task.Environment,
		Requires:         mapToTaskIds(task.Requires),
		RequiredBy:       mapToTaskIds(task.RequiredBy),
		Script:           task.Script,
		Quiet:            task.Quiet,
		When:             task.When,
	}
}

func mapToTaskIds(taskIds []string) []core.TaskId {
	result := []core.TaskId{}
	for _, taskIdStr := range taskIds {
		taskId := core.TaskId(taskIdStr)
		taskId.MustBeValid()
		result = append(result, taskId)
	}
	return result
}

var taskReferenceRegex = regexp.MustCompile(`^:?[a-zA-Z0-9-_\.]+(:[a-zA-Z0-9-_\.]+)*\??$`)

type taskReference struct {
	Path       []string
	IsRelative bool
	IsOptional bool
}

func validateTaskReference(text string) error {
	if !taskReferenceRegex.MatchString(text) {
		return fmt.Errorf("task reference is invalid")
	}
	return nil
}

func mustParseTaskReference(text string) taskReference {
	result := taskReference{
		Path:       []string{},
		IsRelative: true,
		IsOptional: false,
	}

	if err := validateTaskReference(text); err != nil {
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

func (tp taskReference) Absolute(parentPath []string) taskReference {
	if !tp.IsRelative {
		return tp
	}

	return taskReference{
		Path:       append(parentPath, tp.Path...),
		IsRelative: false,
		IsOptional: tp.IsOptional,
	}
}

func (tp taskReference) Concat(extraPath ...string) taskReference {
	return taskReference{
		Path:       append(tp.Path, extraPath...),
		IsRelative: tp.IsRelative,
		IsOptional: tp.IsOptional,
	}
}

func (tp taskReference) TaskId() core.TaskId {
	if tp.IsRelative {
		panic("cannot build TaskId from relative TaskReference")
	}
	return core.MakeTaskId(tp.Path[:len(tp.Path)-1], tp.Path[len(tp.Path)-1])
}

func (tp taskReference) PathString() string {
	chunks := []string{}
	if !tp.IsRelative {
		chunks = append(chunks, ":")
	}
	chunks = append(chunks, strings.Join(tp.Path, ":"))
	return strings.Join(chunks, "")
}

func (tp taskReference) String() string {
	chunks := []string{}
	chunks = append(chunks, tp.PathString())
	if tp.IsOptional {
		chunks = append(chunks, "?")
	}
	return strings.Join(chunks, "")
}

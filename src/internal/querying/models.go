package querying

import (
	"strings"

	"github.com/sirikon/ebro/internal/core"
)

type QueryEnvironment struct {
	Tasks   []Task   `expr:"tasks"`
	Modules []Module `expr:"modules"`
}

type Module struct {
	Id               string            `expr:"id"`
	WorkingDirectory string            `expr:"working_directory"`
	Environment      map[string]string `expr:"environment"`
	Labels           map[string]string `expr:"labels"`
}

type Task struct {
	Id               string            `expr:"id"`
	Module           string            `expr:"module"`
	Name             string            `expr:"name"`
	WorkingDirectory string            `expr:"working_directory"`
	Environment      map[string]string `expr:"environment"`
	Labels           map[string]string `expr:"labels"`
	Extends          []string          `expr:"extends"`
	Requires         []string          `expr:"requires"`
	RequiredBy       []string          `expr:"required_by"`
	Script           []string          `expr:"script"`
	Quiet            *bool             `expr:"quiet"`
	Interactive      *bool             `expr:"interactive"`
	When             *When             `expr:"when"`
}

type When struct {
	CheckFails    []string `expr:"check_fails"`
	OutputChanges []string `expr:"output_changes"`
}

func buildQueryEnvironment(tasks []*core.Task, modules []*core.Module) QueryEnvironment {
	queryEnv := QueryEnvironment{
		Tasks:   []Task{},
		Modules: []Module{},
	}
	for _, task := range tasks {
		queryEnv.Tasks = append(queryEnv.Tasks, mapTask(task))
	}
	for _, module := range modules {
		queryEnv.Modules = append(queryEnv.Modules, mapModule(module))
	}
	return queryEnv
}

func mapTask(task *core.Task) Task {
	return Task{
		Id:               string(task.Id()),
		Module:           ":" + strings.Join(task.Module.Path(), ":"),
		Name:             task.Name,
		WorkingDirectory: task.WorkingDirectory,
		Environment:      task.Environment.Map(),
		Labels:           task.Labels,
		Extends:          mapTaskIds(task.ExtendsIds),
		Requires:         mapTaskIds(task.RequiresIds),
		RequiredBy:       mapTaskIds(task.RequiredByIds),
		Script:           task.Script,
		Quiet:            task.Quiet,
		Interactive:      task.Interactive,
		When:             mapWhen(task.When),
	}
}

func mapModule(module *core.Module) Module {
	return Module{
		Id:               ":" + strings.Join(module.Path(), ":"),
		WorkingDirectory: module.WorkingDirectory,
		Environment:      module.Environment.Map(),
		Labels:           module.Labels,
	}
}

func mapTaskIds(taskIds []core.TaskId) []string {
	result := []string{}
	for _, taskId := range taskIds {
		result = append(result, string(taskId))
	}
	return result
}

func mapWhen(when *core.When) *When {
	if when == nil {
		return nil
	}
	return &When{
		CheckFails:    when.CheckFails,
		OutputChanges: when.OutputChanges,
	}
}

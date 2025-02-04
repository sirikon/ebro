package querying

import (
	"strings"

	"github.com/sirikon/ebro/internal/core"
)

type QueryEnvironment struct {
	Tasks []Task `expr:"tasks"`
}

type Task struct {
	Id               string            `expr:"id"`
	Module           string            `expr:"module"`
	Name             string            `expr:"name"`
	Labels           map[string]string `expr:"labels"`
	WorkingDirectory string            `expr:"working_directory"`
	Extends          []string          `expr:"extends"`
	Environment      map[string]string `expr:"environment"`
	Requires         []string          `expr:"requires"`
	RequiredBy       []string          `expr:"required_by"`
	Script           string            `expr:"script"`
	Quiet            *bool             `expr:"quiet"`
	Interactive      *bool             `expr:"interactive"`
	When             *When             `expr:"when"`
}

type When struct {
	CheckFails    string `expr:"check_fails"`
	OutputChanges string `expr:"output_changes"`
}

func buildQueryEnvironment(tasks []*core.Task) QueryEnvironment {
	queryEnv := QueryEnvironment{
		Tasks: []Task{},
	}
	for _, task := range tasks {
		queryEnv.Tasks = append(queryEnv.Tasks, mapTask(task))
	}
	return queryEnv
}

func mapTask(task *core.Task) Task {
	return Task{
		Id:               string(task.Id),
		Module:           ":" + strings.Join(task.Id.ModulePath(), ":"),
		Name:             task.Name,
		Labels:           task.Labels,
		WorkingDirectory: task.WorkingDirectory,
		Extends:          mapTaskIds(task.ExtendsIds),
		Environment:      task.Environment.Map(),
		Requires:         mapTaskIds(task.RequiresIds),
		RequiredBy:       mapTaskIds(task.RequiredByIds),
		Script:           task.Script,
		Quiet:            task.Quiet,
		Interactive:      task.Interactive,
		When:             mapWhen(task.When),
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

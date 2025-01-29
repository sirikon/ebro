package filtering

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/sirikon/ebro/internal/cli"
	"github.com/sirikon/ebro/internal/core"
)

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
	When             *When             `expr:"when"`
}

type When struct {
	CheckFails    string `expr:"labels"`
	OutputChanges string `expr:"labels"`
}

func BuildTaskFilter(code string) (func(_ core.TaskId, _ *core.Task) bool, error) {
	program, err := expr.Compile(code, expr.Env(Task{}), expr.AsBool())
	if err != nil {
		return nil, fmt.Errorf("compiling filter expression: %w", err)
	}

	return func(taskId core.TaskId, task *core.Task) bool {
		taskInFilter := Task{
			Id:               string(taskId),
			Module:           ":" + strings.Join(taskId.ModuleTrail(), ":"),
			Name:             taskId.TaskName(),
			Labels:           task.Labels,
			WorkingDirectory: task.WorkingDirectory,
			Extends:          taskIdListToStringList(task.Extends),
			Environment:      task.Environment,
			Requires:         taskIdListToStringList(task.Requires),
			RequiredBy:       taskIdListToStringList(task.RequiredBy),
			Script:           task.Script,
			Quiet:            task.Quiet,
			When:             mapWhen(task.When),
		}
		output, err := expr.Run(program, taskInFilter)
		if err != nil {
			cli.ExitWithError(fmt.Errorf("running filter expression: %w", err))
		}
		if reflect.TypeOf(output).Kind() != reflect.Bool {
			cli.ExitWithError(fmt.Errorf("filter expression did not return a boolean when running with task %v. returned: %v", taskId, output))
		}
		return output.(bool)
	}, nil
}

func taskIdListToStringList(taskIds []core.TaskId) []string {
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

package querying

import (
	"fmt"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/sirikon/ebro/internal/cli"
	"github.com/sirikon/ebro/internal/core"
)

type QueryEnv struct {
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

func BuildQuery(code string) (func([]*core.Task) any, error) {
	program, err := expr.Compile(code, expr.Env(QueryEnv{}))
	if err != nil {
		return nil, fmt.Errorf("compiling query expression: %w", err)
	}

	return func(tasks []*core.Task) any {
		queryEnv := QueryEnv{
			Tasks: []Task{},
		}

		for _, task := range tasks {
			queryEnv.Tasks = append(queryEnv.Tasks, Task{
				Id:               string(task.Id),
				Module:           ":" + strings.Join(task.Id.ModulePath(), ":"),
				Name:             task.Name,
				Labels:           task.Labels,
				WorkingDirectory: task.WorkingDirectory,
				Extends:          taskIdListToStringList(task.ExtendsIds),
				Environment:      task.Environment.Map(),
				Requires:         taskIdListToStringList(task.RequiresIds),
				RequiredBy:       taskIdListToStringList(task.RequiredByIds),
				Script:           task.Script,
				Quiet:            task.Quiet,
				Interactive:      task.Interactive,
				When:             mapWhen(task.When),
			})
		}

		output, err := expr.Run(program, queryEnv)
		if err != nil {
			cli.ExitWithError(fmt.Errorf("running query expression: %w", err))
		}

		return output
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

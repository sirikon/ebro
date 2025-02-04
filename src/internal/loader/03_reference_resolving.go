package loader

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/sirikon/ebro/internal/core"
)

func (ctx *loadCtx) referenceResolvingPhase() error {
	var err error

	for task := range ctx.inventory.Tasks() {

		if result, err := resolveExpressions(ctx.inventory, task.RequiresExpressions); err != nil {
			return fmt.Errorf("resolving expressions in 'requires' for task '%v': %w", task.Id, err)
		} else {
			task.Requires = append(task.Requires, result...)
		}

		if result, err := resolveExpressions(ctx.inventory, task.RequiredByExpressions); err != nil {
			return fmt.Errorf("resolving expressions in 'required_by' for task '%v': %w", task.Id, err)
		} else {
			task.RequiredBy = append(task.RequiredBy, result...)
		}

		if task.RequiresIds, err = resolveReferences(ctx.inventory, task, task.Requires); err != nil {
			return fmt.Errorf("normalizing 'requires' for task '%v': %w", task.Id, err)
		}
		if task.RequiredByIds, err = resolveReferences(ctx.inventory, task, task.RequiredBy); err != nil {
			return fmt.Errorf("normalizing 'required_by' for task '%v': %w", task.Id, err)
		}
		if task.ExtendsIds, err = resolveReferences(ctx.inventory, task, task.Extends); err != nil {
			return fmt.Errorf("normalizing 'extends' for task '%v': %w", task.Id, err)
		}
	}

	return nil
}

func resolveExpressions(inventory *core.Inventory, expressions []string) ([]string, error) {
	result := []string{}
	for _, expression := range expressions {
		query, err := buildReferenceQuery(expression)
		if err != nil {
			return nil, err
		}

		queryResult, err := query(slices.Collect(inventory.Tasks()))
		if err != nil {
			return nil, err
		}
		ids := queryResult.([]interface{})
		for _, id := range ids {
			result = append(result, id.(string))
		}
	}
	return result, nil
}

func resolveReferences(inventory *core.Inventory, task *core.Task, taskReferences []string) ([]core.TaskId, error) {
	result := []core.TaskId{}
	for _, taskReference := range taskReferences {
		if err := core.ValidateTaskReference(taskReference); err != nil {
			return nil, fmt.Errorf("validating '%v': %w", taskReference, err)
		}

		ref := core.MustParseTaskReference(taskReference)
		if ref.IsRelative {
			ref = ref.Absolute(task.Id.ModulePath())
		}

		referencedTaskId, _ := inventory.FindTask(ref)
		if referencedTaskId == nil {
			if ref.IsOptional {
				continue
			} else {
				return nil, fmt.Errorf("referenced task '%v' does not exist", ref.TaskId())
			}
		}

		result = append(result, *referencedTaskId)
	}
	return result, nil
}

type ReferenceQueryEnvironment struct {
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
	Script           string            `expr:"script"`
	Quiet            *bool             `expr:"quiet"`
	Interactive      *bool             `expr:"interactive"`
	When             *When             `expr:"when"`
}

type When struct {
	CheckFails    string `expr:"check_fails"`
	OutputChanges string `expr:"output_changes"`
}

func buildQueryEnvironment(tasks []*core.Task) ReferenceQueryEnvironment {
	queryEnv := ReferenceQueryEnvironment{
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

func buildReferenceQuery(code string) (func([]*core.Task) (any, error), error) {
	program, err := expr.Compile(code, expr.Env(ReferenceQueryEnvironment{}), expr.AsKind(reflect.Slice))
	if err != nil {
		return nil, fmt.Errorf("compiling query expression: %w", err)
	}

	return func(tasks []*core.Task) (any, error) {
		queryEnv := buildQueryEnvironment(tasks)
		output, err := expr.Run(program, queryEnv)
		if err != nil {
			return nil, fmt.Errorf("running query expression: %w", err)
		}
		return output, nil
	}, nil
}

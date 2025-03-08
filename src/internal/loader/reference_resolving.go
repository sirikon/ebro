package loader

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/utils"
)

func (ctx *loadCtx) requirementExpressionReferenceResolvingPhase() error {
	for task := range ctx.inventory.Tasks() {
		if result, err := resolveExpressions(ctx.inventory, task.RequiresExpressions); err != nil {
			return fmt.Errorf("resolving expressions in 'requires' for task '%v': %w", task.Id, err)
		} else {
			task.RequiresIds = utils.Dedupe(append(task.RequiresIds, result...))
		}

		if result, err := resolveExpressions(ctx.inventory, task.RequiredByExpressions); err != nil {
			return fmt.Errorf("resolving expressions in 'required_by' for task '%v': %w", task.Id, err)
		} else {
			task.RequiredByIds = utils.Dedupe(append(task.RequiredByIds, result...))
		}
	}

	return nil
}

func (ctx *loadCtx) requirementReferenceResolvingPhase(taskId core.TaskId) error {
	task := ctx.inventory.Task(taskId)
	var err error

	if task.RequiresIds, err = core.ResolveReferences(ctx.inventory, task, task.Requires); err != nil {
		return fmt.Errorf("normalizing 'requires' for task '%v': %w", task.Id, err)
	}

	if task.RequiredByIds, err = core.ResolveReferences(ctx.inventory, task, task.RequiredBy); err != nil {
		return fmt.Errorf("normalizing 'required_by' for task '%v': %w", task.Id, err)
	}

	return nil
}

func (ctx *loadCtx) extensionReferenceResolvingPhase() error {
	var err error
	for task := range ctx.inventory.Tasks() {
		if task.ExtendsIds, err = core.ResolveReferences(ctx.inventory, task, task.Extends); err != nil {
			return fmt.Errorf("normalizing 'extends' for task '%v': %w", task.Id, err)
		}
	}
	return nil
}

func (ctx *loadCtx) abstractPurgingPhase() error {
	for task := range ctx.inventory.Tasks() {
		if task.Abstract != nil && *task.Abstract {
			ctx.inventory.RemoveTask(task.Id)
		}
	}
	return nil
}

func resolveExpressions(inventory *core.Inventory, expressions []string) ([]core.TaskId, error) {
	result := []core.TaskId{}
	for _, expression := range expressions {
		query, err := buildReferenceQuery(expression)
		if err != nil {
			return nil, err
		}

		queryResult, err := query(slices.Collect(inventory.Tasks()), slices.Collect(inventory.Modules()))
		if err != nil {
			return nil, err
		}
		ids := queryResult.([]interface{})
		for _, id := range ids {
			if reflect.TypeOf(id).Kind() != reflect.String {
				return nil, fmt.Errorf("wrong type returned from expression: %v", id)
			}
			refStr := id.(string)
			if err = core.ValidateTaskReference(refStr); err != nil {
				return nil, fmt.Errorf("checking %v: %w", refStr, err)
			}
			ref := core.MustParseTaskReference(refStr)
			if ref.IsRelative || ref.IsOptional {
				return nil, fmt.Errorf("checking %v: only task IDs are accepted, not relative or optional references", refStr)
			}

			result = append(result, ref.TaskId())
		}
	}
	return result, nil
}

type ReferenceQueryEnvironment struct {
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
	Script           []string          `expr:"script"`
	Quiet            *bool             `expr:"quiet"`
	Interactive      *bool             `expr:"interactive"`
	When             *When             `expr:"when"`
}

type When struct {
	CheckFails    []string `expr:"check_fails"`
	OutputChanges []string `expr:"output_changes"`
}

func buildQueryEnvironment(tasks []*core.Task, modules []*core.Module) ReferenceQueryEnvironment {
	queryEnv := ReferenceQueryEnvironment{
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
		Id:               string(task.Id),
		Module:           ":" + strings.Join(task.Id.ModulePath(), ":"),
		Name:             task.Name,
		Labels:           task.Labels,
		WorkingDirectory: task.WorkingDirectory,
		Environment:      task.Environment.Map(),
		Script:           task.Script,
		Quiet:            task.Quiet,
		Interactive:      task.Interactive,
		When:             mapWhen(task.When),
	}
}

func mapModule(module *core.Module) Module {
	return Module{
		Id:               ":" + strings.Join(module.Path, ":"),
		WorkingDirectory: module.WorkingDirectory,
		Environment:      module.Environment.Map(),
		Labels:           module.Labels,
	}
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

func buildReferenceQuery(code string) (func([]*core.Task, []*core.Module) (any, error), error) {
	program, err := expr.Compile(code, expr.Env(ReferenceQueryEnvironment{}), expr.AsKind(reflect.Slice))
	if err != nil {
		return nil, fmt.Errorf("compiling query expression: %w", err)
	}

	return func(tasks []*core.Task, modules []*core.Module) (any, error) {
		queryEnv := buildQueryEnvironment(tasks, modules)
		output, err := expr.Run(program, queryEnv)
		if err != nil {
			return nil, fmt.Errorf("running query expression: %w", err)
		}
		return output, nil
	}, nil
}

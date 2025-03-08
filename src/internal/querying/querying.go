package querying

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/sirikon/ebro/internal/core"
)

func BuildQuery(code string) (func([]*core.Task, []*core.Module) (any, error), error) {
	program, err := expr.Compile(code, expr.Env(QueryEnvironment{}))
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

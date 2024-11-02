package runner

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"github.com/sirikon/ebro/internal/cataloger"
	"github.com/sirikon/ebro/internal/planner"
)

func Run(catalog cataloger.Catalog, plan planner.Plan) error {
	for _, task_name := range plan {
		task := catalog[task_name]
		if task.SkipIf != "" {
			status, err := runScript(task.Script, *task.WorkingDirectory, task.Environment)
			if err != nil {
				return fmt.Errorf("running task %v skip_if: %w", task_name, err)
			}
			if status == 0 {
				color.Green("### skipping " + task_name)
				continue
			}
		}

		if task.Script != "" {
			color.Yellow("### running " + task_name)
			status, err := runScript(task.Script, *task.WorkingDirectory, task.Environment)
			if err != nil {
				return fmt.Errorf("running task %v script: %w", task_name, err)
			}
			if status != 0 {
				return fmt.Errorf("task %v returned status code %v", task_name, status)
			}
		}
	}
	return nil
}

func runScript(script string, working_directory string, environment map[string]string) (uint8, error) {
	file, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		return 1, fmt.Errorf("parsing script: %w", err)
	}

	runner, err := interp.New(
		interp.Env(expand.ListEnviron(append(os.Environ(), environmentToString(environment)...)...)),
		interp.Dir(working_directory),
		interp.StdIO(nil, os.Stdout, os.Stderr),
	)
	if err != nil {
		return 1, fmt.Errorf("runner creation failed: %w", err)
	}

	err = runner.Run(context.TODO(), file)
	if err == nil {
		return 0, nil
	}

	if status, ok := interp.IsExitStatus(err); ok {
		return status, nil
	} else {
		return 1, fmt.Errorf("runner returned error: %w", err)
	}
}

func environmentToString(environment map[string]string) []string {
	result := []string{}
	for key, value := range environment {
		result = append(result, key+"="+value)
	}
	return result
}

package runner

import (
	"context"
	"os"
	"strings"

	"github.com/sirikon/ebro/internal/cataloger"
	"github.com/sirikon/ebro/internal/planner"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Run(catalog cataloger.Catalog, plan planner.Plan) {
	for _, task_name := range plan {
		task := catalog[task_name]
		if task.Script != "" {
			runScript(task.Script, *task.WorkingDirectory, task.Environment)
		}
	}
}

func runScript(script string, working_directory string, environment map[string]string) {
	file, _ := syntax.NewParser().Parse(strings.NewReader(script), "")
	runner, _ := interp.New(
		interp.Env(expand.ListEnviron(append(os.Environ(), environmentToString(environment)...)...)),
		interp.Dir(working_directory),
		interp.StdIO(nil, os.Stdout, os.Stderr),
	)
	runner.Run(context.TODO(), file)
}

func environmentToString(environment map[string]string) []string {
	result := []string{}
	for key, value := range environment {
		result = append(result, key+"="+value)
	}
	return result
}

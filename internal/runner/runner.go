package runner

import (
	"context"
	"os"
	"strings"

	"github.com/sirikon/ebro/internal/indexer"
	"github.com/sirikon/ebro/internal/planner"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Run(index indexer.Index, plan planner.Plan) {
	for _, task_name := range plan {
		task := index[task_name]
		if task.Script != "" {
			runScript(task.Script)
		}
	}

}

func runScript(script string) {
	file, _ := syntax.NewParser().Parse(strings.NewReader(script), "")
	runner, _ := interp.New(
		interp.Env(nil),
		interp.StdIO(nil, os.Stdout, os.Stderr),
	)
	runner.Run(context.TODO(), file)
}

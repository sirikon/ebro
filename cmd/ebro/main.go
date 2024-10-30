package main

import (
	"fmt"
	"os"

	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/dag"
)

func main() {
	module, err := config.ParseFile("Ebro.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dagOutput := resolveDag(module)
	for _, task := range dagOutput.Tasks {
		fmt.Println("=== " + task)
	}
}

func resolveDag(module *config.Module) dag.Output {
	input := dag.Input{Tasks: make(map[string]dag.Task)}
	for name, task := range module.Tasks {
		input.Tasks[name] = dag.Task{
			Requires:   task.Requires,
			RequiredBy: task.RequiredBy,
		}
	}
	return dag.Resolve(input)
}

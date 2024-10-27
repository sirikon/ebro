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

	resolveDag(module, func(s string) error {
		fmt.Println("=== " + s)
		return nil
	})
}

func resolveDag(module *config.Module, cb func(string) error) {
	input := dag.Input{Tasks: make(map[string]dag.Task)}
	for name, task := range module.Tasks {
		input.Tasks[name] = dag.Task{
			Requires:   task.Requires,
			RequiredBy: task.RequiredBy,
		}
	}
	dag.Resolve(input, cb)
}

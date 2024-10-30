package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/sirikon/ebro/cmd/ebro/cli"
	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/dag"
)

func main() {
	module, err := config.ParseFile("Ebro.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	arguments := cli.Parse()

	if arguments.Flags.Config {
		bytes, err := yaml.Marshal(module)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(bytes))
		return
	}

	plan := makePlan(module)

	if arguments.Flags.Plan {
		for _, step := range plan.Steps {
			fmt.Println(step)
		}
		return
	}
}

func makePlan(module *config.Module) dag.Plan {
	input := dag.Input{Steps: make(map[string]dag.Step)}
	for name, step := range module.Tasks {
		input.Steps[name] = dag.Step{
			Requires:   step.Requires,
			RequiredBy: step.RequiredBy,
		}
	}
	return dag.Resolve(input)
}

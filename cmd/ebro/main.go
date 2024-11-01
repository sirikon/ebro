package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/sirikon/ebro/cmd/ebro/cli"
	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/indexer"
	"github.com/sirikon/ebro/internal/planner"
)

func main() {
	arguments := cli.Parse()

	config, err := config.DiscoverConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if arguments.Flags.Config {
		bytes, err := yaml.Marshal(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(string(bytes))
		return
	}

	index := indexer.Index(config)

	if arguments.Flags.Index {
		bytes, err := yaml.Marshal(index)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(string(bytes))
		return
	}

	plan := makePlan(index)

	if arguments.Flags.Plan {
		for _, step := range plan.Steps {
			fmt.Println(step)
		}
		return
	}
}

func makePlan(index map[string]config.Task) planner.Plan {
	input := planner.Input{Steps: make(map[string]planner.Step)}
	for name, step := range index {
		input.Steps[name] = planner.Step{
			Requires:   step.Requires,
			RequiredBy: step.RequiredBy,
		}
	}
	return planner.MakePlan(input)
}

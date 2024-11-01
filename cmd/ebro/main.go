package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/sirikon/ebro/cmd/ebro/cli"
	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/indexer"
	"github.com/sirikon/ebro/internal/planner"
	"github.com/sirikon/ebro/internal/runner"
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

	index := indexer.MakeIndex(config)

	if arguments.Flags.Index {
		bytes, err := yaml.Marshal(index)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(string(bytes))
		return
	}

	indexer.NormalizeTaskReferences(index, arguments.Targets)
	plan := planner.MakePlan(index, arguments.Targets)

	if arguments.Flags.Plan {
		for _, step := range plan {
			fmt.Println(step)
		}
		return
	}

	runner.Run(index, plan)
}

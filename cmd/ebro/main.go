package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/sirikon/ebro/cmd/ebro/cli"
	"github.com/sirikon/ebro/internal/cataloger"
	"github.com/sirikon/ebro/internal/config"
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

	catalog := cataloger.MakeCatalog(config)

	if arguments.Flags.Catalog {
		bytes, err := yaml.Marshal(catalog)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(string(bytes))
		return
	}

	cataloger.NormalizeTaskReferences(catalog, arguments.Targets)
	plan := planner.MakePlan(catalog, arguments.Targets)

	if arguments.Flags.Plan {
		for _, step := range plan {
			fmt.Println(step)
		}
		return
	}

	runner.Run(catalog, plan)
}

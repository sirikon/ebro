package main

import (
	"fmt"
	"os"
	"path"

	"github.com/gofrs/flock"
	"gopkg.in/yaml.v3"

	"github.com/sirikon/ebro/cmd/ebro/cli"
	"github.com/sirikon/ebro/internal/cataloger"
	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/planner"
	"github.com/sirikon/ebro/internal/runner"
)

func main() {
	err := lock()
	if err != nil {
		cli.ExitWithError(err)
	}

	arguments := cli.Parse()

	config, err := config.DiscoverConfig()
	if err != nil {
		cli.ExitWithError(err)
	}

	if arguments.Flags.Config {
		bytes, err := yaml.Marshal(config)
		if err != nil {
			cli.ExitWithError(err)
		}
		fmt.Print(string(bytes))
		return
	}

	catalog := cataloger.MakeCatalog(config)

	if arguments.Flags.Catalog {
		bytes, err := yaml.Marshal(catalog)
		if err != nil {
			cli.ExitWithError(err)
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

	err = runner.Run(catalog, plan)
	if err != nil {
		cli.ExitWithError(err)
	}
}

func lock() error {
	lockPath := path.Join(".ebro", "lock")
	err := os.MkdirAll(lockPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("obtaining lock for process: %w", err)
	}

	lock := flock.New(lockPath)
	err = lock.Lock()
	if err != nil {
		return fmt.Errorf("obtaining lock for process: %w", err)
	}
	return nil
}

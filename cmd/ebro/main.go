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

	config, err := config.ParseModuleFromFile(arguments.File)
	if err != nil {
		cli.ExitWithError(err)
	}

	if arguments.Command == cli.CommandConfig {
		bytes, err := yaml.Marshal(config)
		if err != nil {
			cli.ExitWithError(err)
		}
		fmt.Print(string(bytes))
		return
	}

	catalog, err := cataloger.MakeCatalog(config)
	if err != nil {
		cli.ExitWithError(err)
	}

	err = catalog.Validate()
	if err != nil {
		cli.ExitWithError(err)
	}

	if arguments.Command == cli.CommandCatalog {
		bytes, err := yaml.Marshal(catalog)
		if err != nil {
			cli.ExitWithError(err)
		}
		fmt.Print(string(bytes))
		return
	}

	cataloger.NormalizeTaskReferences(catalog, arguments.Targets)
	plan, err := planner.MakePlan(catalog, arguments.Targets)
	if err != nil {
		cli.ExitWithError(err)
	}

	if arguments.Command == cli.CommandPlan {
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
	err := os.MkdirAll(path.Dir(lockPath), os.ModePerm)
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

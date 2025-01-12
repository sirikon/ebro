package main

import (
	"fmt"
	"maps"
	"os"
	"path"
	"slices"

	"github.com/goccy/go-yaml"
	"github.com/gofrs/flock"

	"github.com/sirikon/ebro/internal/cli"
	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/inventory"
	"github.com/sirikon/ebro/internal/planner"
	"github.com/sirikon/ebro/internal/runner"
)

func main() {
	arguments := cli.Parse()

	if arguments.Command == cli.CommandHelp {
		cli.PrintHelp()
		os.Exit(0)
	}

	if arguments.Command == cli.CommandVersion {
		cli.PrintVersion()
		os.Exit(0)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		cli.ExitWithError(err)
	}

	modulePath := path.Join(workingDirectory, *arguments.GetFlagString(cli.FlagFile))
	rootModule, err := config.ParseModule(modulePath)
	if err != nil {
		cli.ExitWithError(err)
	}

	err = config.ValidateRootModule(rootModule)
	if err != nil {
		cli.ExitWithError(err)
	}

	inv, err := inventory.MakeInventory(arguments)
	if err != nil {
		cli.ExitWithError(err)
	}

	if arguments.Command == cli.CommandInventory {
		bytes, err := yaml.Marshal(inv.Tasks)
		if err != nil {
			cli.ExitWithError(err)
		}
		fmt.Print(string(bytes))
		return
	}

	if arguments.Command == cli.CommandList {
		taskNames := slices.Collect(maps.Keys(inv.Tasks))
		slices.Sort(taskNames)
		for _, taskName := range taskNames {
			fmt.Println(taskName)
		}
		return
	}

	inventory.NormalizeTaskNames(inv, arguments.Targets)
	plan, err := planner.MakePlan(inv, arguments.Targets)
	if err != nil {
		cli.ExitWithError(err)
	}

	if arguments.Command == cli.CommandPlan {
		for _, step := range plan {
			fmt.Println(step)
		}
		return
	}

	err = lock()
	if err != nil {
		cli.ExitWithError(err)
	}

	err = runner.Run(inv, plan, *arguments.GetFlagBool(cli.FlagForce))
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

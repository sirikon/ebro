package main

import (
	"fmt"
	"os"
	"path"

	"github.com/gofrs/flock"
	"gopkg.in/yaml.v3"

	"github.com/sirikon/ebro/internal/cli"
	"github.com/sirikon/ebro/internal/inventory"
	"github.com/sirikon/ebro/internal/planner"
	"github.com/sirikon/ebro/internal/runner"
)

func main() {
	err := lock()
	if err != nil {
		cli.ExitWithError(err)
	}

	arguments := cli.Parse()

	if arguments.Command == cli.CommandHelp {
		cli.PrintHelp()
		os.Exit(0)
	}

	if arguments.Command == cli.CommandVersion {
		cli.PrintVersion()
		os.Exit(0)
	}

	inv, err := inventory.MakeInventory(arguments)
	if err != nil {
		cli.ExitWithError(err)
	}

	if arguments.Command == cli.CommandInventory {
		bytes, err := yaml.Marshal(inv)
		if err != nil {
			cli.ExitWithError(err)
		}
		fmt.Print(string(bytes))
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

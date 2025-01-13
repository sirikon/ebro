package main

import (
	"fmt"
	"os"
	"path"

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

	module, err := config.ParseModule(rootModulePath(arguments))
	if err != nil {
		cli.ExitWithError(err)
	}

	err = config.ValidateRootModule(module)
	if err != nil {
		cli.ExitWithError(err)
	}
	rootModule := config.NewRootModule(module)
	config.PurgeModule(rootModule)

	inv, err := inventory.MakeInventory(rootModule)
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
		for taskName := range inv.TasksSorted() {
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

func rootModulePath(arguments cli.ExecutionArguments) string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		cli.ExitWithError(err)
	}
	filePath := *arguments.GetFlagString(cli.FlagFile)
	if !path.IsAbs(filePath) {
		filePath = path.Join(workingDirectory, filePath)
	}
	return filePath
}

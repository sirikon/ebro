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

	// -help
	if arguments.Command == cli.CommandHelp {
		cli.PrintHelp()
		os.Exit(0)
	}

	// -version
	if arguments.Command == cli.CommandVersion {
		cli.PrintVersion()
		os.Exit(0)
	}

	workingDirectory := getWorkingDirectory()

	indexedRootModule, err := config.ParseRootModule(workingDirectory, rootModulePath(workingDirectory, arguments))
	if err != nil {
		cli.ExitWithError(err)
	}

	inv, err := inventory.MakeInventory(indexedRootModule)
	if err != nil {
		cli.ExitWithError(err)
	}

	// -inventory
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

	targets, err := config.NormalizeTargets(indexedRootModule, arguments.Targets)
	if err != nil {
		cli.ExitWithError(err)
	}

	plan, err := planner.MakePlan(inv, targets)
	if err != nil {
		cli.ExitWithError(err)
	}

	// -plan
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

func getWorkingDirectory() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		cli.ExitWithError(fmt.Errorf("obtaining working directory: %w", err))
	}
	return workingDirectory
}

func rootModulePath(workingDirectory string, arguments cli.ExecutionArguments) string {
	filePath := *arguments.GetFlagString(cli.FlagFile)
	if !path.IsAbs(filePath) {
		filePath = path.Join(workingDirectory, filePath)
	}
	return filePath
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

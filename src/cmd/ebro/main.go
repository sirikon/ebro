package main

import (
	"fmt"
	"os"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/gofrs/flock"

	"github.com/sirikon/ebro/internal/cli"
	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/filtering"
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

	indexedRootModule, err := config.ParseRootModule(rootModulePath(workingDirectory, arguments))
	if err != nil {
		cli.ExitWithError(err)
	}

	baseEnvironment := core.NewEnvironment(
		core.EnvironmentValue{Key: "EBRO_BIN", Value: arguments.Bin},
		core.EnvironmentValue{Key: "EBRO_ROOT", Value: workingDirectory},
	)

	inv, err := inventory.MakeInventory(indexedRootModule, baseEnvironment)
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

	// -list
	if arguments.Command == cli.CommandList {
		taskFilter := buildTaskFilter(arguments)
		for taskId := range inv.TasksSorted() {
			if taskFilter(taskId, inv.Tasks[taskId]) {
				fmt.Println(taskId)
			}
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

func buildTaskFilter(arguments cli.ExecutionArguments) func(_ core.TaskId, _ *core.Task) bool {
	filterExpression := *arguments.GetFlagString(cli.FlagFilter)
	if filterExpression != "" {
		filter, err := filtering.BuildTaskFilter(filterExpression)
		if err != nil {
			cli.ExitWithError(err)
		}
		return filter
	}

	return func(_ core.TaskId, _ *core.Task) bool {
		return true
	}
}

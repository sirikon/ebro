package main

import (
	"fmt"
	"os"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/gofrs/flock"

	"github.com/sirikon/ebro/internal/cli"
	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/core2"
	"github.com/sirikon/ebro/internal/loader"
	"github.com/sirikon/ebro/internal/querying"
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
	rootFile := rootFilePath(workingDirectory, arguments)

	baseEnvironment := &core2.Environment{
		Values: []core2.EnvironmentValue{
			{Key: "EBRO_BIN", Value: arguments.Bin},
			{Key: "EBRO_ROOT", Value: workingDirectory},
			{Key: "EBRO_ROOT_FILE", Value: rootFile},
		},
	}

	inventory, err := loader.Load(baseEnvironment, workingDirectory, rootFile)
	if err != nil {
		cli.ExitWithError(err)
	}

	// -inventory
	if arguments.Command == cli.CommandInventory {
		// inventoryQuery := buildInventoryQuery(arguments)
		// if inventoryQuery != nil {
		// 	result = inventoryQuery(inv.Tasks)
		// }

		// if reflect.TypeOf(result).Kind() == reflect.String {
		// 	fmt.Println(result)
		// 	return
		// }

		inventoryView := InventoryView{}
		for task := range inventory.Tasks() {
			inventoryView[task.Id] = TaskView{
				Labels:           task.Labels,
				WorkingDirectory: task.WorkingDirectory,
				Environment:      task.Environment.YamlMapSlice(),
				Requires:         taskIdsToView(task.RequiresIds),
				RequiredBy:       taskIdsToView(task.RequiredByIds),
				Script:           task.Script,
				Interactive:      task.Interactive,
				Quiet:            task.Quiet,
				When:             whenToView(task.When),
			}
		}

		bytes, err := yaml.Marshal(inventoryView)
		if err != nil {
			cli.ExitWithError(err)
		}
		fmt.Print(string(bytes))
		return
	}

	// -list
	if arguments.Command == cli.CommandList {
		for task := range inventory.Tasks() {
			fmt.Println(task.Id)
		}
		return
	}

	// targets, err := config.NormalizeTargets(indexedRootModule, arguments.Targets)
	// if err != nil {
	// 	cli.ExitWithError(err)
	// }

	// plan, err := planner.MakePlan(inv, targets)
	// if err != nil {
	// 	cli.ExitWithError(err)
	// }

	// // -plan
	// if arguments.Command == cli.CommandPlan {
	// 	for _, step := range plan {
	// 		fmt.Println(step)
	// 	}
	// 	return
	// }

	// err = lock()
	// if err != nil {
	// 	cli.ExitWithError(err)
	// }

	// err = runner.Run(inv, plan, *arguments.GetFlagBool(cli.FlagForce))
	// if err != nil {
	// 	cli.ExitWithError(err)
	// }
}

func getWorkingDirectory() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		cli.ExitWithError(fmt.Errorf("obtaining working directory: %w", err))
	}
	return workingDirectory
}

func rootFilePath(workingDirectory string, arguments cli.ExecutionArguments) string {
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

func buildInventoryQuery(arguments cli.ExecutionArguments) func(map[core.TaskId]*core.Task) any {
	queryExpression := *arguments.GetFlagString(cli.FlagQuery)
	if queryExpression != "" {
		query, err := querying.BuildQuery(queryExpression)
		if err != nil {
			cli.ExitWithError(err)
		}
		return query
	}
	return nil
}

type InventoryView map[core2.TaskId]TaskView

type TaskView struct {
	Labels           map[string]string `yaml:"labels,omitempty"`
	WorkingDirectory string            `yaml:"working_directory,omitempty"`
	Environment      yaml.MapSlice     `yaml:"environment,omitempty"`
	Requires         []string          `yaml:"requires,omitempty"`
	RequiredBy       []string          `yaml:"required_by,omitempty"`
	Script           string            `yaml:"script,omitempty"`
	Interactive      *bool             `yaml:"interactive,omitempty"`
	Quiet            *bool             `yaml:"quiet,omitempty"`
	When             *WhenView         `yaml:"when,omitempty"`
}

type WhenView struct {
	CheckFails    string `yaml:"check_fails,omitempty"`
	OutputChanges string `yaml:"output_changes,omitempty"`
}

func taskIdsToView(taskIds []core2.TaskId) []string {
	if taskIds == nil {
		return nil
	}
	result := []string{}
	for _, taskId := range taskIds {
		result = append(result, string(taskId))
	}
	return result
}

func whenToView(when *core2.When) *WhenView {
	if when == nil {
		return nil
	}
	return &WhenView{
		CheckFails:    when.CheckFails,
		OutputChanges: when.OutputChanges,
	}
}

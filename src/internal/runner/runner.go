package runner

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"github.com/sirikon/ebro/internal/inventory"
	"github.com/sirikon/ebro/internal/logger"
	"github.com/sirikon/ebro/internal/planner"
)

func Run(inv inventory.Inventory, plan planner.Plan, force bool) error {
	for _, taskName := range plan {
		task := inv[taskName]

		if task.Script == "" {
			logger.Info(logLine(taskName, "satisfied"))
			continue
		}

		skip := false
		if task.When != nil && !force {
			skip = true

			if task.When.CheckFails != "" {
				output := bytes.Buffer{}
				outputWriter := bufio.NewWriter(&output)
				status, err := runScriptWithIO(task.When.CheckFails, task.WorkingDirectory, task.Environment, outputWriter, outputWriter)
				if err != nil {
					return fmt.Errorf("running task %v when.check_fails: %w", taskName, err)
				}
				outputWriter.Flush()
				if status > 0 {
					fmt.Print(output.String())
					skip = false
				}
			}

			if task.When.OutputChanges != "" {
				output := bytes.Buffer{}
				outputWriter := bufio.NewWriter(&output)
				status, err := runScriptWithIO(task.When.OutputChanges, task.WorkingDirectory, task.Environment, outputWriter, outputWriter)
				if err != nil {
					return fmt.Errorf("running task %v when.output_changes: %w", taskName, err)
				}
				outputWriter.Flush()
				if status > 0 {
					return fmt.Errorf("task %v when.output_changes returned status code %v. here is the output:\n%v", taskName, status, output.String())
				}

				outputChanged, err := storeTaskOutputAndCheckIfChanged(taskName, output.Bytes())
				if err != nil {
					return fmt.Errorf("storing output for task %v when.output_changes: %w", taskName, err)
				}
				if outputChanged {
					skip = false
				}
			}
		}

		if skip {
			logger.Info(logLine(taskName, "skipping"))
			continue
		}

		logger.Notice(logLine(taskName, "running"))
		status, err := runScript(task.Script, task.WorkingDirectory, task.Environment)

		var final_err error

		if err != nil {
			final_err = fmt.Errorf("running task %v script: %w", taskName, err)
		}

		if status != 0 {
			final_err = fmt.Errorf("task %v returned status code %v", taskName, status)
		}

		if final_err != nil {
			err := removeTaskOutput(taskName)
			if err != nil {
				return fmt.Errorf("removing output after failure of task %v: %w", taskName, err)
			}
			return final_err
		}
	}
	return nil
}

func logLine(taskName string, message string) string {
	return "[" + taskName + "] " + message
}

func runScript(script string, workingDirectory string, environment map[string]string) (uint8, error) {
	return runScriptWithIO(script, workingDirectory, environment, os.Stdout, os.Stdout)
}

func runScriptWithIO(script string, workingDirectory string, environment map[string]string, stdout io.Writer, stderr io.Writer) (uint8, error) {
	script_header := []string{"set -euo pipefail"}

	file, err := syntax.NewParser().Parse(strings.NewReader(strings.Join(script_header, "\n")+"\n"+script), "")
	if err != nil {
		return 1, fmt.Errorf("parsing script: %w", err)
	}

	runner, err := interp.New(
		interp.Env(expand.ListEnviron(append(os.Environ(), environmentToString(environment)...)...)),
		interp.Dir(workingDirectory),
		interp.StdIO(nil, stdout, stderr),
	)
	if err != nil {
		return 1, fmt.Errorf("runner creation failed: %w", err)
	}

	err = runner.Run(context.TODO(), file)
	if err == nil {
		return 0, nil
	}

	if status, ok := interp.IsExitStatus(err); ok {
		return status, nil
	} else {
		return 1, fmt.Errorf("runner returned error: %w", err)
	}
}

func environmentToString(environment map[string]string) []string {
	result := []string{}
	for key, value := range environment {
		result = append(result, key+"="+value)
	}
	return result
}

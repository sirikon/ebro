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

	"github.com/sirikon/ebro/internal/cataloger"
	"github.com/sirikon/ebro/internal/logger"
	"github.com/sirikon/ebro/internal/planner"
)

func Run(catalog cataloger.Catalog, plan planner.Plan, force bool) error {
	for _, task_name := range plan {
		task := catalog[task_name]

		if task.Script == "" {
			logger.Info(logLine(task_name, "satisfied"))
			continue
		}

		skip := false
		if task.When != nil && !force {
			skip = true

			if task.When.CheckFails != "" {
				output := bytes.Buffer{}
				outputWriter := bufio.NewWriter(&output)
				status, err := runScriptWithIo(task.When.CheckFails, *task.WorkingDirectory, task.Environment, outputWriter, outputWriter)
				if err != nil {
					return fmt.Errorf("running task %v when.check_fails: %w", task_name, err)
				}
				outputWriter.Flush()
				if status > 0 {
					fmt.Println(output.String())
					skip = false
				}
			}

			if task.When.OutputChanges != "" {
				output := bytes.Buffer{}
				outputWriter := bufio.NewWriter(&output)
				status, err := runScriptWithIo(task.When.OutputChanges, *task.WorkingDirectory, task.Environment, outputWriter, outputWriter)
				if err != nil {
					return fmt.Errorf("running task %v when.output_changes: %w", task_name, err)
				}
				outputWriter.Flush()
				if status > 0 {
					return fmt.Errorf("task %v when.output_changes returned status code %v. here is the output:\n%v", task_name, status, output.String())
				}

				outputChanged, err := storeTaskOutputAndCheckIfChanged(task_name, output.Bytes())
				if err != nil {
					return fmt.Errorf("storing output for task %v when.output_changes: %w", task_name, err)
				}
				if outputChanged {
					skip = false
				}
			}
		}

		if skip {
			logger.Info(logLine(task_name, "skipping"))
			continue
		}

		logger.Notice(logLine(task_name, "running"))
		status, err := runScript(task.Script, *task.WorkingDirectory, task.Environment)
		if err != nil {
			return fmt.Errorf("running task %v script: %w", task_name, err)
		}
		if status != 0 {
			return fmt.Errorf("task %v returned status code %v", task_name, status)
		}
	}
	return nil
}

func logLine(task_name string, message string) string {
	return "[" + task_name + "] " + message
}

func runScript(script string, working_directory string, environment map[string]string) (uint8, error) {
	return runScriptWithIo(script, working_directory, environment, os.Stdout, os.Stderr)
}

func runScriptWithIo(script string, working_directory string, environment map[string]string, stdout io.Writer, stderr io.Writer) (uint8, error) {
	file, err := syntax.NewParser().Parse(strings.NewReader("set -euo pipefail\n"+script), "")
	if err != nil {
		return 1, fmt.Errorf("parsing script: %w", err)
	}

	runner, err := interp.New(
		interp.Env(expand.ListEnviron(append(os.Environ(), environmentToString(environment)...)...)),
		interp.Dir(working_directory),
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

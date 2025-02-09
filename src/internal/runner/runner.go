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

	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/logger"
	"github.com/sirikon/ebro/internal/planner"
)

func Run(inv *core.Inventory, plan planner.Plan, force bool) error {
	for _, taskId := range plan {
		task := inv.Task(taskId)
		if task == nil {
			return fmt.Errorf("task %v does not exist", taskId)
		}

		if len(task.Script) == 0 {
			logger.Info(logLine(taskId, "satisfied"))
			continue
		}

		skip := false
		if task.When != nil && !force {
			skip = true

			if task.When.CheckFails != "" {
				output := bytes.Buffer{}
				outputWriter := bufio.NewWriter(&output)
				status, err := runScript(task.When.CheckFails, task.WorkingDirectory, task.Environment, nil, outputWriter, outputWriter)
				if err != nil {
					return fmt.Errorf("running task %v when.check_fails: %w", taskId, err)
				}
				err = outputWriter.Flush()
				if err != nil {
					return fmt.Errorf("running task %v when.check_fails: flushing writer: %w", taskId, err)
				}
				if status > 0 {
					skip = false
				}
			}

			if task.When.OutputChanges != "" {
				output := bytes.Buffer{}
				outputWriter := bufio.NewWriter(&output)
				status, err := runScript(task.When.OutputChanges, task.WorkingDirectory, task.Environment, nil, outputWriter, outputWriter)
				if err != nil {
					return fmt.Errorf("running task %v when.output_changes: %w", taskId, err)
				}
				err = outputWriter.Flush()
				if err != nil {
					return fmt.Errorf("running task %v when.check_fails: flushing writer: %w", taskId, err)
				}
				if status > 0 {
					return fmt.Errorf("task %v when.output_changes returned status code %v. here is the output:\n%v", taskId, status, output.String())
				}

				outputChanged, err := storeTaskOutputAndCheckIfChanged(taskId, output.Bytes())
				if err != nil {
					return fmt.Errorf("storing output for task %v when.output_changes: %w", taskId, err)
				}
				if outputChanged {
					skip = false
				}
			}
		}

		if skip {
			logger.Info(logLine(taskId, "skipping"))
			continue
		}

		var err error
		var status uint8
		if task.Quiet != nil && *task.Quiet {
			logger.Info(logLine(taskId, "running"))
			output := bytes.Buffer{}
			outputWriter := bufio.NewWriter(&output)
			for _, script := range task.Script {
				status, err = runScript(script, task.WorkingDirectory, task.Environment, nil, outputWriter, outputWriter)
				outputWriter.Flush()
				if err != nil || status != 0 {
					fmt.Print(output.String())
					break
				}
			}
		} else {
			logger.Notice(logLine(taskId, "running"))
			var stdin io.Reader = nil
			if task.Interactive != nil && *task.Interactive {
				stdin = os.Stdin
			}
			for _, script := range task.Script {
				status, err = runScript(script, task.WorkingDirectory, task.Environment, stdin, os.Stdout, os.Stdout)
				if err != nil || status != 0 {
					break
				}
			}
		}

		var final_err error

		if err != nil {
			final_err = fmt.Errorf("running task %v script: %w", taskId, err)
		}

		if status != 0 {
			final_err = fmt.Errorf("task %v returned status code %v", taskId, status)
		}

		if final_err != nil {
			err := removeTaskOutput(taskId)
			if err != nil {
				return fmt.Errorf("removing output after failure of task %v: %w", taskId, err)
			}
			return final_err
		}
	}
	return nil
}

func logLine(taskId core.TaskId, message string) string {
	return "[" + string(taskId) + "] " + message
}

func runScript(script string, workingDirectory string, environment *core.Environment, stdin io.Reader, stdout io.Writer, stderr io.Writer) (uint8, error) {
	script_header := []string{"set -euo pipefail"}

	file, err := syntax.NewParser().Parse(strings.NewReader(strings.Join(script_header, "\n")+"\n"+script), "")
	if err != nil {
		return 1, fmt.Errorf("parsing script: %w", err)
	}

	runner, err := interp.New(
		interp.Env(expand.ListEnviron(append(os.Environ(), environmentToString(environment)...)...)),
		interp.Dir(workingDirectory),
		interp.StdIO(stdin, stdout, stderr),
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

func environmentToString(environment *core.Environment) []string {
	result := []string{}
	for _, envValue := range environment.Values {
		result = append(result, envValue.Key+"="+envValue.Value)
	}
	return result
}

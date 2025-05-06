package loader

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirikon/ebro/internal/core"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func (ctx *loadCtx) moduleIterationResolvingPhase(module *core.Module) error {
	if module.ForEach != "" {
		var env = &core.Environment{}
		if module.Parent != nil {
			env = module.Parent.Environment
		}
		names, err := runForEachScript(module.ForEach, module.WorkingDirectory, env)
		if err != nil {
			return err
		}

		submodules := map[string]*core.Module{}

		for _, name := range names {
			// TODO: Prevent name collission if another submodule already exists called the same
			submodule := module.Clone(module)
			submodule.ForEach = ""
			submodule.Name = name
			submodule.Parent = module
			submodule.Environment.Values = append([]core.EnvironmentValue{{Key: "EBRO_EACH", Value: name}}, submodule.Environment.Values...)
			submodules[name] = submodule
		}

		module.ForEach = ""
		module.Environment = &core.Environment{}
		module.Imports = make(map[string]*core.Import)
		module.Labels = make(map[string]string)
		module.Tasks = make(map[string]*core.Task)
		module.Modules = submodules

		ctx.inventory.RefreshIndex()
	}
	return nil
}

func runForEachScript(script string, workingDirectory string, environment *core.Environment) ([]string, error) {
	output := bytes.Buffer{}
	outputWriter := bufio.NewWriter(&output)
	runner, err := interp.New(
		interp.Env(expand.ListEnviron(append(os.Environ(), environmentToString(environment)...)...)),
		interp.Dir(workingDirectory),
		interp.StdIO(nil, outputWriter, nil),
	)

	header_file, err := syntax.NewParser().Parse(strings.NewReader("set -euo pipefail"), "")
	if err != nil {
		return nil, fmt.Errorf("parsing header script: %w", err)
	}

	err = runner.Run(context.TODO(), header_file)
	if err != nil {
		if status, ok := interp.IsExitStatus(err); ok {
			return nil, fmt.Errorf("exited while applying header: (code %v) %w", status, err)
		} else {
			return nil, fmt.Errorf("error while applying header: %w", err)
		}
	}

	file, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		return nil, fmt.Errorf("parsing script: %w", err)
	}

	err = runner.Run(context.TODO(), file)
	if err != nil {
		if status, ok := interp.IsExitStatus(err); ok {
			return nil, fmt.Errorf("script exited with non-zero code: (code %v) %w", status, err)
		} else {
			return nil, fmt.Errorf("runner returned error: %w", err)
		}
	}

	err = outputWriter.Flush()
	if err != nil {
		return nil, fmt.Errorf("flushing stdout writer: %w", err)
	}

	output_text := string(output.Bytes())

	return strings.FieldsFunc(output_text, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n'
	}), nil
}

func environmentToString(environment *core.Environment) []string {
	result := []string{}
	for _, envValue := range environment.Values {
		result = append(result, envValue.Key+"="+envValue.Value)
	}
	return result
}

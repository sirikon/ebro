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
		var env = ctx.baseEnvironment
		if module.Parent != nil {
			env = module.Parent.Environment
		}
		names, err := runScriptReturnWords(module.ForEach, module.WorkingDirectory, env)
		if err != nil {
			return err
		}

		submodules := map[string]*core.Module{}

		for _, name := range names {
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

func runScriptReturnWords(script string, workingDirectory string, environment *core.Environment) ([]string, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	stdoutWriter := bufio.NewWriter(&stdout)
	stderrWriter := bufio.NewWriter(&stderr)
	runner, err := interp.New(
		interp.Env(expand.ListEnviron(append(os.Environ(), environmentToString(environment)...)...)),
		interp.Dir(workingDirectory),
		interp.StdIO(nil, stdoutWriter, stderrWriter),
	)
	if err != nil {
		return nil, fmt.Errorf("instantiating interpreter: %w", err)
	}

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
			if err = stdoutWriter.Flush(); err != nil {
				return nil, fmt.Errorf("flushing stdout writer: %w", err)
			}
			if err = stderrWriter.Flush(); err != nil {
				return nil, fmt.Errorf("flushing stderr writer: %w", err)
			}
			return nil, fmt.Errorf("script exited with code %v: stdout:\n%v███ stderr:\n%v", status, stdout.String(), stderr.String())
		} else {
			return nil, fmt.Errorf("runner returned error: %w", err)
		}
	}

	if err = stdoutWriter.Flush(); err != nil {
		return nil, fmt.Errorf("flushing stdout writer: %w", err)
	}

	return strings.FieldsFunc(stdout.String(), func(r rune) bool {
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

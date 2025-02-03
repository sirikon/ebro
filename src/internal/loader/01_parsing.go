package loader

import (
	"fmt"
	"iter"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/sirikon/ebro/internal/core2"
)

func (ctx *loadCtx) parsingPhase() error {
	var err error
	if ctx.inventory.RootModule, err = parseModuleFile(ctx.rootFile, []string{}); err != nil {
		return fmt.Errorf("parsing: %w", err)
	}
	return nil
}

func parseModuleFile(filePath string, modulePath []string) (*core2.Module, error) {
	file, err := parser.ParseFile(filePath, 0)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}

	result, err := parseModule(file.Docs[0].Body, modulePath)
	if err != nil {
		return nil, fmt.Errorf("parsing module: %w", err)
	}

	return result, nil
}

func parseModule(node ast.Node, modulePath []string) (*core2.Module, error) {
	var err error
	module := core2.NewModule()

	mapping, err := parseStringMapping(node)
	if err != nil {
		return nil, err
	}

	for key, value := range mapping {
		switch key {
		case "environment":
		case "imports":
		case "tasks":
			module.Tasks, err = parseTasks(value, modulePath)
		case "modules":
		default:
			return nil, fmt.Errorf("unexpected key %v", key)
		}
		if err != nil {
			return nil, fmt.Errorf("parsing value of %v", key)
		}
	}

	return module, nil
}

func parseTasks(node ast.Node, modulePath []string) (map[string]*core2.Task, error) {
	var err error
	tasks := map[string]*core2.Task{}

	mapping, err := parseStringMapping(node)
	if err != nil {
		return nil, err
	}

	for name, value := range mapping {
		if tasks[name], err = parseTask(value, modulePath, name); err != nil {
			return nil, fmt.Errorf("parsing task %v: %w", name, err)
		}
	}

	return tasks, nil
}

func parseTask(node ast.Node, modulePath []string, name string) (*core2.Task, error) {
	task := &core2.Task{}
	task.Name = name
	task.Id = core2.NewTaskId(modulePath, name)

	mapping, err := parseStringMapping(node)
	if err != nil {
		return nil, err
	}

	for key, _ := range mapping {
		switch key {
		case "labels":
		case "requires":
		case "required_by":
		case "script":
		case "interactive":
		case "quiet":
		case "when":
		default:
			return nil, fmt.Errorf("unexpected key %v", key)
		}
	}

	return task, nil
}

func parseStringMapping(node ast.Node) (iter.Seq2[string, ast.Node], error) {
	if node.Type() != ast.MappingType {
		return nil, fmt.Errorf("wrong type: %v", node.Type())
	}

	mapping := node.(*ast.MappingNode)
	for _, mappingValue := range mapping.Values {
		if mappingValue.Key.Type() != ast.StringType {
			return nil, fmt.Errorf("wrong type for key %v in mapping", mappingValue.Key)
		}
	}

	return func(yield func(string, ast.Node) bool) {
		for _, mappingValue := range mapping.Values {
			if !yield(mappingValue.Key.(*ast.StringNode).Value, mappingValue.Value) {
				return
			}
		}
	}, nil
}

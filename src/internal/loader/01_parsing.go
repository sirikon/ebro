package loader

import (
	"fmt"
	"iter"
	"path"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/sirikon/ebro/internal/core2"
)

func (ctx *loadCtx) parsingPhase() error {
	var err error
	if ctx.inventory.RootModule, err = parseModuleFile(ctx.rootFile, ctx.workingDirectory, []string{}); err != nil {
		return fmt.Errorf("parsing: %w", err)
	}
	ctx.inventory.RefreshIndex()
	return nil
}

func parseModuleFile(filePath string, workingDirectory string, modulePath []string) (*core2.Module, error) {
	file, err := parser.ParseFile(filePath, 0)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}

	result, err := parseModule(file.Docs[0].Body, workingDirectory, modulePath)
	if err != nil {
		return nil, fmt.Errorf("parsing module: %w", err)
	}

	return result, nil
}

func parseModule(node ast.Node, workingDirectory string, modulePath []string) (*core2.Module, error) {
	var err error
	module := &core2.Module{}
	module.Path = modulePath

	mapping, err := parseStringToAstMapping(node)
	if err != nil {
		return nil, err
	}

	for key, value := range mapping {
		switch key {
		case "working_directory":
			module.WorkingDirectory, err = parseString(value)
		case "environment":
			module.Environment, err = parseEnvironment(value)
		case "imports":
			module.Imports, err = parseImports(value, workingDirectory, modulePath)
		case "tasks":
			module.Tasks, err = parseTasks(value, modulePath)
		case "modules":
			module.Modules, err = parseModules(value, workingDirectory, modulePath)
		default:
			return nil, fmt.Errorf("unexpected key %v", key)
		}
		if err != nil {
			return nil, fmt.Errorf("parsing value of %v: %w", key, err)
		}
	}

	if module.Imports != nil {
		if module.Modules == nil {
			module.Modules = map[string]*core2.Module{}
		}
		for importName, importObj := range module.Imports {
			if _, ok := module.Modules[importName]; ok {
				return nil, fmt.Errorf("cannot import module %v because there is already a module called equally", importName)
			}
			module.Modules[importName] = importObj.Module
		}
	}

	module.Imports = nil

	if !path.IsAbs(module.WorkingDirectory) {
		module.WorkingDirectory = path.Join(workingDirectory, module.WorkingDirectory)
	}

	return module, nil
}

func parseImports(node ast.Node, workingDirectory string, modulePath []string) (map[string]*core2.Import, error) {
	var err error
	imports := map[string]*core2.Import{}

	mapping, err := parseStringToAstMapping(node)
	if err != nil {
		return nil, err
	}

	for name, value := range mapping {
		if err = core2.ValidateName(name); err != nil {
			return nil, fmt.Errorf("validating import name %v: %w", name, err)
		}
		if imports[name], err = parseImport(value, workingDirectory, append(modulePath, name)); err != nil {
			return nil, fmt.Errorf("parsing task %v: %w", name, err)
		}
	}

	return imports, nil
}

func parseImport(node ast.Node, workingDirectory string, modulePath []string) (*core2.Import, error) {
	var err error
	importObj := &core2.Import{}

	mapping, err := parseStringToAstMapping(node)
	if err != nil {
		return nil, err
	}

	for key, value := range mapping {
		switch key {
		case "from":
			importObj.From, err = parseString(value)
		case "environment":
			importObj.Environment, err = parseEnvironment(value)
		default:
			return nil, fmt.Errorf("unexpected key %v", key)
		}
		if err != nil {
			return nil, fmt.Errorf("parsing value of %v: %w", key, err)
		}
	}

	if !path.IsAbs(importObj.From) {
		importObj.From = path.Join(workingDirectory, importObj.From)
	}

	importObj.Module, err = parseModuleFile(path.Join(importObj.From, "Ebro.yaml"), importObj.From, modulePath)
	if err != nil {
		return nil, fmt.Errorf("importing: %w", err)
	}

	return importObj, nil
}

func parseModules(node ast.Node, workingDirectory string, modulePath []string) (map[string]*core2.Module, error) {
	var err error
	modules := map[string]*core2.Module{}

	mapping, err := parseStringToAstMapping(node)
	if err != nil {
		return nil, err
	}

	for name, value := range mapping {
		if err = core2.ValidateName(name); err != nil {
			return nil, fmt.Errorf("validating module name %v: %w", name, err)
		}
		if modules[name], err = parseModule(value, workingDirectory, append(modulePath, name)); err != nil {
			return nil, fmt.Errorf("parsing module %v: %w", name, err)
		}
	}

	return modules, nil
}

func parseTasks(node ast.Node, modulePath []string) (map[string]*core2.Task, error) {
	var err error
	tasks := map[string]*core2.Task{}

	mapping, err := parseStringToAstMapping(node)
	if err != nil {
		return nil, err
	}

	for name, value := range mapping {
		if err = core2.ValidateName(name); err != nil {
			return nil, fmt.Errorf("validating task name %v: %w", name, err)
		}
		if tasks[name], err = parseTask(value, modulePath, name); err != nil {
			return nil, fmt.Errorf("parsing task %v: %w", name, err)
		}
	}

	return tasks, nil
}

func parseTask(node ast.Node, modulePath []string, name string) (*core2.Task, error) {
	var err error
	task := &core2.Task{}
	task.Name = name
	task.Id = core2.NewTaskId(modulePath, name)

	mapping, err := parseStringToAstMapping(node)
	if err != nil {
		return nil, err
	}

	for key, value := range mapping {
		switch key {
		case "working_directory":
			task.WorkingDirectory, err = parseString(value)
		case "if_tasks_exist":
			task.IfTasksExist, err = parseStringSequence(value)
		case "labels":
			task.Labels, err = parseTaskLabels(value)
		case "requires":
			task.Requires, err = parseStringSequence(value)
		case "required_by":
			task.RequiredBy, err = parseStringSequence(value)
		case "abstract":
			task.Abstract, err = parseBoolPtr(value)
		case "extends":
			task.Extends, err = parseStringSequence(value)
		case "script":
			task.Script, err = parseString(value)
		case "interactive":
			task.Interactive, err = parseBoolPtr(value)
		case "quiet":
			task.Quiet, err = parseBoolPtr(value)
		case "when":
			task.When, err = parseWhen(value)
		case "environment":
			task.Environment, err = parseEnvironment(value)
		default:
			return nil, fmt.Errorf("unexpected key %v", key)
		}
		if err != nil {
			return nil, fmt.Errorf("parsing value of %v: %w", key, err)
		}
	}

	return task, nil
}

func parseTaskLabels(node ast.Node) (map[string]string, error) {
	result := map[string]string{}
	mapping, err := parseStringToAstMapping(node)
	if err != nil {
		return nil, err
	}

	for key, value := range mapping {
		if value.Type() != ast.StringType {
			return nil, fmt.Errorf("wrong type for value of %v in mapping: %v", key, value)
		}
		result[key] = value.(*ast.StringNode).Value
	}

	return result, nil
}

func parseWhen(node ast.Node) (*core2.When, error) {
	var err error
	when := &core2.When{}

	mapping, err := parseStringToAstMapping(node)
	if err != nil {
		return nil, err
	}

	for key, value := range mapping {
		switch key {
		case "output_changes":
			when.OutputChanges, err = parseString(value)
		case "check_fails":
			when.CheckFails, err = parseString(value)
		default:
			return nil, fmt.Errorf("unexpected key %v", key)
		}
		if err != nil {
			return nil, fmt.Errorf("parsing value of %v: %w", key, err)
		}
	}

	return when, nil
}

func parseEnvironment(node ast.Node) (*core2.Environment, error) {
	var err error
	environment := &core2.Environment{
		Values: []core2.EnvironmentValue{},
	}

	mapping, err := parseStringToAstMapping(node)
	if err != nil {
		return nil, err
	}

	for key, value := range mapping {
		if value.Type() != ast.StringType {
			return nil, fmt.Errorf("wrong type for value of %v in mapping: %v", key, value)
		}
		environment.Values = append(environment.Values, core2.EnvironmentValue{
			Key:   key,
			Value: value.(*ast.StringNode).Value,
		})
	}

	return environment, nil
}

func parseStringSequence(node ast.Node) ([]string, error) {
	result := []string{}

	if node.Type() != ast.SequenceType {
		return nil, fmt.Errorf("wrong type: %v", node.Type())
	}

	sequence := node.(*ast.SequenceNode)
	for i, value := range sequence.Values {
		if value.Type() != ast.StringType {
			return nil, fmt.Errorf("wrong type for item %v in sequence: %v", i, value)
		}
		result = append(result, value.(*ast.StringNode).Value)
	}

	return result, nil
}

func parseString(node ast.Node) (string, error) {
	switch node.Type() {
	case ast.StringType:
		return node.(*ast.StringNode).Value, nil
	case ast.LiteralType:
		return node.(*ast.LiteralNode).Value.Value, nil
	default:
		return "", fmt.Errorf("wrong type: %v", node.Type())
	}
}

func parseBoolPtr(node ast.Node) (*bool, error) {
	if node.Type() != ast.BoolType {
		return nil, fmt.Errorf("wrong type: %v", node.Type())
	}
	return &node.(*ast.BoolNode).Value, nil
}

func parseStringToAstMapping(node ast.Node) (iter.Seq2[string, ast.Node], error) {
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

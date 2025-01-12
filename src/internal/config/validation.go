package config

import (
	"fmt"
	"maps"
	"slices"
)

type rootModuleValidationContext struct {
	rootModule *Module
}

func ValidateRootModule(module *Module) error {
	ctx := rootModuleValidationContext{
		rootModule: module,
	}
	return ctx.validateModule(ctx.rootModule)
}

func (ctx *rootModuleValidationContext) validateModule(module *Module) error {
	taskNames := slices.Collect(maps.Keys(module.Tasks))
	slices.Sort(taskNames)

	for _, taskName := range taskNames {
		task := module.Tasks[taskName]
		err := ctx.validateTask(task)
		if err != nil {
			return fmt.Errorf("validating task %v: %w", taskName, err)
		}
	}

	moduleNames := slices.Collect(maps.Keys(module.Modules))
	slices.Sort(moduleNames)

	for _, moduleName := range moduleNames {
		module := module.Modules[moduleName]
		err := ctx.validateModule(module)
		if err != nil {
			return fmt.Errorf("validating module %v: %w", moduleName, err)
		}
	}

	return nil
}

func (ctx *rootModuleValidationContext) validateTask(task *Task) error {
	if len(task.Requires) == 0 && task.Script == "" && len(task.Extends) == 0 && !task.Abstract {
		return fmt.Errorf("task has nothing to do (no requires, script, extends nor abstract)")
	}

	for _, taskReferenceString := range task.Requires {
		_, err := ParseTaskReference(taskReferenceString)
		if err != nil {
			return fmt.Errorf("parsing reference %v in requires: %w", taskReferenceString, err)
		}
	}

	for _, taskReferenceString := range task.RequiredBy {
		_, err := ParseTaskReference(taskReferenceString)
		if err != nil {
			return fmt.Errorf("parsing reference %v in required_by: %w", taskReferenceString, err)
		}
	}

	for _, taskReferenceString := range task.Extends {
		_, err := ParseTaskReference(taskReferenceString)
		if err != nil {
			return fmt.Errorf("parsing reference %v in extends: %w", taskReferenceString, err)
		}
	}

	return nil
}

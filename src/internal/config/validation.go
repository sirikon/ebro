package config

import (
	"fmt"
)

func (t Task) Validate() error {
	if len(t.Requires) == 0 && t.Script == "" && len(t.Extends) == 0 && !t.Abstract {
		return fmt.Errorf("task has nothing to do (no requires, script, extends nor abstract)")
	}
	return nil
}

type rootModuleValidationContext struct {
	rootModule Module
}

func ValidateRootModule(module Module) error {
	ctx := rootModuleValidationContext{
		rootModule: module,
	}
	return ctx.validateModule(ctx.rootModule)
}

func (ctx *rootModuleValidationContext) validateModule(module Module) error {
	for taskName, task := range module.Tasks {
		err := ctx.validateTask(task)
		if err != nil {
			return fmt.Errorf("validating task %v: %w", taskName, err)
		}
	}
	for moduleName, module := range module.Modules {
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
		taskReference, err := ParseTaskReference(taskReferenceString)
		if err != nil {
			return fmt.Errorf("parsing reference %v in requires: %w", taskReferenceString, err)
		}

		requitedTask := ctx.rootModule.GetTask(taskReference)
		if requitedTask == nil && !taskReference.IsOptional {
			return fmt.Errorf("required task %v does not exist", taskReference.PartsString())
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

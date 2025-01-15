package config

import (
	"fmt"
)

type ctxValidateModule struct {
	module *Module
}

func ValidateRootModule(module *Module) error {
	ctx := ctxValidateModule{
		module: module,
	}
	return ctx.validateModule(ctx.module)
}

func (ctx *ctxValidateModule) validateModule(module *Module) error {
	for taskName, task := range module.TasksSorted() {
		if err := ctx.validateTask(task); err != nil {
			return fmt.Errorf("validating task %v: %w", taskName, err)
		}
	}

	for moduleName, module := range module.ModulesSorted() {
		if err := ctx.validateModule(module); err != nil {
			return fmt.Errorf("validating module %v: %w", moduleName, err)
		}
	}

	return nil
}

func (ctx *ctxValidateModule) validateTask(task *Task) error {
	if len(task.Requires) == 0 && task.Script == "" && len(task.Extends) == 0 && !task.Abstract {
		return fmt.Errorf("task has nothing to do (no requires, script, extends nor abstract)")
	}

	for _, taskReferenceString := range task.Requires {
		if err := ValidateTaskReference(taskReferenceString); err != nil {
			return fmt.Errorf("parsing reference %v in requires: %w", taskReferenceString, err)
		}
	}

	for _, taskReferenceString := range task.RequiredBy {
		if err := ValidateTaskReference(taskReferenceString); err != nil {
			return fmt.Errorf("parsing reference %v in required_by: %w", taskReferenceString, err)
		}
	}

	for _, taskReferenceString := range task.Extends {
		if err := ValidateTaskReference(taskReferenceString); err != nil {
			return fmt.Errorf("parsing reference %v in extends: %w", taskReferenceString, err)
		}
	}

	return nil
}

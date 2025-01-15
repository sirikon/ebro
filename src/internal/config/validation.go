package config

import (
	"fmt"

	"github.com/sirikon/ebro/internal/core"
)

type ctxValidateModule struct {
	module *Module
}

func ValidateModule(module *Module) error {
	ctx := ctxValidateModule{
		module: module,
	}
	return ctx.validateModule(ctx.module)
}

func (ctx *ctxValidateModule) validateModule(module *Module) error {
	for taskName, task := range module.TasksSorted() {
		if err := core.ValidateName(taskName); err != nil {
			return fmt.Errorf("validating task name '%v': %w", taskName, err)
		}
		if err := ctx.validateTask(task); err != nil {
			return fmt.Errorf("validating task %v: %w", taskName, err)
		}
	}

	for moduleName, module := range module.ModulesSorted() {
		if err := core.ValidateName(moduleName); err != nil {
			return fmt.Errorf("validating module name '%v': %w", moduleName, err)
		}
		if err := ctx.validateModule(module); err != nil {
			return fmt.Errorf("validating module %v: %w", moduleName, err)
		}
	}

	for importName, _ := range module.ImportsSorted() {
		if err := core.ValidateName(importName); err != nil {
			return fmt.Errorf("validating import name '%v': %w", importName, err)
		}
	}

	return nil
}

func (ctx *ctxValidateModule) validateTask(task *Task) error {
	if len(task.Requires) == 0 && task.Script == "" && len(task.Extends) == 0 && !task.Abstract {
		return fmt.Errorf("task has nothing to do (no requires, script, extends nor abstract)")
	}

	for _, taskReferenceString := range task.Requires {
		if err := validateTaskReference(taskReferenceString); err != nil {
			return fmt.Errorf("parsing reference %v in requires: %w", taskReferenceString, err)
		}
	}

	for _, taskReferenceString := range task.RequiredBy {
		if err := validateTaskReference(taskReferenceString); err != nil {
			return fmt.Errorf("parsing reference %v in required_by: %w", taskReferenceString, err)
		}
	}

	for _, taskReferenceString := range task.Extends {
		if err := validateTaskReference(taskReferenceString); err != nil {
			return fmt.Errorf("parsing reference %v in extends: %w", taskReferenceString, err)
		}
	}

	return nil
}

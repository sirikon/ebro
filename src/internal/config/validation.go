package config

import (
	"fmt"
	"strings"
)

func (t Task) Validate() error {
	if len(t.Requires) == 0 && t.Script == "" && len(t.Extends) == 0 && !t.Abstract {
		return fmt.Errorf("task has nothing to do (no requires, script, extends nor abstract)")
	}
	return nil
}

type rootModuleValidationContext struct {
	rootModule     Module
	validatedTasks map[string]bool
}

func ValidateRootModule(module Module) error {
	ctx := rootModuleValidationContext{rootModule: module, validatedTasks: make(map[string]bool)}
	return ctx.run()
}

func (ctx *rootModuleValidationContext) run() error {
	return nil
}

func (ctx *rootModuleValidationContext) validateTask(taskReference TaskReference) error {
	if ctx.validatedTasks[strings.Join(taskReference.Parts, ":")] {
		return nil
	}
	return nil
}

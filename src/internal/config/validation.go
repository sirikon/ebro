package config

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
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

type ctxValidateReferenceChain struct {
	rootModule *Module
}

func ValidateReferenceChain(module *Module, targets []string) error {
	ctx := ctxValidateReferenceChain{
		rootModule: module,
	}
	return ctx.run(targets)
}

func (ctx *ctxValidateReferenceChain) run(targets []string) error {
	for _, target := range targets {
		ref, _ := ParseTaskReference(target)
		err := ctx.checkReferenceChain(ref, []string{})
		if err != nil {
			return fmt.Errorf("checking reference chain: %w", err)
		}
	}
	return nil
}

func (ctx *ctxValidateReferenceChain) checkReferenceChain(taskReference TaskReference, idTrail []string) error {
	taskId, task := ctx.rootModule.GetTask(taskReference)
	if taskId == nil {
		return WrapErrTaskNotFound(taskReference)
	}

	if slices.Contains(idTrail, taskId.String()) {
		return fmt.Errorf("cyclic reference detected:\n%v", strings.Join(append(idTrail, taskId.String()), " -> "))
	}
	newIdTrail := append(idTrail, taskId.String())

	for _, taskReferenceString := range task.Requires {
		taskReference, err := ParseTaskReference(taskReferenceString)
		if err != nil {
			return fmt.Errorf("parsing task reference %v: %w", taskReferenceString, err)
		}
		taskReference = taskReference.Absolute(taskId.ModuleTrail)

		err = ctx.checkReferenceChain(taskReference, newIdTrail)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) && taskReference.IsOptional {
				continue
			}
			return fmt.Errorf("requires %v: %w", taskReference.String(), err)
		}
	}

	for _, taskReferenceString := range task.RequiredBy {
		taskReference, err := ParseTaskReference(taskReferenceString)
		if err != nil {
			return fmt.Errorf("parsing task reference %v: %w", taskReferenceString, err)
		}
		taskReference = taskReference.Absolute(taskId.ModuleTrail)

		err = ctx.checkReferenceChain(taskReference, newIdTrail)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) && taskReference.IsOptional {
				continue
			}
			return fmt.Errorf("required by %v: %w", taskReference.String(), err)
		}
	}

	for _, taskReferenceString := range task.Extends {
		taskReference, err := ParseTaskReference(taskReferenceString)
		if err != nil {
			return fmt.Errorf("parsing task reference %v: %w", taskReferenceString, err)
		}
		taskReference = taskReference.Absolute(taskId.ModuleTrail)

		err = ctx.checkReferenceChain(taskReference, newIdTrail)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) && taskReference.IsOptional {
				continue
			}
			return fmt.Errorf("extends %v: %w", taskReference.String(), err)
		}
	}

	return nil
}

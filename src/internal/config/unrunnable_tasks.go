package config

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type ctxGetUnrunnableTasks struct {
	rootModule *Module
	result     map[string]error
}

func GetUnrunnableTasks(module *Module) map[string]error {
	ctx := ctxGetUnrunnableTasks{
		rootModule: module,
		result:     map[string]error{},
	}
	ctx.run()
	return ctx.result
}

func (ctx *ctxGetUnrunnableTasks) run() {
	ctx.processModule(ctx.rootModule, []string{})
}

func (ctx *ctxGetUnrunnableTasks) processModule(module *Module, moduleTrail []string) {
	taskNames := slices.Collect(maps.Keys(module.Tasks))
	slices.Sort(taskNames)

	for _, taskName := range taskNames {
		ref := TaskReference{Path: append(moduleTrail, taskName), IsRelative: false, IsOptional: false}
		err := ctx.checkReferenceChain(TaskReference{Path: append(moduleTrail, taskName), IsRelative: false, IsOptional: false}, []string{})
		if err != nil {
			ctx.result[ref.String()] = err
		}
	}

	moduleNames := slices.Collect(maps.Keys(module.Modules))
	slices.Sort(moduleNames)

	for _, moduleName := range moduleNames {
		module := module.Modules[moduleName]
		ctx.processModule(module, append(moduleTrail, moduleName))
	}
}

func (ctx *ctxGetUnrunnableTasks) checkReferenceChain(taskReference TaskReference, idTrail []string) error {
	taskId, task := ctx.rootModule.GetTask(taskReference)
	if taskId == nil {
		return ErrTaskNotFound
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

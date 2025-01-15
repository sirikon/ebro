package config

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/sirikon/ebro/internal/core"
)

type ctxNormalizeModule struct {
	rootModule *RootModule
}

func NormalizeRootModule(rootModule *RootModule) (*core.Module, error) {
	ctx := ctxNormalizeModule{
		rootModule: rootModule,
	}
	err := ctx.normalizeModule(ctx.rootModule.Module, []string{})
	if err != nil {
		return nil, err
	}
	return toCoreModulePtr(ctx.rootModule.Module), nil
}

func (ctx *ctxNormalizeModule) normalizeModule(module *Module, moduleTrail []string) error {
	for taskName, task := range module.TasksSorted() {
		if err := ctx.normalizeTask(task, moduleTrail); err != nil {
			return fmt.Errorf("normalizing task '%v': %w", taskName, err)
		}
	}

	for moduleName, module := range module.ModulesSorted() {
		if err := ctx.normalizeModule(module, append(moduleTrail, moduleName)); err != nil {
			return fmt.Errorf("normalizing module '%v': %w", moduleName, err)
		}
	}

	return nil
}

func (ctx *ctxNormalizeModule) normalizeTask(task *Task, moduleTrail []string) error {
	var err error
	task.Requires, err = ctx.resolveRefs(task.Requires, moduleTrail)
	if err != nil {
		return fmt.Errorf("resolving 'requires': %w", err)
	}

	task.RequiredBy, err = ctx.resolveRefs(task.RequiredBy, moduleTrail)
	if err != nil {
		return fmt.Errorf("resolving 'required_by': %w", err)
	}

	task.Extends, err = ctx.resolveRefs(task.Extends, moduleTrail)
	if err != nil {
		return fmt.Errorf("resolving 'extends': %w", err)
	}

	return nil
}

func (ctx *ctxNormalizeModule) resolveRefs(s []string, moduleTrail []string) ([]string, error) {
	result := []string{}
	for _, taskReferenceString := range s {
		ref := MustParseTaskReference(taskReferenceString)
		taskId, _ := FindTask(ctx.rootModule, ref.Absolute(moduleTrail))
		if taskId != nil {
			result = append(result, string(*taskId))
		} else if !ref.IsOptional {
			return nil, fmt.Errorf("referenced task %v does not exist", ref.String())
		}
	}
	return result, nil
}

func toCoreModulePtr(module *Module) *core.Module {
	result := &core.Module{}
	castUsingYaml(module, result)
	return result
}

func castUsingYaml(from interface{}, to interface{}) {
	data, err := yaml.Marshal(from)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(data, to)
}

func NormalizeTargets(rootModule *RootModule, targets []string) ([]core.TaskId, error) {
	result := []core.TaskId{}
	for _, target := range targets {
		taskId, err := normalizeTarget(rootModule, target)
		if err != nil {
			return nil, fmt.Errorf("validating target '%v': %w", target, err)
		}
		result = append(result, *taskId)
	}
	return result, nil
}

func normalizeTarget(rootModule *RootModule, target string) (*core.TaskId, error) {
	if err := ValidateTaskReference(target); err != nil {
		return nil, err
	}
	ref := MustParseTaskReference(target).Absolute([]string{})
	taskId, _ := FindTask(rootModule, ref)
	if taskId == nil {
		return nil, fmt.Errorf("task does not exist")
	}
	return taskId, nil
}

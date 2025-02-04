package loader

import (
	"fmt"

	"github.com/sirikon/ebro/internal/core2"
)

func (ctx *loadCtx) referenceResolvingPhase() error {
	var err error

	for task := range ctx.inventory.Tasks() {
		if task.RequiresIds, err = resolveReferences(ctx.inventory, task, task.Requires); err != nil {
			return fmt.Errorf("normalizing 'requires' for task %v: %w", task.Id, err)
		}
		if task.RequiredByIds, err = resolveReferences(ctx.inventory, task, task.RequiredBy); err != nil {
			return fmt.Errorf("normalizing 'required_by' for task %v: %w", task.Id, err)
		}
		if task.ExtendsIds, err = resolveReferences(ctx.inventory, task, task.Extends); err != nil {
			return fmt.Errorf("normalizing 'extends' for task %v: %w", task.Id, err)
		}
	}

	return nil
}

func resolveReferences(inventory *core2.Inventory, task *core2.Task, taskReferences []string) ([]core2.TaskId, error) {
	result := []core2.TaskId{}
	for _, taskReference := range taskReferences {
		if err := core2.ValidateTaskReference(taskReference); err != nil {
			return nil, fmt.Errorf("validating '%v': %w", taskReference, err)
		}

		ref := core2.MustParseTaskReference(taskReference)
		if ref.IsRelative {
			ref = ref.Absolute(task.Id.ModulePath())
		}

		referencedTaskId, _ := inventory.FindTask(ref)
		if referencedTaskId == nil {
			if ref.IsOptional {
				continue
			} else {
				return nil, fmt.Errorf("referenced task %v does not exist", ref.TaskId())
			}
		}

		result = append(result, *referencedTaskId)
	}
	return result, nil
}

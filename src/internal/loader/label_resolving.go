package loader

import (
	"fmt"
	"strings"

	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/utils"
)

func (ctx *loadCtx) taskLabelResolvingPhase(taskId core.TaskId) error {
	task := ctx.inventory.Task(taskId)
	var err error

	for label, value := range task.Labels {
		task.Labels[label], err = utils.ExpandString(value, task.Environment)
		if err != nil {
			return fmt.Errorf("expanding label %v in task %v: %w", label, task.Id, err)
		}
	}

	return nil
}

func (ctx *loadCtx) moduleLabelResolvingPhase() error {
	for module := range ctx.inventory.Modules() {
		moduleEnvironment, err := resolveModuleEnvironment(ctx.inventory, ctx.baseEnvironment, module)
		if err != nil {
			return fmt.Errorf("resolving module '%v' environment: %w", ":"+strings.Join(module.Path, ":"), err)
		}

		for label, value := range module.Labels {
			module.Labels[label], err = utils.ExpandString(value, moduleEnvironment)
			if err != nil {
				return fmt.Errorf("expanding label %v in module %v: %w", label, ":"+strings.Join(module.Path, ":"), err)
			}
		}
	}
	return nil
}

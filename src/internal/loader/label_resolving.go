package loader

import (
	"fmt"
	"strings"

	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/utils"
)

func (ctx *loadCtx) moduleLabelResolvingPhase(module *core.Module) error {
	var err error
	for label, value := range module.Labels {
		module.Labels[label], err = utils.ExpandString(value, module.Environment)
		if err != nil {
			return fmt.Errorf("expanding label %v in module %v: %w", label, ":"+strings.Join(module.Path(), ":"), err)
		}
	}
	return nil
}

func (ctx *loadCtx) taskLabelResolvingPhase() error {
	for task := range ctx.inventory.Tasks() {
		var err error
		for label, value := range task.Labels {
			task.Labels[label], err = utils.ExpandString(value, task.Environment)
			if err != nil {
				return fmt.Errorf("expanding label %v in task %v: %w", label, task.Id, err)
			}
		}
	}
	return nil
}

package loader

import (
	"fmt"

	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/utils"
)

func (ctx *loadCtx) labelResolvingPhase(task *core.Task) error {
	var err error

	for label, value := range task.Labels {
		task.Labels[label], err = utils.ExpandString(value, task.Environment)
		if err != nil {
			return fmt.Errorf("expanding label %v in task %v: %w", label, task.Id, err)
		}
	}

	return nil
}

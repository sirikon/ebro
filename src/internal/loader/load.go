package loader

import (
	"github.com/sirikon/ebro/internal/core"
)

type loadCtx struct {
	workingDirectory string
	rootFile         string
	baseEnvironment  *core.Environment
	inventory        *core.Inventory
}

type phase = func() error
type taskPhase = func(core.TaskId) error

func Load(baseEnvironment *core.Environment, workingDirectory string, rootFile string) (*core.Inventory, error) {
	ctx := &loadCtx{
		workingDirectory: workingDirectory,
		rootFile:         rootFile,
		baseEnvironment:  baseEnvironment,
		inventory:        core.NewInventory(),
	}

	phases := []phase{
		ctx.parsingPhase,
		ctx.conditionalExistencePurgingPhase,
		ctx.workdirResolvingPhase,
		ctx.moduleLabelResolvingPhase,
		ctx.extensionReferenceResolvingPhase,
		ctx.perTaskByExtensionOrder(
			ctx.requirementReferenceResolvingPhase,
			ctx.extendingPhase,
			ctx.environmentResolvingPhase,
			ctx.taskLabelResolvingPhase,
		),
		ctx.abstractPurgingPhase,
		ctx.requirementExpressionReferenceResolvingPhase,
	}

	for _, phase := range phases {
		if err := phase(); err != nil {
			return nil, err
		}
	}

	return ctx.inventory, nil
}

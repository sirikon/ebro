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
type modulePhase = func(*core.Module) error
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
		ctx.perModuleByHierarchicalOrder(
			ctx.workdirResolvingPhase,
			ctx.moduleEnvironmentResolvingPhase,
			ctx.moduleIterationResolvingPhase,
			ctx.moduleLabelResolvingPhase,
		),
		ctx.conditionalExistencePurgingPhase,
		ctx.extensionReferenceResolvingPhase,
		ctx.perTaskByExtensionOrder(
			ctx.requirementReferenceResolvingPhase,
			ctx.extendingPhase,
			ctx.taskEnvironmentResolvingPhase,
		),
		ctx.taskLabelResolvingPhase,
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

func (ctx *loadCtx) perModuleByHierarchicalOrder(modulePhases ...modulePhase) phase {
	return func() error {
		for module := range ctx.inventory.WalkDownModuleTree() {
			for _, modulePhase := range modulePhases {
				if err := modulePhase(module); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

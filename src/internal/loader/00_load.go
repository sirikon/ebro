package loader

import (
	"fmt"

	"github.com/sirikon/ebro/internal/core"
)

type loadCtx struct {
	workingDirectory string
	rootFile         string
	baseEnvironment  *core.Environment
	inventory        *core.Inventory
}

type phase = func() error

func Load(baseEnvironment *core.Environment, workingDirectory string, rootFile string) (*core.Inventory, error) {
	ctx := &loadCtx{
		workingDirectory: workingDirectory,
		rootFile:         rootFile,
		baseEnvironment:  baseEnvironment,
		inventory:        core.NewInventory(),
	}

	phases := []phase{
		ctx.parsingPhase,
		ctx.purgingPhase,
		ctx.referenceResolvingPhase,
		ctx.workdirResolvingPhase,
		ctx.extendingPhase,
		ctx.labelResolvingPhase,
	}

	for _, phase := range phases {
		if err := phase(); err != nil {
			return nil, fmt.Errorf("loading: %w", err)
		}
	}

	return ctx.inventory, nil
}

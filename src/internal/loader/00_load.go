package loader

import (
	"fmt"

	"github.com/sirikon/ebro/internal/core2"
)

type loadCtx struct {
	workingDirectory string
	rootFile         string
	inventory        *core2.Inventory
}

type phase = func() error

func Load(workingDirectory string, rootFile string) (*core2.Inventory, error) {
	ctx := &loadCtx{
		workingDirectory: workingDirectory,
		rootFile:         rootFile,
		inventory:        core2.NewInventory(),
	}

	phases := []phase{
		ctx.parsingPhase,
		ctx.purgingPhase,
	}

	for _, phase := range phases {
		if err := phase(); err != nil {
			return nil, fmt.Errorf("loading: %w", err)
		}
	}

	return ctx.inventory, nil
}

package config2

import (
	"fmt"

	"github.com/sirikon/ebro/internal/core2"
)

type LoadCtx struct {
	WorkingDirectory string
	RootFile         string
}

func Load(ctx LoadCtx) (*core2.Inventory, error) {
	var err error
	inventory := core2.NewInventory()

	if inventory.RootModule, err = parseModuleFile(ctx.RootFile); err != nil {
		return nil, fmt.Errorf("loading: %w", err)
	}

	return inventory, nil
}

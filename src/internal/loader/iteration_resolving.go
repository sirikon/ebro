package loader

import "github.com/sirikon/ebro/internal/core"

func (ctx *loadCtx) moduleIterationResolvingPhase(module *core.Module) error {
	if module.ForEach != "" {
		names, err := runForEachScript(module.ForEach)
		if err != nil {
			return err
		}

		for _, name := range names {
			// TODO: Prevent name collission if another submodule already exists called the same
			submodule := module.Clone()
			submodule.ForEach = ""
			submodule.Name = name
			submodule.Parent = module
			if module.Modules == nil {
				module.Modules = make(map[string]*core.Module)
			}
			module.Modules[name] = submodule
		}

		module.Tasks = make(map[string]*core.Task)

		ctx.inventory.RefreshIndex()
	}
	return nil
}

func runForEachScript(script string) ([]string, error) {
	return []string{"a", "b"}, nil
}

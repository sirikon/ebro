package loader

import "github.com/sirikon/ebro/internal/core"

func (ctx *loadCtx) moduleIterationResolvingPhase(module *core.Module) error {
	if module.ForEach != "" {
		names, err := runForEachScript(module.ForEach)
		if err != nil {
			return err
		}

		submodules := map[string]*core.Module{}

		for _, name := range names {
			// TODO: Prevent name collission if another submodule already exists called the same
			submodule := module.Clone(module)
			submodule.ForEach = ""
			submodule.Name = name
			submodule.Parent = module
			submodule.Environment.Values = append([]core.EnvironmentValue{{Key: "EBRO_EACH", Value: name}}, submodule.Environment.Values...)
			submodules[name] = submodule
		}

		module.ForEach = ""
		module.Environment = &core.Environment{}
		module.Imports = make(map[string]*core.Import)
		module.Labels = make(map[string]string)
		module.Tasks = make(map[string]*core.Task)
		module.Modules = submodules

		ctx.inventory.RefreshIndex()
	}
	return nil
}

func runForEachScript(script string) ([]string, error) {
	return []string{"a", "b"}, nil
}

package config

type ctxPurgeModule struct {
	rootModule *Module
}

func PurgeModule(module *Module) {
	ctx := ctxPurgeModule{
		rootModule: module,
	}
	ctx.processModule(module, []string{})
}

func (ctx *ctxPurgeModule) processModule(module *Module, moduleTrail []string) {
	for taskName, task := range module.TasksSorted() {
		if len(task.IfTasksExist) > 0 {
			purge := false
			for _, t := range task.IfTasksExist {
				ref, _ := ParseTaskReference(t)
				ref = ref.Absolute(moduleTrail)
				taskId, _ := ctx.rootModule.GetTask(ref)
				if taskId == nil {
					purge = true
				}
			}
			if purge {
				delete(module.Tasks, taskName)
			} else {
				task.IfTasksExist = []string{}
			}
		}
	}

	for submoduleName, submodule := range module.ModulesSorted() {
		ctx.processModule(submodule, append(moduleTrail, submoduleName))
	}
}

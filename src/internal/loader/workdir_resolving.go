package loader

import (
	"path"
	"slices"

	"github.com/sirikon/ebro/internal/core"
)

func (ctx *loadCtx) workdirResolvingPhase(module *core.Module) error {
	for _, task := range module.Tasks {
		workDirs := []string{task.WorkingDirectory}
		for _, module := range ctx.inventory.WalkUpModulePath(task.Module.Path()) {
			workDirs = append(workDirs, module.WorkingDirectory)
		}
		slices.Reverse(workDirs)

		currentWorkDir := ""
		for _, workDir := range workDirs {
			if path.IsAbs(workDir) {
				currentWorkDir = workDir
			} else {
				currentWorkDir = path.Join(currentWorkDir, workDir)
			}
		}

		task.WorkingDirectory = currentWorkDir
	}
	return nil
}

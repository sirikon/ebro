package loader

import (
	"path"
	"slices"
)

func (ctx *loadCtx) workdirResolvingPhase() error {
	for task := range ctx.inventory.Tasks() {
		workDirs := []string{task.WorkingDirectory}
		for module := range ctx.inventory.WalkUpModulePath(task.Id) {
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

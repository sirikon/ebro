package planner

import (
	"fmt"

	"github.com/sirikon/ebro/internal/core"
	"github.com/sirikon/ebro/internal/dag"
	"github.com/sirikon/ebro/internal/utils"

	"github.com/goccy/go-yaml"
)

type Plan []core.TaskId

func MakePlan(inv *core.Inventory, targets []core.TaskId) (Plan, error) {
	tasksToRun := utils.NewSet[core.TaskId]()
	taskDag := dag.NewDag[core.TaskId]()

	tasksToRun.Add(targets...)

	for i := 0; i < tasksToRun.Length(); i++ {
		taskName := tasksToRun.Get(i)
		task := inv.Task(taskName)
		if task == nil {
			return nil, fmt.Errorf("task %v does not exist", taskName)
		}
		tasksToRun.Add(task.RequiresIds...)
		taskDag.Link(taskName, task.RequiresIds...)
		for _, requiredByTaskName := range task.RequiredByIds {
			taskDag.Link(requiredByTaskName, taskName)
		}
	}

	result, remains := taskDag.Resolve(targets)
	if remains != nil {
		remainsData, err := yaml.Marshal(remains)
		if err != nil {
			return nil, fmt.Errorf("planning could not complete. error while turning requirement index to yaml: %w", err)
		}
		return nil, fmt.Errorf("planning could not complete. "+
			"there could be a cyclic dependency. "+
			"here is the list of tasks remaining to be planned and their requirements:\n%s", string(remainsData))
	}

	return result, nil
}

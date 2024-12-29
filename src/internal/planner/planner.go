package planner

import (
	"fmt"

	"github.com/sirikon/ebro/internal/dag"
	"github.com/sirikon/ebro/internal/inventory"
	"github.com/sirikon/ebro/utils"

	"gopkg.in/yaml.v3"
)

type Plan []string

func MakePlan(inv inventory.Inventory, targets []string) (Plan, error) {
	tasksToRun := utils.NewSet[string]()
	taskDag := dag.NewDag()

	tasksToRun.Add(targets...)

	for i := 0; i < tasksToRun.Length(); i++ {
		taskName := tasksToRun.Get(i)
		task, ok := inv.Tasks[taskName]
		if !ok {
			return nil, fmt.Errorf("task %v does not exist", taskName)
		}
		tasksToRun.Add(task.Requires...)
		taskDag.Link(taskName, task.Requires...)
		for _, requiredByTaskName := range task.RequiredBy {
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

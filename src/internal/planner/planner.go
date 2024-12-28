package planner

import (
	"fmt"

	"github.com/sirikon/ebro/internal/dag"
	"github.com/sirikon/ebro/internal/inventory"

	"gopkg.in/yaml.v3"
)

type Plan []string

func MakePlan(inv inventory.Inventory, targets []string) (Plan, error) {
	taskDag := dag.NewDag()

	for taskName, task := range inv.Tasks {
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

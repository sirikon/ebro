package planner

import (
	"fmt"

	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/dag"
	"github.com/sirikon/ebro/internal/inventory"
	"github.com/sirikon/ebro/internal/utils"

	"github.com/goccy/go-yaml"
)

type Plan []string

func MakePlan(inv inventory.Inventory, unrunnableTasks map[string]error, targets []string) (Plan, error) {
	tasksToRun := utils.NewSet[string]()
	taskDag := dag.NewDag()

	tasksToRun.Add(targets...)

	checkTaskCanRun := func(taskReferenceString string) (bool, error) {
		ref, _ := config.ParseTaskReference(taskReferenceString)
		if err := unrunnableTasks[ref.PathString()]; err != nil {
			if !ref.IsOptional {
				return false, fmt.Errorf("task %v cannot run: %w", taskReferenceString, err)
			}
			return false, nil
		}
		return true, nil
	}

	for i := 0; i < tasksToRun.Length(); i++ {
		taskName := tasksToRun.Get(i)
		canRun, err := checkTaskCanRun(taskName)
		if err != nil {
			return nil, err
		}
		if !canRun {
			continue
		}

		task, ok := inv.Tasks[taskName]
		if !ok {
			return nil, fmt.Errorf("task %v does not exist", taskName)
		}

		for _, taskReferenceString := range task.Requires {
			canRun, err := checkTaskCanRun(taskReferenceString)
			if err != nil {
				return nil, err
			}
			if !canRun {
				continue
			}

			ref, _ := config.ParseTaskReference(taskReferenceString)
			tasksToRun.Add(ref.PathString())
			taskDag.Link(taskName, ref.PathString())
		}

		for _, taskReferenceString := range task.RequiredBy {
			ref, _ := config.ParseTaskReference(taskReferenceString)
			taskDag.Link(ref.PathString(), taskName)
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

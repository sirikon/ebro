package planner

import (
	"fmt"
	"slices"

	"github.com/sirikon/ebro/internal/inventory2"
	"gopkg.in/yaml.v3"
)

type Plan []string

func MakePlan(inv inventory2.Inventory, targets []string) (Plan, error) {
	result := Plan{}
	tasksToRun := []string{}
	requirementsIndex := make(map[string][]string)

	addTasksToRun := func(taskNames ...string) {
		for _, taskName := range taskNames {
			if i := slices.Index(tasksToRun, taskName); i == -1 {
				tasksToRun = append(tasksToRun, taskName)
				if _, ok := requirementsIndex[taskName]; !ok {
					requirementsIndex[taskName] = []string{}
				}
			}
		}
	}

	addRequirements := func(task string, requirements ...string) {
		for _, requirement := range requirements {
			if i := slices.Index(requirementsIndex[task], requirement); i == -1 {
				requirementsIndex[task] = append(requirementsIndex[task], requirement)
			}
		}
	}

	for _, target := range targets {
		if _, ok := inv[target]; !ok {
			return nil, fmt.Errorf("target task %v does not exist", target)
		}
		addTasksToRun(target)
	}

	for i := 0; i < len(tasksToRun); i++ {
		taskName := tasksToRun[i]
		task := inv[taskName]
		addTasksToRun(task.Requires...)
		addRequirements(taskName, task.Requires...)
		for _, parent := range task.RequiredBy {
			addRequirements(parent, taskName)
		}
	}

	shouldContinue := true
	for shouldContinue {
		shouldContinue = false
		batch := []string{}
		for name, requires := range requirementsIndex {
			if len(requires) == 0 {
				batch = append(batch, name)
				shouldContinue = true
			}
		}

		slices.Sort(batch)
		for _, name := range batch {
			delete(requirementsIndex, name)
			for parent := range requirementsIndex {
				i := slices.Index(requirementsIndex[parent], name)
				if i >= 0 {
					requirementsIndex[parent] = append(requirementsIndex[parent][:i], requirementsIndex[parent][i+1:]...)
				}
			}
		}

		result = append(result, batch...)
	}

	if len(requirementsIndex) > 0 {
		detail, err := yaml.Marshal(requirementsIndex)
		if err != nil {
			return nil, fmt.Errorf("planning could not complete. error while turning requirement index to yaml: %w", err)
		}
		return nil, fmt.Errorf("planning could not complete. " +
			"there could be a cyclic dependency. " +
			"here is the list of tasks remaining to be planned and their requirements:\n" + string(detail))
	}

	return result, nil
}

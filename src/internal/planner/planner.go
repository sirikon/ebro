package planner

import (
	"fmt"
	"maps"
	"slices"

	"github.com/sirikon/ebro/internal/inventory"
	"gopkg.in/yaml.v3"
)

type Plan []string

func MakePlan(inv inventory.Inventory, targets []string) (Plan, error) {
	result := Plan{}
	tasksToRun := make(map[string]bool)
	tasksToRunSlice := []string{}
	requirementsIndex := make(map[string]map[string]bool)

	addTasksToRun := func(taskNames ...string) {
		for _, taskName := range taskNames {
			if _, ok := tasksToRun[taskName]; !ok {
				tasksToRun[taskName] = true
				tasksToRunSlice = append(tasksToRunSlice, taskName)
			}
		}
	}

	addRequirements := func(taskName string, requirementsToAdd ...string) {
		requirements, ok := requirementsIndex[taskName]
		if !ok {
			requirements = make(map[string]bool)
			requirementsIndex[taskName] = requirements
		}
		for _, requirementToAdd := range requirementsToAdd {
			requirements[requirementToAdd] = true
		}
	}

	for _, target := range targets {
		addTasksToRun(target)
	}

	for i := 0; i < len(tasksToRunSlice); i++ {
		taskName := tasksToRunSlice[i]
		task, ok := inv.Tasks[taskName]
		if !ok {
			return nil, fmt.Errorf("task %v does not exist", taskName)
		}
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

		for taskName := range tasksToRun {
			requires := requirementsIndex[taskName]
			if len(requires) == 0 {
				batch = append(batch, taskName)
				shouldContinue = true
			}
		}

		slices.Sort(batch)

		for _, taskName := range batch {
			delete(tasksToRun, taskName)
			delete(requirementsIndex, taskName)
			for parentTaskName := range requirementsIndex {
				delete(requirementsIndex[parentTaskName], taskName)
			}
		}

		result = append(result, batch...)
	}

	if len(tasksToRun) > 0 {
		detailObj := make(map[string][]string)
		for parentTaskName := range tasksToRun {
			taskNames := slices.Collect(maps.Keys(requirementsIndex[parentTaskName]))
			slices.Sort(taskNames)
			detailObj[parentTaskName] = taskNames
		}
		detail, err := yaml.Marshal(detailObj)
		if err != nil {
			return nil, fmt.Errorf("planning could not complete. error while turning requirement index to yaml: %w", err)
		}
		return nil, fmt.Errorf("planning could not complete. "+
			"there could be a cyclic dependency. "+
			"here is the list of tasks remaining to be planned and their requirements:\n%s", string(detail))
	}

	return result, nil
}

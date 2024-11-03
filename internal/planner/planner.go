package planner

import (
	"fmt"
	"slices"

	"github.com/sirikon/ebro/internal/cataloger"
	"gopkg.in/yaml.v3"
)

type Plan []string

func MakePlan(catalog cataloger.Catalog, targets []string) (Plan, error) {
	result := Plan{}
	tasksToRun := []string{}
	requirementsIndex := make(map[string][]string)

	addTasksToRun := func(task_names ...string) {
		for _, task_name := range task_names {
			if i := slices.Index(tasksToRun, task_name); i == -1 {
				tasksToRun = append(tasksToRun, task_name)
				if _, ok := requirementsIndex[task_name]; !ok {
					requirementsIndex[task_name] = []string{}
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
		if _, ok := catalog[target]; !ok {
			return nil, fmt.Errorf("target task %v does not exist", target)
		}
		addTasksToRun(target)
	}

	for i := 0; i < len(tasksToRun); i++ {
		task_name := tasksToRun[i]
		task := catalog[task_name]
		addTasksToRun(task.Requires...)
		addRequirements(task_name, task.Requires...)
		for _, parent := range task.RequiredBy {
			addRequirements(parent, task_name)
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

package planner

import (
	"slices"

	"github.com/sirikon/ebro/internal/cataloger"
)

type Plan []string

func MakePlan(catalog cataloger.Catalog, targets []string) Plan {
	result := Plan{}
	tasksToRun := []string{}
	reqIndex := make(map[string][]string)

	for _, target := range targets {
		if _, ok := catalog[target]; ok {
			tasksToRun = append(tasksToRun, target)
		}
	}

	for i := 0; i < len(tasksToRun); i++ {
		task_name := tasksToRun[i]
		task := catalog[task_name]
		tasksToRun = append(tasksToRun, task.Requires...)
		reqIndex[task_name] = append(reqIndex[task_name], task.Requires...)
		for _, parent := range task.RequiredBy {
			reqIndex[parent] = append(reqIndex[parent], task_name)
		}
	}

	shouldContinue := true
	for shouldContinue {
		shouldContinue = false
		batch := []string{}
		for name, requires := range reqIndex {
			if len(requires) == 0 {
				batch = append(batch, name)
				shouldContinue = true
			}
		}

		slices.Sort(batch)
		for _, name := range batch {
			delete(reqIndex, name)
			for parent := range reqIndex {
				i := slices.Index(reqIndex[parent], name)
				if i >= 0 {
					reqIndex[parent] = append(reqIndex[parent][:i], reqIndex[parent][i+1:]...)
				}
			}
		}

		result = append(result, batch...)
	}

	return result
}

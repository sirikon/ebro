package planner

import (
	"slices"

	"github.com/sirikon/ebro/internal/indexer"
)

type Plan struct {
	Steps []string
}

type Step struct {
	Requires   []string
	RequiredBy []string
}

func MakePlan(index indexer.Index) Plan {
	result := Plan{}

	reqIndex := make(map[string][]string)

	for task_name, task := range index {
		reqIndex[task_name] = append(reqIndex[task_name], task.Requires...)
		for _, parent := range task.RequiredBy {
			reqIndex[parent] = append(reqIndex[parent], task_name)
		}
	}

	shouldContinue := true
	for shouldContinue {
		shouldContinue = false
		for name, requires := range reqIndex {
			if len(requires) == 0 {
				result.Steps = append(result.Steps, name)
				delete(reqIndex, name)
				for parent := range reqIndex {
					i := slices.Index(reqIndex[parent], name)
					if i >= 0 {
						reqIndex[parent] = append(reqIndex[parent][:i], reqIndex[parent][i+1:]...)
					}
				}
				shouldContinue = true
			}
		}
	}

	return result
}

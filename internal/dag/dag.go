package dag

import "slices"

type Input struct {
	Steps map[string]Step
}

type Plan struct {
	Steps []string
}

type Step struct {
	Requires   []string
	RequiredBy []string
}

func Resolve(input Input) Plan {
	result := Plan{}

	index := make(map[string][]string)

	for name, step := range input.Steps {
		index[name] = append(index[name], step.Requires...)
		for _, parent := range step.RequiredBy {
			index[parent] = append(index[parent], name)
		}
	}

	shouldContinue := true
	for shouldContinue {
		shouldContinue = false
		for name, requires := range index {
			if len(requires) == 0 {
				result.Steps = append(result.Steps, name)
				delete(index, name)
				for parent := range index {
					i := slices.Index(index[parent], name)
					if i >= 0 {
						index[parent] = append(index[parent][:i], index[parent][i+1:]...)
					}
				}
				shouldContinue = true
			}
		}
	}

	return result
}

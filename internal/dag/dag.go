package dag

import "slices"

type Input struct {
	Tasks map[string]Task
}

type Output struct {
	Tasks []string
}

type Task struct {
	Requires   []string
	RequiredBy []string
}

func Resolve(input Input) Output {
	result := Output{}

	index := make(map[string][]string)

	for name, task := range input.Tasks {
		index[name] = append(index[name], task.Requires...)
		for _, parent := range task.RequiredBy {
			index[parent] = append(index[parent], name)
		}
	}

	shouldContinue := true
	for shouldContinue {
		shouldContinue = false
		for name, requires := range index {
			if len(requires) == 0 {
				result.Tasks = append(result.Tasks, name)
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

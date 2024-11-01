package indexer

import (
	"strings"

	"github.com/sirikon/ebro/internal/config"
)

type Index map[string]config.Task

func MakeIndex(module *config.Module) Index {
	index := indexModule(module, []string{})
	for _, task := range index {
		NormalizeTaskReferences(index, task.Requires)
		NormalizeTaskReferences(index, task.RequiredBy)
	}
	return index
}

func NormalizeTaskReferences(index Index, task_names []string) {
	for i, task_name := range task_names {
		defaulted_task_name := task_name + ":default"
		_, taskExists := index[task_name]
		_, defaultedTaskExists := index[defaulted_task_name]
		if !taskExists && defaultedTaskExists {
			task_names[i] = defaulted_task_name
		}
	}
}

func indexModule(module *config.Module, trail []string) Index {
	result := make(Index)
	prefix := ":" + strings.Join(append(trail, ""), ":")

	for task_name, task := range module.Tasks {
		for i := range task.Requires {
			task.Requires[i] = prefix + task.Requires[i]
		}
		for i := range task.RequiredBy {
			task.RequiredBy[i] = prefix + task.RequiredBy[i]
		}
		result[prefix+task_name] = task
	}

	for module_name, module := range module.Modules {
		module_tasks := indexModule(&module, append(trail, module_name))
		for task_name, task := range module_tasks {
			result[task_name] = task
		}
	}

	return result
}

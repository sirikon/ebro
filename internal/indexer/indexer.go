package indexer

import (
	"strings"

	"github.com/sirikon/ebro/internal/config"
)

type Index map[string]config.Task

func MakeIndex(module *config.Module) Index {
	index := indexModule(module, []string{})

	for _, task := range index {
		for i := range task.Requires {
			_, taskExists := index[task.Requires[i]]
			_, defaultTaskExists := index[task.Requires[i]+":default"]
			if !taskExists && defaultTaskExists {
				task.Requires[i] = task.Requires[i] + ":default"
			}
		}
		for i := range task.RequiredBy {
			_, taskExists := index[task.RequiredBy[i]]
			_, defaultTaskExists := index[task.RequiredBy[i]+":default"]
			if !taskExists && defaultTaskExists {
				task.RequiredBy[i] = task.RequiredBy[i] + ":default"
			}
		}
	}

	return index
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

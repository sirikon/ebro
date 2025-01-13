package config

import (
	"iter"
	"maps"
	"slices"
)

func (m *Module) TasksSorted() iter.Seq2[string, *Task] {
	taskNames := slices.Sorted(maps.Keys(m.Tasks))
	return func(yield func(string, *Task) bool) {
		for _, taskName := range taskNames {
			if !yield(taskName, m.Tasks[taskName]) {
				return
			}
		}
	}
}

func (m *Module) ModulesSorted() iter.Seq2[string, *Module] {
	moduleNames := slices.Sorted(maps.Keys(m.Modules))
	return func(yield func(string, *Module) bool) {
		for _, moduleName := range moduleNames {
			if !yield(moduleName, m.Modules[moduleName]) {
				return
			}
		}
	}
}

func (m *Module) ImportsSorted() iter.Seq2[string, *Import] {
	importNames := slices.Sorted(maps.Keys(m.Imports))
	return func(yield func(string, *Import) bool) {
		for _, importName := range importNames {
			if !yield(importName, m.Imports[importName]) {
				return
			}
		}
	}
}

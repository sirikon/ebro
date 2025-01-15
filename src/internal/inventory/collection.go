package inventory

import (
	"iter"
	"maps"
	"slices"

	"github.com/sirikon/ebro/internal/core"
)

func (inv Inventory) TasksSorted() iter.Seq2[core.TaskId, *core.Task] {
	taskNames := slices.Sorted(maps.Keys(inv.Tasks))
	return func(yield func(core.TaskId, *core.Task) bool) {
		for _, taskName := range taskNames {
			if !yield(taskName, inv.Tasks[taskName]) {
				return
			}
		}
	}
}

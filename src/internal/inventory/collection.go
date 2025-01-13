package inventory

import (
	"iter"
	"maps"
	"slices"

	"github.com/sirikon/ebro/internal/config"
)

func (inv Inventory) TasksSorted() iter.Seq2[string, *config.Task] {
	taskNames := slices.Sorted(maps.Keys(inv.Tasks))
	return func(yield func(string, *config.Task) bool) {
		for _, taskName := range taskNames {
			if !yield(taskName, inv.Tasks[taskName]) {
				return
			}
		}
	}
}

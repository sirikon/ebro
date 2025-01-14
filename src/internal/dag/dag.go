package dag

import (
	"maps"
	"slices"

	"github.com/sirikon/ebro/internal/utils"
)

type Dag[T ~string] struct {
	dependencyMap map[T]map[T]bool
}

func NewDag[T ~string]() Dag[T] {
	return Dag[T]{dependencyMap: map[T]map[T]bool{}}
}

func (d *Dag[T]) Link(parentNode T, childNodes ...T) {
	dependencies, ok := d.dependencyMap[parentNode]
	if !ok {
		dependencies = make(map[T]bool)
		d.dependencyMap[parentNode] = dependencies
	}
	for _, childNode := range childNodes {
		dependencies[childNode] = true
	}
}

func (d *Dag[T]) Resolve(targets []T) ([]T, map[T][]T) {
	result := []T{}
	nodesToSolve := utils.NewSet[T]()

	for _, node := range targets {
		nodesToSolve.Add(node)
	}

	for i := 0; i < nodesToSolve.Length(); i++ {
		node := nodesToSolve.Get(i)
		nodesToSolve.Add(slices.Collect(maps.Keys(d.dependencyMap[node]))...)
	}

	shouldContinue := true
	for shouldContinue {
		shouldContinue = false
		batch := []T{}

		for _, node := range nodesToSolve.List() {
			dependencies := d.dependencyMap[node]
			if len(dependencies) == 0 {
				batch = append(batch, node)
				shouldContinue = true
			}
		}

		slices.Sort(batch)

		for _, node := range batch {
			nodesToSolve.Delete(node)
			delete(d.dependencyMap, node)
			for parentNode := range d.dependencyMap {
				delete(d.dependencyMap[parentNode], node)
			}
		}

		result = append(result, batch...)
	}

	if nodesToSolve.Length() > 0 {
		remains := make(map[T][]T)
		for _, parentNode := range nodesToSolve.List() {
			nodes := slices.Collect(maps.Keys(d.dependencyMap[parentNode]))
			slices.Sort(nodes)
			remains[parentNode] = nodes
		}
		return result, remains
	}

	return result, nil
}

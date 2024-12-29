package dag

import (
	"maps"
	"slices"

	"github.com/sirikon/ebro/internal/utils"
)

type Dag struct {
	dependencyMap map[string]map[string]bool
}

func NewDag() Dag {
	return Dag{dependencyMap: map[string]map[string]bool{}}
}

func (d *Dag) Link(parentNode string, childNodes ...string) {
	dependencies, ok := d.dependencyMap[parentNode]
	if !ok {
		dependencies = make(map[string]bool)
		d.dependencyMap[parentNode] = dependencies
	}
	for _, childNode := range childNodes {
		dependencies[childNode] = true
	}
}

func (d *Dag) Resolve(targets []string) ([]string, map[string][]string) {
	result := []string{}
	nodesToSolve := utils.NewSet[string]()

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
		batch := []string{}

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
		remains := make(map[string][]string)
		for _, parentNode := range nodesToSolve.List() {
			nodes := slices.Collect(maps.Keys(d.dependencyMap[parentNode]))
			slices.Sort(nodes)
			remains[parentNode] = nodes
		}
		return result, remains
	}

	return result, nil
}

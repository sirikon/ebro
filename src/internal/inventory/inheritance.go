package inventory

import (
	"fmt"
	"maps"
	"slices"

	"github.com/sirikon/ebro/internal/config"
	"github.com/sirikon/ebro/internal/dag"
	"github.com/sirikon/ebro/internal/utils"

	"github.com/goccy/go-yaml"
)

func resolveInheritanceOrder(inv Inventory, unrunnableTasks map[string]error) ([]string, error) {
	inheritanceDag := dag.NewDag()

	checkTaskCanExtend := func(taskReferenceString string) (bool, error) {
		ref, _ := config.ParseTaskReference(taskReferenceString)
		if err := unrunnableTasks[ref.PathString()]; err != nil {
			if !ref.IsOptional {
				return false, fmt.Errorf("task %v can't extend: %w", taskReferenceString, err)
			}
			return false, nil
		}
		return true, nil
	}

	for taskName, task := range inv.Tasks {
		canRun, err := checkTaskCanExtend(taskName)
		if err != nil {
			return nil, err
		}
		if !canRun {
			continue
		}

		for _, taskReferenceString := range task.Extends {
			canRun, err := checkTaskCanExtend(taskReferenceString)
			if err != nil {
				return nil, err
			}
			if !canRun {
				continue
			}
			ref, _ := config.ParseTaskReference(taskReferenceString)
			inheritanceDag.Link(taskName, ref.PathString())
		}
	}

	result, remains := inheritanceDag.Resolve(slices.Collect(maps.Keys(inv.Tasks)))
	if remains != nil {
		remainsData, err := yaml.Marshal(remains)
		if err != nil {
			return nil, fmt.Errorf("inheritance order resolution could not complete. error while turning dependency index to yaml: %w", err)
		}
		return nil, fmt.Errorf("inheritance order resolution could not complete. "+
			"there could be a cyclic dependency. "+
			"here is the list of tasks and their inheritance data:\n%s", string(remainsData))
	}

	return result, nil
}

func applyInheritance(childTask *config.Task, parentTask *config.Task) {
	childTask.Extends = nil
	childTask.Requires = utils.Dedupe(slices.Concat(childTask.Requires, parentTask.Requires))
	childTask.RequiredBy = utils.Dedupe(slices.Concat(childTask.RequiredBy, parentTask.RequiredBy))
	if childTask.Script == "" {
		childTask.Script = parentTask.Script
	}
	if childTask.Quiet == nil {
		childTask.Quiet = parentTask.Quiet
	}
	if parentTask.When != nil {
		if childTask.When == nil {
			when := config.When{}
			childTask.When = &when
		}
		if childTask.When.CheckFails == "" {
			childTask.When.CheckFails = parentTask.When.CheckFails
		}
		if childTask.When.OutputChanges == "" {
			childTask.When.OutputChanges = parentTask.When.OutputChanges
		}
	}
}

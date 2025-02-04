package core

import (
	"fmt"
	"regexp"
	"strings"
)

type TaskReference struct {
	Path       []string
	IsRelative bool
	IsOptional bool
}

var taskReferenceRegex = regexp.MustCompile(`^:?[a-zA-Z0-9-_\.]+(:[a-zA-Z0-9-_\.]+)*\??$`)

func ValidateTaskReference(text string) error {
	if !taskReferenceRegex.MatchString(text) {
		return fmt.Errorf("task reference is invalid")
	}
	return nil
}

func MustParseTaskReference(text string) TaskReference {
	result := TaskReference{
		Path:       []string{},
		IsRelative: true,
		IsOptional: false,
	}

	if err := ValidateTaskReference(text); err != nil {
		panic(err)
	}

	if strings.HasPrefix(text, ":") {
		text = strings.TrimPrefix(text, ":")
		result.IsRelative = false
	}

	if strings.HasSuffix(text, "?") {
		text = strings.TrimSuffix(text, "?")
		result.IsOptional = true
	}

	result.Path = strings.Split(text, ":")

	return result
}

func (tp TaskReference) Absolute(parentPath []string) TaskReference {
	if !tp.IsRelative {
		return tp
	}

	return TaskReference{
		Path:       append(parentPath, tp.Path...),
		IsRelative: false,
		IsOptional: tp.IsOptional,
	}
}

func (tp TaskReference) Concat(extraPath ...string) TaskReference {
	return TaskReference{
		Path:       append(tp.Path, extraPath...),
		IsRelative: tp.IsRelative,
		IsOptional: tp.IsOptional,
	}
}

func (tp TaskReference) TaskId() TaskId {
	if tp.IsRelative {
		panic("cannot build TaskId from relative TaskReference")
	}
	return NewTaskId(tp.Path[:len(tp.Path)-1], tp.Path[len(tp.Path)-1])
}

func ResolveReferences(inventory *Inventory, task *Task, taskReferences []string) ([]TaskId, error) {
	result := []TaskId{}
	for _, taskReference := range taskReferences {
		if err := ValidateTaskReference(taskReference); err != nil {
			return nil, fmt.Errorf("validating '%v': %w", taskReference, err)
		}

		ref := MustParseTaskReference(taskReference)
		if ref.IsRelative {
			ref = ref.Absolute(task.Id.ModulePath())
		}

		referencedTaskId, _ := inventory.FindTask(ref)
		if referencedTaskId == nil {
			if ref.IsOptional {
				continue
			} else {
				return nil, fmt.Errorf("referenced task '%v' does not exist", ref.TaskId())
			}
		}

		result = append(result, *referencedTaskId)
	}
	return result, nil
}

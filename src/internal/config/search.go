package config

import (
	"fmt"
	"regexp"
	"strings"
)

type TaskId struct {
	ModuleTrail []string
	TaskName    string
}

func (tid TaskId) String() string {
	chunks := []string{""}
	chunks = append(chunks, tid.ModuleTrail...)
	chunks = append(chunks, tid.TaskName)
	return strings.Join(chunks, ":")
}

type TaskReference struct {
	Path       []string
	IsRelative bool
	IsOptional bool
}

func ParseTaskReference(text string) (TaskReference, error) {
	result := TaskReference{
		Path:       []string{},
		IsRelative: true,
		IsOptional: false,
	}

	re, err := regexp.Compile(`^:?[a-zA-Z0-9-_\.]+(:[a-zA-Z0-9-_\.]+)*\??$`)
	if !re.MatchString(text) {
		return result, fmt.Errorf("task reference is invalid")
	}

	if err != nil {
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

	return result, nil
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

func (tp TaskReference) PathString() string {
	chunks := []string{}
	if !tp.IsRelative {
		chunks = append(chunks, ":")
	}
	chunks = append(chunks, strings.Join(tp.Path, ":"))
	return strings.Join(chunks, "")
}

func (tp TaskReference) String() string {
	chunks := []string{}
	chunks = append(chunks, tp.PathString())
	if tp.IsOptional {
		chunks = append(chunks, "?")
	}
	return strings.Join(chunks, "")
}

func (m *Module) GetTask(taskReference TaskReference) (*TaskId, *Task) {
	if taskReference.IsRelative {
		panic("cannot call getTask with a relative taskPath")
	}

	moduleTrail := []string{}
	currentModule := m
	for i, part := range taskReference.Path {
		if i >= (len(taskReference.Path) - 1) {
			module, ok := currentModule.Modules[part]
			if ok {
				if task, ok := module.Tasks["default"]; ok {
					return &TaskId{
						ModuleTrail: append(moduleTrail, part),
						TaskName:    "default",
					}, task
				}
			} else {
				if task, ok := currentModule.Tasks[part]; ok {
					return &TaskId{
						ModuleTrail: moduleTrail,
						TaskName:    part,
					}, task
				}
			}
		} else {
			module, ok := currentModule.Modules[part]
			if ok {
				moduleTrail = append(moduleTrail, part)
				currentModule = module
			} else {
				return nil, nil
			}
		}
	}

	return nil, nil
}

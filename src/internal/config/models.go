package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirikon/ebro/internal/core"
)

/* ==== Models reflected in Ebro.yaml files ==== */

type Module = core.ModuleBase[Task, Import]
type Task = core.TaskBase[string, When]
type Import = core.Import
type When = core.When

/* ============================================= */

var taskReferenceRegex = regexp.MustCompile(`^:?[a-zA-Z0-9-_\.]+(:[a-zA-Z0-9-_\.]+)*\??$`)

type TaskReference struct {
	Path       []string
	IsRelative bool
	IsOptional bool
}

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

func (tp TaskReference) TaskId() core.TaskId {
	if tp.IsRelative {
		panic("cannot build TaskId from relative TaskReference")
	}
	return core.MakeTaskId(tp.Path[:len(tp.Path)-1], tp.Path[len(tp.Path)-1])
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

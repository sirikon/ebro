package config

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

func MakeTaskReference(path []string) TaskReference {
	return TaskReference{
		Path:       path,
		IsRelative: false,
		IsOptional: false,
	}
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
	return strings.Join(tp.Path, ":")
}

func (tp TaskReference) String() string {
	chunks := []string{}
	if !tp.IsRelative {
		chunks = append(chunks, ":")
	}
	chunks = append(chunks, tp.PathString())
	if tp.IsOptional {
		chunks = append(chunks, "?")
	}
	return strings.Join(chunks, "")
}

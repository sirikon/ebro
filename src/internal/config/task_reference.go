package config

import (
	"fmt"
	"regexp"
	"strings"
)

type TaskReference struct {
	Parts      []string
	IsRelative bool
	IsOptional bool
}

func ParseTaskReference(text string) (TaskReference, error) {
	result := TaskReference{
		Parts:      []string{},
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

	result.Parts = strings.Split(text, ":")

	return result, nil
}

func (tp TaskReference) Absolute(parts []string) TaskReference {
	if !tp.IsRelative {
		return tp
	}

	return TaskReference{
		Parts:      append(parts, tp.Parts...),
		IsRelative: false,
		IsOptional: tp.IsOptional,
	}
}

func (tp TaskReference) PartsString() string {
	return strings.Join(tp.Parts, ":")
}

func (tp TaskReference) String() string {
	chunks := []string{}
	if !tp.IsRelative {
		chunks = append(chunks, ":")
	}
	chunks = append(chunks, tp.PartsString())
	if tp.IsOptional {
		chunks = append(chunks, "?")
	}
	return strings.Join(chunks, "")
}

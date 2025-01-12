package config

import (
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

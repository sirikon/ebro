package core2

import (
	"regexp"
	"strings"
)

type TaskId string

var TaskIdRe = regexp.MustCompile("^(:" + nameValidCharsRe + "+)+$")

func NewTaskId(modulePath []string, taskName string) TaskId {
	chunks := []string{""}
	chunks = append(chunks, modulePath...)
	chunks = append(chunks, taskName)
	result := TaskId(strings.Join(chunks, ":"))
	result.MustBeValid()
	return result
}

func (tid TaskId) MustBeValid() {
	if !TaskIdRe.MatchString(string(tid)) {
		panic("TaskId ismalformed: " + string(tid))
	}
}

func (tid TaskId) ModulePath() []string {
	parts := tid.parts()
	return parts[:len(parts)-1]
}

func (tid TaskId) TaskName() string {
	parts := tid.parts()
	return parts[len(parts)-1]
}

func (tid TaskId) parts() []string {
	tid.MustBeValid()
	return strings.Split(strings.TrimPrefix(string(tid), ":"), ":")
}

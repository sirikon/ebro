package core2

import (
	"fmt"
	"regexp"
	"strings"
)

type TaskId string

var nameValidCharsRe = `[a-zA-Z0-9-_\.]`
var NameRe = regexp.MustCompile("^" + nameValidCharsRe + "+$")
var TaskIdRe = regexp.MustCompile("^(:" + nameValidCharsRe + "+)+$")

func ValidateName(name string) error {
	if !NameRe.MatchString(name) {
		return fmt.Errorf("name does not match the following regex: %v", NameRe.String())
	}
	return nil
}

func NewTaskId(moduleTrail []string, taskName string) TaskId {
	chunks := []string{""}
	chunks = append(chunks, moduleTrail...)
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

func (tid TaskId) ModuleTrail() []string {
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

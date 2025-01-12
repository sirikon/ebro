package config

import "fmt"

type TaskNotFoundError struct {
	TaskReference TaskReference
}

func (e TaskNotFoundError) Error() string {
	return fmt.Sprintf("task %v does not exist", e.TaskReference.PathString())
}

func NewTaskNotFoundError(taskReference TaskReference) TaskNotFoundError {
	return TaskNotFoundError{TaskReference: taskReference}
}

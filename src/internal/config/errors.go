package config

import (
	"errors"
	"fmt"
)

var ErrTaskNotFound = errors.New("task not found")

func WrapErrTaskNotFound(taskReference TaskReference) error {
	return fmt.Errorf("task %v: %w", taskReference.String(), ErrTaskNotFound)
}

package cataloger

import "fmt"

func (c Catalog) Validate() error {
	for taskName, task := range c {
		for _, otherTaskName := range append(task.Requires, task.RequiredBy...) {
			if _, ok := c[otherTaskName]; !ok {
				return fmt.Errorf("task %v: referenced task %v does not exist", taskName, otherTaskName)
			}
		}
	}
	return nil
}

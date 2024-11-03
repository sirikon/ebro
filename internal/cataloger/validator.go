package cataloger

import "fmt"

func (c Catalog) Validate() error {
	for task_name, task := range c {
		for _, other_task_name := range append(task.Requires, task.RequiredBy...) {
			if _, ok := c[other_task_name]; !ok {
				return fmt.Errorf("task %v: referenced task %v does not exist", task_name, other_task_name)
			}
		}
	}
	return nil
}

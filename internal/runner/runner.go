package runner

import (
	"fmt"

	"github.com/sirikon/ebro/internal/indexer"
	"github.com/sirikon/ebro/internal/planner"
)

func Run(index indexer.Index, plan planner.Plan) {
	for _, task_name := range plan {
		fmt.Println(task_name)
	}
}

package agent

import (
	"DistributedCalculator/internal/worker"
	"DistributedCalculator/pkg/database"
	"DistributedCalculator/pkg/env"
	"DistributedCalculator/pkg/expression"
	"errors"

	"go.uber.org/zap"
)

type Agent struct {
	Tasks  []chan expression.Expression
	IsFree []bool
	Work   []expression.Expression
}


func NewAgent(logger *zap.Logger) *Agent {
	tasks := make([]chan expression.Expression, 0)
	IsFree := make([]bool, 0)
	Work := make([]expression.Expression, 0)
	for i := 0; i < env.Workers; i++ {
		tasks = append(tasks, make(chan expression.Expression))

		worker.StartWorker(tasks[i], logger)
		IsFree = append(IsFree, true)
		Work = append(Work, expression.Expression{})

	}
	return &Agent{Tasks: tasks, IsFree: IsFree, Work: Work}
}

func (ag *Agent) AddTask(expr expression.Expression) error {
	for ind, i := range ag.IsFree {
		if i {
			ag.Tasks[ind] <- expr
			ag.Work[ind] = expr
			expr.Status = 1
			database.UpdateExpr(expr)
			ag.IsFree[ind] = false
			go func() {
				newexp := <-ag.Tasks[ind]
				newexp.Status = 0
				database.UpdateExpr(newexp)
				ag.IsFree[ind] = true
				ag.Work[ind] = expression.Expression{}
			}()
			return nil
		}
	}
	return errors.New("")
}

package worker

import (
	"DistributedCalculator/pkg/ExpParser"
	"DistributedCalculator/pkg/database"
	"DistributedCalculator/pkg/expression"
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)


func StartWorker(exprchan chan expression.Expression, logger *zap.Logger) {
	go func() {
		for {
			expr := <-exprchan
			node := expr.Node
			if ExpParser.Length(&node) > 20 {
				newexpr := expression.NewExpression("added")
				newexpr.Node = *node.Right
				b, _ := json.Marshal(newexpr)
				r := bytes.NewReader(b)
				http.Post("http://localhost:8081", "application/json", r)
				leftres := calcNode(node.Left)
				rightres := 0.0
				for {
					g := database.ReadExp(newexpr.Id, logger)
					if g != nil {
						rightres = g.Result
						break
					}
				}

				expr.Result = ExpParser.PerformOperation(node.Operator, leftres, rightres)
			} else {
				expr.Result = calcNode(&expr.Node)
			}
			exprchan <- expr
		}
	}()
}

func calcNode(node *ExpParser.Node) float64 {
	if node.Operator == "" {

		return node.Value
	} else {
		if node.Left == nil || node.Right == nil {
		} else {
			return ExpParser.PerformOperation(node.Operator, calcNode(node.Left), calcNode(node.Right))
		}
	}
	return 0
}

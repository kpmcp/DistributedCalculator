package expression

import (
	"DistributedCalculator/pkg/ExpParser"
	"fmt"

	"github.com/google/uuid"
)

type Expression struct {
	Name   string // Изначальное значение выражения
	Status int    // Статус выражения
	Id     int
	Result float64     // результат выражения
	Node   ExpParser.Node
}

func NewExpression(Name string) *Expression {
	return &Expression{Name: Name, Status: 2, Id: int(uuid.New().ID())}
}

func (exp *Expression) ForTemplate() string {
	var stat string
	if exp.Status == 0 {
		stat = "Результат:"
	} else if exp.Status == 1 {
		stat = "Считается"
	} else if exp.Status == 2 {
		stat = "Ожидает рассчёта"
	} else if exp.Status == 3 {
		stat = "Выражение неккоректно"
	}
	if exp.Status == 0 {
		return fmt.Sprintf("id: %d, %s %s %.4f", exp.Id, exp.Name, stat, exp.Result)
	} else {
		return fmt.Sprintf("id: %d, %s %s", exp.Id, exp.Name, stat)
	}

}

package ExpParser

import (
	constants "DistributedCalculator/pkg/env"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Node struct {
	Left     *Node
	Right    *Node
	Operator string
	Value    float64
}

func NewNode() *Node {
	return &Node{nil, nil, "", 0}
}

func tokenize(str string) []string {
	re := regexp.MustCompile(`\d+|\D`)
	tokens := re.FindAllString(str, -1)
	return tokens
}


func ParseExpr(s string) (*Node, error) {
	var (
		tokens    = tokenize(s)
		stack     []*Node
		operators []string
	)
	var err error

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			for len(operators) > 0 && prec(operators[len(operators)-1]) >= prec(token) {
				err = popOperator(&stack, &operators)
				if err != nil {
					return nil, err
				}
			}
			operators = append(operators, token)
		default:
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, err
			}
			stack = append(stack, &Node{Value: value})
		}
	}
	for len(operators) > 0 {
		popOperator(&stack, &operators)
	}

	if len(stack) != 1 {
		return nil, errors.New("err")
	}

	return stack[0], nil
}

func prec(s string) int {
	if s == "^" {
		return 3
	} else if (s == "/") || (s == "*") {
		return 2
	} else if (s == "+") || (s == "-") {
		return 1
	} else {
		return -1
	}
}


func popOperator(stack *[]*Node, operators *[]string) error {
	operator := (*operators)[len(*operators)-1]
	*operators = (*operators)[:len(*operators)-1]
	if len(*stack) == 0 {
		return errors.New("ss")
	}
	right := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	if len(*stack) == 0 {
		return errors.New("ф")
	}
	left := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]

	node := &Node{Right: right, Left: left, Operator: operator}
	*stack = append(*stack, node)
	return nil
}


func EvaluatePostOrder(node *Node, subExpressions *map[int]string, counter *int) error {
	if node == nil {
		return nil
	}

	if node.Left != nil {
		err := EvaluatePostOrder(node.Left, subExpressions, counter)
		if err != nil {
			return err
		}
	}

	if node.Right != nil {
		err := EvaluatePostOrder(node.Right, subExpressions, counter)
		if err != nil {
			return err
		}
	}

	if node.Left == nil && node.Right == nil {
		(*subExpressions)[*counter] = fmt.Sprintf("%.2f", node.Value)
		*counter++
	}

	if node.Operator != "" {
		lastIndex := *counter - 1
		secondLastIndex := lastIndex - 1
		subExpression := fmt.Sprintf("%s %s %s", (*subExpressions)[secondLastIndex], node.Operator, (*subExpressions)[lastIndex])
		(*subExpressions)[*counter] = subExpression
		*counter++
	}
	return nil
}

func ValidatedPostOrder(s string) (map[int]string, error) {
	node, err := ParseExpr(s)
	if err != nil {
		return nil, err
	}
	subExps := make(map[int]string)
	var counter int
	err = EvaluatePostOrder(node, &subExps, &counter)
	if err != nil {
		return nil, err
	}
	for key, val := range subExps {
		if len(val) == 4 {
			delete(subExps, key)
		}
	}
	return subExps, nil
}


func CalcNode(node *Node) float64 {
	if node.Operator == "" {
		return node.Value
	} else {
		if node.Left == nil || node.Right == nil {
		} else {
			return PerformOperation(node.Operator, CalcNode(node.Left), CalcNode(node.Right))
		}
	}
	return 0
}


func PerformOperation(operator string, operand1, operand2 float64) float64 {
	switch operator {
		case "+":
			time.Sleep(time.Duration(constants.Plus) * time.Second)
			return operand1 + operand2
		case "-":
			time.Sleep(time.Duration(constants.Minus) * time.Second)
			return operand1 - operand2
		case "*":
			time.Sleep(time.Duration(constants.Mul) * time.Second)
			return operand1 * operand2
		case "/":
			time.Sleep(time.Duration(constants.Div) * time.Second)
			return operand1 / operand2
		default:
			panic(errors.New("not an operator"))
	}
}

func Length(node *Node) int {
	if node.Left != nil && node.Right != nil {
		return Length(node.Left) + Length(node.Right) + 1
	}
	return 1
}

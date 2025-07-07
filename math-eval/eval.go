package main

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

type Expression struct {
	Expression string
	Parameters map[string]interface{}
}

var expressions = []Expression{
	{"a + b", map[string]interface{}{"a": 1, "b": 2}},
	{"a - b", map[string]interface{}{"a": 1, "b": 2}},
	{"a * b", map[string]interface{}{"a": 1, "b": 2}},
	{"a / b", map[string]interface{}{"a": 1, "b": 2}},
	{"a % b", map[string]interface{}{"a": 1, "b": 2}},
	{"a ** b", map[string]interface{}{"a": 1, "b": 2}},
	{"a == b", map[string]interface{}{"a": 1, "b": 2}},
	{"a != b", map[string]interface{}{"a": 1, "b": 2}},
	{"a > b", map[string]interface{}{"a": 1, "b": 2}},
	{"a < b", map[string]interface{}{"a": 1, "b": 2}},
	{"a >= b", map[string]interface{}{"a": 1, "b": 2}},
	{"a <= b", map[string]interface{}{"a": 1, "b": 2}},
	{"a && b", map[string]interface{}{"a": true, "b": true}},
	{"a || b", map[string]interface{}{"a": true, "b": true}},
	{"!a", map[string]interface{}{"a": true}},
	{"area * 320", map[string]interface{}{"area": 100}},
	{"area * gas_price", map[string]interface{}{"area": 100, "gas_price": 10}},
	{"usage / sum_usage * price", map[string]interface{}{"usage": 100, "sum_usage": 1200, "price": 30000}},
}

func main() {
	for _, expr := range expressions {
		result, err := evalExpression(expr.Expression, expr.Parameters)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	}
}

func evalExpression(expr string, parameters map[string]interface{}) (interface{}, error) {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return nil, err
	}
	result, err := expression.Evaluate(parameters)
	if err != nil {
		return nil, err
	}
	return result, nil
}

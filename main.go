package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Operator rune

func (op Operator) IsOperator() bool {
	switch op {
	case '+', '-', '*', '/':
		return true
	default:
		return false
	}
}

func Calc(expression string) (float64, error) {
	var numbers []float64
	var operator Operator
	var num string

	for _, symb := range expression {
		if Operator(symb).IsOperator() {
			operator = Operator(symb)
			f, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, errors.New("invalid number")
			}
			numbers = append(numbers, f)
			num = ""
		} else {
			num += string(symb)
		}
	}

	if num != "" {
		f, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return 0, errors.New("invalid number")
		}
		numbers = append(numbers, f)
	}

	if len(numbers) != 2 {
		return 0, errors.New("invalid expression")
	}

	switch operator {
	case '+':
		return numbers[0] + numbers[1], nil
	case '-':
		return numbers[0] - numbers[1], nil
	case '*':
		return numbers[0] * numbers[1], nil
	case '/':
		if numbers[1] == 0 {
		return 0, errors.New("division by zero")
	}
		return numbers[0] / numbers[1], nil
	default:
		return 0, errors.New("unknown operator")
	}
}

func main() {
	fmt.Println("Hello")
}

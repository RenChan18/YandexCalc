package main

import (
	"errors"
	"strconv"
	"unicode"
)

type Operator rune

var operatorPrecedence = map[Operator]int{
	'(': 0,
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
	'^': 3,
}

func (op Operator) Precedence() int {
	return operatorPrecedence[op]
}
func (op Operator) getOp() rune {
	return rune(op)
}
func IsOperator(r rune) bool {
	_, exists := operatorPrecedence[Operator(r)]
	return exists
}

func contains(slice []Operator, value Operator) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func pop(s []Operator) ([]Operator, Operator) {
	lastElement := s[len(s)-1]
	s = s[:len(s)-1]
	return s, lastElement
}

func pops(s []float64) ([]float64, float64) {
	lastElement := s[len(s)-1]
	s = s[:len(s)-1]
	return s, lastElement
}

func compute(a, b float64, op Operator) (float64, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	default:
		return 0, errors.New("unknown operator")
	}
}

func toPostfix(expression string) string {
	var postfixEpr string
	var sNum string
	var operator []Operator
	var last Operator
	for _, symb := range expression {
		if unicode.IsDigit(symb) {
			sNum += string(symb)
		} else if IsOperator(symb) {
			if sNum != "" {
				postfixEpr += sNum + " "
				sNum = ""
			}

			for len(operator) > 0 && operator[len(operator)-1].Precedence() >= Operator(symb).Precedence() {
				operator, last = pop(operator)
				postfixEpr += string(last.getOp()) + " "
			}
			operator = append(operator, Operator(symb))
		} else if symb == '(' {
			operator = append(operator, Operator(symb))
		} else if symb == ')' {
			for len(operator) > 0 && operator[len(operator)-1] != '(' {
				operator, last = pop(operator)
				postfixEpr += string(last.getOp()) + " "
			}
			if len(operator) > 0 {
				operator, _ = pop(operator)
			}
		}
	}

	if sNum != "" {
		postfixEpr += sNum + " "
	}

	for len(operator) > 0 {
		operator, last = pop(operator)
		postfixEpr += string(last.getOp()) + " "
	}

	return postfixEpr
}
func Calc(expression string) (float64, error) {
	var numbers []float64
	var sNum string
	postfix := toPostfix(expression)

	for _, token := range postfix {
		if unicode.IsDigit(token) {
			sNum += string(token)
		} else if IsOperator(token) {
			if sNum != "" {
				f, err := strconv.ParseFloat(string(token), 64)
				if err != nil {

				} else {
					numbers = numbers.append(f)
					sNum = ""
				}
			}
			var second, first float64
			numbers, second = pops(numbers)
			numbers, first = pops(numbers)

		}
	}

	return result, nil
}

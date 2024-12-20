package calculator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)


func Calc(expression string) (float64, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}

	return parseExpression(tokens)
}

func tokenize(expression string) ([]string, error) {
	var tokens []string
	var number strings.Builder

	for _, ch := range expression {
		switch {
		case ch >= '0' && ch <= '9' || ch == '.':
			number.WriteRune(ch)
		case ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '(' || ch == ')':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			tokens = append(tokens, string(ch))
		case ch == ' ':
			continue
		default:
			return nil, fmt.Errorf("invalid character: %c", ch)
		}
	}

	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	return tokens, nil
}

func parseExpression(tokens []string) (float64, error) {
	var stack []float64
	var operators []string

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if isNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				err := applyOperator(&stack, &operators)
				if err != nil {
					return 0, err
				}
			}
			if len(operators) == 0 {
				return 0, errors.New("mismatched parentheses")
			}
			operators = operators[:len(operators)-1] // Удаляем "("
		} else if isOperator(token) {
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[token] {
				err := applyOperator(&stack, &operators)
				if err != nil {
					return 0, err
				}
			}
			operators = append(operators, token)
		} else {
			return 0, fmt.Errorf("invalid token: %s", token)
		}
	}

	for len(operators) > 0 {
		err := applyOperator(&stack, &operators)
		if err != nil {
			return 0, err
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}

func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func applyOperator(stack *[]float64, operators *[]string) error {
	if len(*stack) < 2 {
		return errors.New("not enough operands")
	}

	b := (*stack)[len(*stack)-1]
	a := (*stack)[len(*stack)-2]
	op := (*operators)[len(*operators)-1]

	*stack = (*stack)[:len(*stack)-2]
	*operators = (*operators)[:len(*operators)-1]

	var result float64
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return errors.New("division by zero")
		}
		result = a / b
	default:
		return fmt.Errorf("unknown operator: %s", op)
	}

	*stack = append(*stack, result)
	return nil
}

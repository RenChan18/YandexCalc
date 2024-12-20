package main

import "testing"

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{"simple", "1+1", 2},
		{"priority", "(2+2)*2", 8},
		{"mixed", "2+2*2", 6},
		{"division", "1/2", 0.5},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error: %v", testCase.expression, err)
			}
			if val != testCase.expectedResult {
				t.Fatalf("expected %f, got %f", testCase.expectedResult, val)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr string
	}{
		{"invalid operator", "1+1*", "not enough operands"},
		{"mismatched parenthesis", "(1+2", "mismatched parentheses"},
		{"division by zero", "1/0", "division by zero"},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Calc(testCase.expression)
			if err == nil || err.Error() != testCase.expectedErr {
				t.Fatalf("expected error %s, got %v", testCase.expectedErr, err)
			}
		})
	}
}


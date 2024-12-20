package tests

import (
	"calc_service/internal/calculator"
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		hasError   bool
	}{
		{"3+5", 8, false},
		{"10 + 2 * 6", 22, false},
		{"100 * 2 + 12", 212, false},
		{"100 * ( 2 + 12 )", 1400, false},
		{"100 * ( 2 + 12 ) / 14", 100, false},
		{"2 / 0", 0, true}, // Деление на ноль
		{"3 + a", 0, true}, // Некорректный символ
	}

	for _, tt := range tests {
		result, err := calculator.Calc(tt.expression)
		if (err != nil) != tt.hasError {
			t.Errorf("Calc(%s) error = %v, want error: %v", tt.expression, err, tt.hasError)
		}
		if !tt.hasError && result != tt.expected {
			t.Errorf("Calc(%s) = %f; want %f", tt.expression, result, tt.expected)
		}
	}
}


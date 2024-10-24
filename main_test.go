package main

import (
	"errors"
	"testing"
)

type Test struct {
	in      string
	result  float64
	exp_err error
}

var tests = []Test{
	{"2+3", 5.0, nil},
	{"2/0", 0.0, errors.New("division by zero")},
	{"2*3", 6.0, nil},
	{"3*0", 0.0, nil},
	{"3/2", 1.5, nil},
    {"1+1+1", 3, nil},
    {"2+3*0", 2, nil},
    {"20+1", 21, nil},
}

func TestCalc(t *testing.T) {
	for i, test := range tests {
		got, err := Calc(test.in)

		// Сравниваем значение результата
		if got != test.result {
			t.Errorf("#%d: got %v want %v", i, got, test.result)
		}

		// Сравниваем ошибки
		if (err != nil && test.exp_err == nil) || (err == nil && test.exp_err != nil) {
			t.Errorf("#%d: got %v want %v", i, err, test.exp_err)
		} else if err != nil && test.exp_err != nil && err.Error() != test.exp_err.Error() {
			t.Errorf("#%d: got %v want %v", i, err, test.exp_err)
		}
	}
}


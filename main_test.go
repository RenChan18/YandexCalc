package main

import "testing"

func TestCalc(t *testing.T) {
	got := ""
	want := 0, nil

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

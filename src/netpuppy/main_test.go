package main

import "testing"

func TestSum(t *testing.T) {
	got := sum(2, 3)
	want := 5

	if got != want {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

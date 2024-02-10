package main

import "testing"

func TestSum(t *testing.T) {
	var placeholder int = 69
	if placeholder != 69 {
		t.Errorf("Placeholder not 69\n")
	}
}

package utils

import "testing"

func TestInterface(t *testing.T) {
	var dummyName string = "Porschia"
	var dummyNumber int = 4444
	var dummyContact string = "Porschia: 4444"

	got := Tiddies()

	if got.name != dummyName {
		t.Errorf("WRONG")
	}

	if got.number != dummyNumber {
		t.Errorf("WTRkadsfnafn")
	}

	if got.contact() != dummyContact {
		t.Errorf("Got Contact: %v, Expected: %v\n", got.contact(), dummyContact)
	}
}

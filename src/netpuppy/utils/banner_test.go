package utils

import (
	"fmt"
	"testing"
)

func TestBanner(t *testing.T) {
	banner := Banner()

	stringCheck := fmt.Sprintf("%T", banner)
	if stringCheck != "string" {
		t.Errorf("Testing First Banner - Got: %v with type: %v, Expected: string\n", banner, stringCheck)
	}
}

func TestUserSelectionBanner(t *testing.T) {
	var testChoice string = "testServer"
	var testHost string = "0.0.0.0"
	var testRPort int = 69
	var testLPort int = 96

	banner := UserSelectionBanner(testChoice, testHost, testRPort, testLPort)

	stringCheck := fmt.Sprintf("%T", banner)
	if stringCheck != "string" {
		t.Errorf("Got: %v\n", banner)
	}
}

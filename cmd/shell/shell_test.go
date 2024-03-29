package shell

import (
	"io"
	"reflect"
	"testing"
)

func TestCBShellGetter(t *testing.T) {
	var fakeShellGetter TestShellGetter

	//expected := reflect.TypeOf(ShellInterface)
	actual, err := fakeShellGetter.GetConnectBackInitiatedShell()

	if err != nil {
		t.Errorf("Test Shell Getter returned error: %v\n", err)
	}

	_, success := actual.(*TestShell)
	if !success {
		t.Errorf("Test Shell Getter Failed - Expected: *TestShell, Actual: %v\n", actual)
	}
}

func TestGetStdoutReader(t *testing.T) {
	var fakeShell TestShell

	expected := reflect.TypeOf((*io.ReadCloser)(nil)).Elem()
	actual, err := fakeShell.GetStdoutReader()
	if err != nil {
		t.Errorf("Test Stdout Getter - Error returned from function: %v\n", err)
	}

	if reflect.TypeOf(actual).Elem() != expected {
		t.Errorf("Test Stdout Getter - Expected %T, Got: %v\n", expected, actual)
	}
}

func TestGetStderrReader(t *testing.T) {
	var fakeShell TestShell

	expected := reflect.TypeOf((*io.ReadCloser)(nil)).Elem()
	actual, err := fakeShell.GetStderrReader()
	if err != nil {
		t.Errorf("Test Stderr Getter - Error returned from function: %v\n", err)
	}

	if reflect.TypeOf(actual).Elem() != expected {
		t.Errorf("Test Stderr Getter - Expected %T, Got: %v\n", expected, actual)
	}
}

func TestGetStdinWriter(t *testing.T) {
	var fakeShell TestShell

	expected := reflect.TypeOf((*io.WriteCloser)(nil)).Elem()
	actual, err := fakeShell.GetStdinWriter()
	if err != nil {
		t.Errorf("Test Stdin Getter - Error returned from function: %v\n", err)
	}

	if reflect.TypeOf(actual).Elem() != expected {
		t.Errorf("Test Stdin Getter - Expected %T, Got: %v\n", expected, actual)
	}
}

func TestStartShell(t *testing.T) {
	var fakeShell TestShell
	var fakeError error

	expected := reflect.TypeOf(fakeError)
	actual := fakeShell.StartShell()

	if reflect.TypeOf(actual) != expected {
		t.Errorf("Test Start Shell - Expected: nil, Got: %v\n", actual)
	}
}

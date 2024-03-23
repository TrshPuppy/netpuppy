package conn

import (
	"io"
	"reflect"
	"slices"
	"testing"
)

func TestGetConnectionFromClient(t *testing.T) {
	// Arrange
	var rPort int = 69
	var address string = "69.69.69.69"
	var shell bool = true
	//var testConnectionGetter = TestConnectionGetter{}
	var testConnection TestConnectionGetter

	// Act
	var socket = testConnection.GetConnectionFromClient(rPort, address, shell)

	// Assert
	testClientSocket, success := socket.(TestSocket) // Have to do type assertion to make sure TestSocket & not Socket is returned:
	if !success {
		t.Errorf("Test Client Socket Type Assertion - Got: %v, Expected: TestSocket\n", socket)
	}

	if testClientSocket.Port != rPort {
		t.Errorf("Test Client Socket Port - Got: %v, Expected: %v\n", testClientSocket.Port, rPort)
	}

	if testClientSocket.Address != address {
		t.Errorf("Test Client Socket Address - Got %v, Expected: %v\n", testClientSocket.Address, address)
	}
}

func TestGetConnectionFromListener(t *testing.T) {
	var rPort int = 69
	var address string = "0.0.0.0"
	var testConnection TestConnectionGetter

	var socket = testConnection.GetConnectionFromListener(rPort, address)
	// Type assertion:
	testListenerSocket, success := socket.(TestSocket)
	if !success {
		t.Errorf("Test Listener Socket Type Assertion - Got: %v, Expected: TestSocket\n", socket)
	}

	if testListenerSocket.Port != rPort {
		t.Errorf("Test Listener Socket Port - Got: %v, Expected: %v\n", testListenerSocket.Port, rPort)
	}

	if testListenerSocket.Address != address {
		t.Errorf("Test Listener Socket Address = Got: %v, Expected: %v\n", testListenerSocket.Address, address)
	}
}

func TestSocketRead(t *testing.T) {
	testReadByteArr := []byte("tiddies")
	var testReadErr error

	var fakeSocket TestSocket

	readReturn, readErr := fakeSocket.Read()
	if readErr != testReadErr {
		t.Errorf("Test Error readErr - Got: %v, Expected: error\n", readErr)
	}

	if !slices.Equal(readReturn, testReadByteArr) {
		t.Errorf("Test Read readReturn - Got: %v, Expected: []byte\n", readReturn)
	}
}

func TestSocketClose(t *testing.T) {
	var fakeSocket TestSocket

	expected := reflect.TypeOf((error)(nil))
	actual := fakeSocket.Close()

	if reflect.TypeOf(actual) != expected {
		t.Errorf("Test Socket Close - expected: %T, actual: %v\n", expected, actual)
	}
}

func TestSocketWrite(t *testing.T) {
	var fakeSocket TestSocket
	testWriteThis := []byte("tiddies")

	expected := reflect.TypeOf(69)
	writeReturn, err := fakeSocket.Write(testWriteThis)
	if err != nil {
		t.Errorf("Test Write return - Got: %v, Expected: error\n", err)
	}

	if reflect.TypeOf(writeReturn) != expected {
		t.Errorf("Test Write return - Got: %T, Expected: %v\n", writeReturn, expected)
	}
}

func TestGetReader(t *testing.T) {
	var fakeSocket TestSocket

	actualReturn := fakeSocket.GetReader()
	expected := reflect.TypeOf((*io.Reader)(nil)).Elem()

	if reflect.TypeOf(actualReturn).Elem() != expected {
		t.Errorf("Test GetReader return - Got: %T, Expected: %v\n", actualReturn, expected)
	}
}

func TestGetWriter(t *testing.T) {
	var fakeSocket TestSocket

	actual := fakeSocket.GetWriter()
	expected := reflect.TypeOf((*io.Writer)(nil)).Elem()

	if reflect.TypeOf(actual).Elem() != expected {
		t.Errorf("Test GetWriter return - Got: %T, Expected: %v\n", actual, expected)
	}
}

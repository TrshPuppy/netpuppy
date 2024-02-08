package utils

import (
	"testing"
)

func TestCreateConnectBackPeer(t *testing.T) {
	// Arrange (create dummy connectBack peer)
	var dummyPort int = 69
	var dummyAddress string = "69.69.69.69"
	var dummyListen bool = false
	var dummyConnectionType string = "connect_back"

	// Act
	got := CreatePeer(dummyPort, dummyAddress, dummyListen)

	// Assert
	if got.RPort != dummyPort {
		t.Errorf("got %v, wanted %v\n", got.LPort, dummyPort)
	}

	if got.Address != dummyAddress {
		t.Errorf("Got %v, wanted %v\n", got.Address, dummyAddress)
	}

	if got.ConnectionType != dummyConnectionType {
		t.Errorf("Got %v, wanted %v\n", got.ConnectionType, dummyConnectionType)
	}
}

func TestCreateOffensePeer(t *testing.T) {
	// Arrange (create dummy connectBack peer)
	var dummyPort int = 40404
	var dummyAddress string = "0.0.0.0"
	var dummyListen bool = true
	var dummyConnectionType string = "offense"

	// Act
	got := CreatePeer(dummyPort, dummyAddress, dummyListen)

	// Assert
	if got.RPort != dummyPort {
		t.Errorf("got %v, wanted %v\n", got.LPort, dummyPort)
	}

	if got.Address != dummyAddress {
		t.Errorf("Got %v, wanted %v\n", got.Address, dummyAddress)
	}

	if got.ConnectionType != dummyConnectionType {
		t.Errorf("Got %v, wanted %v\n", got.ConnectionType, dummyConnectionType)
	}
}

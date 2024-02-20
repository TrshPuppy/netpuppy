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
	var dummyShell bool = true

	// Act
	got := CreatePeer(dummyPort, dummyAddress, dummyListen, dummyShell)

	// Assert
	if got.RPort != dummyPort {
		t.Errorf("CBPeer.RPort - Got: %v, Expected: %v\n", got.RPort, dummyPort)
	}

	if got.Address != dummyAddress {
		t.Errorf("CBPeer.Address - Got: %v, Expected: %v\n", got.Address, dummyAddress)
	}

	if got.ConnectionType != dummyConnectionType {
		t.Errorf("CBPeer.ConnectionType - Got: %v, Expected: %v\n", got.ConnectionType, dummyConnectionType)
	}

	if got.Shell != dummyShell {
		t.Errorf("CBPeer.Shell - Got: %v, Expected: %v\n", got.Shell, dummyShell)
	}
}

func TestCreateOffensePeer(t *testing.T) {
	// Arrange (create dummy connectBack peer)
	var dummyPort int = 40404
	var dummyAddress string = "0.0.0.0"
	var dummyListen bool = true
	var dummyConnectionType string = "offense"
	var dummyShell bool = false

	// Act
	got := CreatePeer(dummyPort, dummyAddress, dummyListen, dummyShell)

	// Assert
	if got.LPort != dummyPort {
		t.Errorf("OffensePeer.LPort - Got: %v, Expected: %v\n", got.LPort, dummyPort)
	}

	if got.Address != dummyAddress {
		t.Errorf("OffensePeer.Address - Got: %v, Expected: %v\n", got.Address, dummyAddress)
	}

	if got.ConnectionType != dummyConnectionType {
		t.Errorf("OffencePeer.ConnectionType - Got: %v, Expected: %v\n", got.ConnectionType, dummyConnectionType)
	}

	if got.Shell != dummyShell {
		t.Errorf("OffensePeer.Shell - Got: %v, Expected: %v\n", got.Shell, dummyShell)
	}
}

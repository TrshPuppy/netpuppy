package utils

import (
	"net"
	"os/exec"
)

// Depending on input, create this peer's type:
type Peer struct {
	ConnectionType string
	RPort          int
	LPort          string
	Address        string
	Connection     net.Conn
	CbShell        *exec.Cmd
}

func CreatePeer(port int, address string, listen bool) Peer {
	var thisPeer Peer

	if listen {
		thisPeer = getOffense(port)
	} else {
		thisPeer = getConnectBack(port, address)
	}
	return thisPeer
}

func getOffense(port int) Peer {
	offensePeer := Peer{RPort: port, Address: "0.0.0.0", ConnectionType: "offense"}
	return offensePeer
}

func getConnectBack(port int, address string) Peer {
	connectBackPeer := Peer{RPort: port, Address: address, ConnectionType: "connect_back"}
	return connectBackPeer
}

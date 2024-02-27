package utils

import "fmt"

// Depending on input, create this peer's type:
type Peer struct {
	ConnectionType string
	RPort          int
	LPort          int
	Address        string
	Connection     Socket
	Shell          bool
	ShellProcess   BashShell
	ReportTo       string
}

func CreatePeer(port int, address string, listen bool, shell bool) *Peer {
	var thisPeer *Peer // POINTER: the functions which initialize thiePeer are returning addresses to the instances of Peer they create

	if listen { // Offense peer
		thisPeer = getOffense(port, shell)
	} else { // Connect-back peer
		thisPeer = getConnectBack(port, address, shell)
	}

	// If we are the connect-back & --shell was given, report to bash (socket data redirected to bash)
	if shell && !listen {
		thisPeer.ReportTo = "bash"
	} else {
		thisPeer.ReportTo = "user"
	}

	fmt.Printf("Address of thisPeer in CreatePeer: %v\n", &thisPeer)
	return thisPeer
}

func getOffense(port int, shell bool) *Peer {
	var offensePeer Peer = Peer{LPort: port, Address: "0.0.0.0", ConnectionType: "offense", Shell: shell}

	return &offensePeer // POINTER: return the address of the offensePeer instance (instead of a copy)
}

func getConnectBack(port int, address string, shell bool) *Peer {
	var connectBackPeer Peer = Peer{RPort: port, Address: address, ConnectionType: "connect-back", Shell: shell}

	return &connectBackPeer // POINTER: return the address of the connectBackPeer instance (instead of copy)
}

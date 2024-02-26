package utils

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
	var thisPeer Peer

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
	return &thisPeer
}

func getOffense(port int, shell bool) Peer {
	offensePeer := Peer{LPort: port, Address: "0.0.0.0", ConnectionType: "offense", Shell: shell}
	return offensePeer
}

func getConnectBack(port int, address string, shell bool) Peer {
	connectBackPeer := Peer{RPort: port, Address: address, ConnectionType: "connect-back", Shell: shell}
	return connectBackPeer
}

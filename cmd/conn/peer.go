package conn

// Depending on user flags given, create this peer's type:
type Peer struct {
	ConnectionType string
	RPort          int
	LPort          int
	Address        string
	Shell          bool
}

// Create the peer struct which will guide the rest of the logic for NetPuppy:
func CreatePeer(port int, address string, listen bool, shell bool) *Peer {
	// POINTER: the functions which initialize thisPeer are returning addresses to the instances of Peer they create
	var thisPeer *Peer

	// If listen is true, we are the Offense peer (server):
	if listen {
		thisPeer = getOffense(port, shell)
	} else { // else, we are the Connect-back peer (client):
		thisPeer = getConnectBack(port, address, shell)
	}

	return thisPeer
}

// The Offense-type peer is the server and will listen on the any address:
func getOffense(port int, shell bool) *Peer {
	// Create offense-type instance of Peer:
	var offensePeer *Peer = &Peer{LPort: port, Address: "0.0.0.0", ConnectionType: "offense", Shell: shell}

	return offensePeer
}

// The Connect-back peer is the client & will connect to the given IP address & port:
func getConnectBack(port int, address string, shell bool) *Peer {
	// Create connect-back instance of Peer:
	var connectBackPeer *Peer = &Peer{RPort: port, Address: address, ConnectionType: "connect-back", Shell: shell}

	return connectBackPeer
}

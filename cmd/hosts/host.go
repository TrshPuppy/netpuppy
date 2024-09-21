package hosts

import (
	"errors"
	"os"

	"github.com/trshpuppy/netpuppy/cmd/conn"
	"github.com/trshpuppy/netpuppy/cmd/shell"
)

type Host interface {
	// WriteTo()
	// ReadFrom()
	Start() error
}

type OffensiveHost struct {
	socket conn.SocketInterface
	stdin  *os.File
	stdout *os.File
}

type ConnectBackHost struct {
	socket       conn.SocketInterface
	masterDevice *os.File
	ptsDevice    *os.File
	shell        shell.RealShell
}

func NewHost(peer *conn.Peer) (Host, error) {
	// Based on the connection type requested by user,
	// get the socket:
	var c conn.ConnectionGetter
	// var host Host

	switch peer.ConnectionType {
	case "offense":
		socket, err := c.GetConnectionFromListener(peer.LPort, peer.Address)
		if err != nil {
			return nil, errors.New("Error getting Listener socket: " + err.Error())
		}

		// We've got the socket, attach it to .... host struct?
		host := OffensiveHost{socket: socket}
		return &host, nil

	case "connect-back":
		socket, err := c.GetConnectionFromClient(peer.RPort, peer.Address, peer.Shell)
		if err != nil {
			return nil, errors.New("Error getting Connect Back socket: " + err.Error())
		}

		// Attach socket to new ConnectBackHost struct:
		host := ConnectBackHost{socket: socket}
		return &host, nil

	default:
		return nil, errors.New("Invalid Shell type when attempting to create host.")
	}

	return nil, nil
}

func (off *OffensiveHost) Start() error {
	return nil
}

func (cb *ConnectBackHost) Start() error {
	return nil
}

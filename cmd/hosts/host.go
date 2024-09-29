package hosts

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/trshpuppy/netpuppy/cmd/conn"
	"github.com/trshpuppy/netpuppy/cmd/shell"
	"github.com/trshpuppy/netpuppy/pkg/ioctl"
)

type Host interface {
	Start(context.Context) error
}

type OffensiveHost struct {
	socket conn.SocketInterface
	stdin  *os.File
	stdout *os.File
	stderr *os.File
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

	switch peer.ConnectionType {
	case "offense":
		socket, err := c.GetConnectionFromListener(peer.LPort, peer.Address)
		if err != nil {
			return nil, errors.New("Error getting Listener socket: " + err.Error())
		}

		// We've got the socket, attach it to .... host struct?
		host := OffensiveHost{socket: socket, stdin: os.Stdin, stdout: os.Stdout, stderr: os.Stderr}
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
}

func (off *OffensiveHost) Start(pCtx context.Context) error {
	// Enable raw mode:
	oGTermios, errno := ioctl.EnableRawMode(int(off.stdin.Fd()))
	if errno != 0 {
		return fmt.Errorf("Error enabling termios raw mode; returned error code: %s\n,", errno)
	}
	// Use the oGTermios structure w/ a defer of DisableRawMode() to reset the terminal before exiting:
	defer ioctl.DisableRawMode(int(off.stdin.Fd()), oGTermios)

	// Create child context:
	childContext, chCancel := context.WithCancel(pCtx)
	defer chCancel()

	// Create wait group for go routines:
	var wg sync.WaitGroup
	wg.Add(2)

	// GO ROUTINES:
	go func() { // Read socket & copy to stdout:
		defer wg.Done()
		defer chCancel()

		for {
			select {
			case <-childContext.Done():
				fmt.Printf("child context done in offense go routine\n")
				return
			default:
				dataFromSocket, err := off.socket.Read()
				if len(dataFromSocket) > 0 {
					// Write data to stdout:
					_, err = off.stdout.Write(dataFromSocket)
					if err != nil {
						fmt.Printf("Error writing data from socket to Offense stdout: %v\n", err)
						return
					}
				}

				if err != nil {
					if errors.Is(err, io.EOF) {
						fmt.Printf("Socket closed by other peer\n")
						return
					}
					fmt.Printf("Error reading from socket to Offense stdout: %v\n", err)
					return
				}
			}
		}
	}()

	go func() { // Read stdin & write to socket:
		defer wg.Done()
		defer chCancel()

		// Make buffer
		buffer := make([]byte, 1)

		// Read stdin byte by byte
		for {
			select {
			case <-childContext.Done():
				fmt.Printf("Child received done signal, quitting out of offensive start routine\n")
				return
			default:
				i, TIT_BY_BOO := off.stdin.Read(buffer)
				if TIT_BY_BOO != nil {
					if errors.Is(TIT_BY_BOO, io.EOF) {
						continue
					} else {
						fmt.Printf("Error reading from stdin: %v\n", TIT_BY_BOO)
						return
					}
				}

				// Write to socket:
				_, err := off.socket.WriteShit(buffer[:i])
				if err != nil {
					fmt.Printf("Error writing Stdin to socket: %v\n", err)
					return
				}
			}
		}
	}()

	// Call wait after go routines b/c it's going to block:
	wg.Wait()

	errno = ioctl.DisableRawMode(int(off.stdin.Fd()), oGTermios)
	if errno != 0 {
		return errors.New("Error disabling raw mode on Offense termios, error code: " + errno.Error())
	}

	err := off.socket.Close()
	if err != nil {
		return errors.New("Error closing socket on Offense host: " + err.Error())
	}

	err = off.stdin.Close()
	if err != nil {
		return errors.New("Error closing stdin on Offense host: " + err.Error())
	}

	err = off.stdout.Close()
	if err != nil {
		return errors.New("Error closing stdout on Offense host: " + err.Error())
	}

	err = off.stderr.Close()
	if err != nil {
		return errors.New("Error closing stderr on Offense host: " + err.Error())
	}

	return nil
}

func (cb *ConnectBackHost) Start(pCtx context.Context) error {
	// read from socket master
	// read from master to socket
	return nil
}

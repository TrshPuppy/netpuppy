/* TP U ARE HERE:
- re-organize the channels given new cb shell
- decide if the shell should be an option vs automatic
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

func readFromSocket(socketReader chan<- []byte, connection utils.Socket) {
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
		socketReader <- dataBytes
	}
}

func startHelperShell() (*exec.Cmd, error) {
	bashPath, err := exec.LookPath(`/bin/bash`)
	if err != nil {
		fmt.Printf("Error finding bash path: %v\n", err)
		os.Stderr.WriteString(" " + err.Error() + "\n")
		return nil, err
	}
	bCmd := exec.Command(bashPath)
	var erR error = bCmd.Start()

	return bCmd, erR
}

func runApp(c utils.ConnectionGetter) {
	// Parse flags from user, attach to struct:
	flagStruct := utils.GetFlags()

	// Print banner:
	fmt.Printf("%s", utils.Banner())

	// Create peer instance based on user input:
	var socket utils.Socket
	thisPeer := utils.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen, flagStruct.Shell)

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

	// Update banner w/ missing port:
	// var missingPortInBanner = utils.PrintMissingPortToBanner(thisPeer.ConnectionType, thisPeer.Connection)
	// fmt.Println(missingPortInBanner)

	// Start SIGINT go routine:
	// Start channel to listen for SIGINT:
	listenForSIGINT(thisPeer)

	// Start go routines and channels to read socket and user input:
	// IO read & socket write channels (user input will be written to socket)
	ioReader := make(chan string)

	// IO write & socket read channels (messages from socket will be printed to stdout)
	socketReader := make(chan []byte)

	go readUserInput(ioReader)
	go readFromSocket(socketReader, thisPeer.Connection)

	for {
		select {
		case userInput := <-ioReader:
			_, err := thisPeer.Connection.Write([]byte(userInput))

			if err != nil {
				// Quit here?
				fmt.Printf("Error in userInput select: %v\n", err)
				os.Stderr.WriteString(" " + err.Error() + "\n")
			}
		case socketIncoming := <-socketReader:
			_, err := os.Stdout.Write(socketIncoming)
			if err != nil {
				// Quit here?
				fmt.Printf("Error in writing to stdout: %v\n", err)
				os.Stderr.WriteString(" " + err.Error() + "\n")
			}
		default:
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func listenForSIGINT(thisPeer utils.Peer) {
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

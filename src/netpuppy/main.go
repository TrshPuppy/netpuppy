/* TP U ARE HERE:
- re-organize the channels given new cb shell
- decide if the shell should be an option vs automatic
*/

package main

import (
	"bufio"
	"netpuppy/utils"
	"time"

	//"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	// NetPuppy modules:
)

func sum(a int, b int) int {
	s := a + b
	return s
}

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

func startHelperShell() (*exec.Cmd, error) { // @Trauma_X_Sella 'connection'
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

	// Make connection:
	if thisPeer.ConnectionType == "offense" {
		socket = c.GetConnectionFromListener(thisPeer.RPort, thisPeer.Address)
	} else {
		socket = c.GetConnectionFromClient(thisPeer.RPort, thisPeer.Address)
	}

	// Connect socket connection to peer
	thisPeer.Connection = socket

	// Update user:
	var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
	fmt.Println(updateUserBanner)

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
	//var fakeConnection utils.FakeConnectionGetter
	var realConnection utils.RealConnectionGetter
	runApp(realConnection)
	// s := sum(2, 3)
	// st := fmt.Sprintf("%v", s)
	// fmt.Printf("sum: %v\n", st)
	//
	// // Get flags passed to NetPuppy by user:

	//
	// // Print banner (runs if flags are parsed w/p error):
	// fmt.Printf("%s", utils.Banner())
	//
	// // Get STDIN and save to a variable we can use if we need:
	// // stdinReader := bufio.NewReader(os.Stdin)
	// // stdin, _ := stdinReader.ReadString('\n')
	// // fmt.Printf("STDIN = %v", stdin) // Keep for now to avoid golang complaints about unused vars.
	//
	// // Create and return peer based on user's input:
	// thisPeer := utils.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen, flagStruct.Shell)
	//
	// // Now that we have our peer: try to make connection:
	// asyncio_rocks := utils.GetConnection(thisPeer.ConnectionType, thisPeer.RPort, thisPeer.Address) // @0xTib3rius 'connection'
	//
	// // Attach connection to peer struct and get local port number from connection:
	// thisPeer.Connection = asyncio_rocks
	// localPortArr := strings.Split(thisPeer.Connection.LocalAddr().String(), ":")
	// localPort := localPortArr[len(localPortArr)-1]
	// thisPeer.LPort = localPort
	//
	// // Update user:
	// var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
	// fmt.Println(updateUserBanner)
	//
	// // Start channel to listen for SIGINT:
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	//
	//	go func() {
	//		// If SIGINT: close connection, exit w/ code 2
	//		for sig := range signalChan {
	//			if sig.String() == "interrupt" {
	//				fmt.Printf("signal: %v\n", sig)
	//				thisPeer.Connection.Close()
	//				os.Exit(2)
	//			}
	//		}
	//	}()
	//
	// // If we're the connect_back peer, start 'helper' shell:
	//
	//	if thisPeer.ConnectionType == "connect_back" {
	//		connectBackShell, shellStartErr := startHelperShell()
	//		if shellStartErr != nil {
	//			fmt.Printf("Error starting shell process: %v\n", shellStartErr)
	//			thisPeer.Connection.Close()
	//			os.Stderr.WriteString(" " + shellStartErr.Error() + "\n")
	//			os.Exit(1)
	//		}
	//		thisPeer.ShellProcess = connectBackShell
	//	}
	//
	// // IO read & socket write channels (user input will be written to socket)
	// ioReader := make(chan string)
	//
	// // IO write & socket read channels (messages from socket will be printed to stdout)
	// socketReader := make(chan []byte)
	//
	// go readUserInput(ioReader)
	// go readFromSocket(socketReader, thisPeer.Connection)
	//
	//	for {
	//		select {
	//		case userInput := <-ioReader:
	//			_, err := thisPeer.Connection.Write([]byte(userInput))
	//			if err != nil {
	//				// Quit here?
	//				fmt.Printf("Error in userInput select: %v\n", err)
	//				os.Stderr.WriteString(" " + err.Error() + "\n")
	//			}
	//		case socketIncoming := <-socketReader:
	//			_, err := os.Stdout.Write(socketIncoming)
	//			if err != nil {
	//				// Quit here?
	//				fmt.Printf("Error in writing to stdout: %v\n", err)
	//				os.Stderr.WriteString(" " + err.Error() + "\n")
	//			}
	//		default:
	//			time.Sleep(300 * time.Millisecond)
	//		}
	//	}
	//
	// return
}

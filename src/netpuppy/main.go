package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
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

//func writeToSocket(ioReader <-chan string, socketWriter chan<- bool, connection net.Conn) {
//	inputToSend := <-ioReader
//
//	_, err := connection.Write([]byte(inputToSend))
//	if err != nil {
//		socketWriter <- false
//		return
//	}
//	socketWriter <- true
//}

func readFromSocket(socketReader chan<- []byte, connection net.Conn) {
	// Read from connection socket:

	for {
		dataBytes, err := bufio.NewReader(connection).ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error reading from socket: %v\n", err)
			return
		}
		socketReader <- dataBytes
	}
}

func writeToStdout(ioWriter chan<- bool, socketReader <-chan []byte) {
	//fmt.Printf("Socket writer started...\n")
	bytesToWrite := <-socketReader

	_, err := os.Stdout.Write(bytesToWrite)
	if err != nil {
		ioWriter <- false
		return
	}

	ioWriter <- true
}

func main() {
	// Set flag values based on input:
	listenFlag := flag.Bool("l", false, "put NetPuppy in listen mode")
	hostFlag := flag.String("H", "0.0.0.0", "target host IP address to connect to")
	turdnuggies := flag.Int("p", 40404, "target port") // portFlag @Trauma_x_Sella

	// Parse command line arguments:
	//                                            error?
	flag.Parse()

	// Print banner:
	fmt.Printf("%s", utils.Banner())

	// Depending on input, create this peer's type:
	type peer struct {
		connection_type string
		rPort           int
		lPort           string
		address         string
		connection      net.Conn
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
			os.Stderr.WriteString(err1.Error())
			os.Exit(1)
		}

		defer listener.Close() // Ensure the listener closes when main() returns

		asyncio_rocks, err = listener.Accept()
		if err != nil {
			os.Stderr.WriteString(err.Error())
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
				os.Stderr.WriteString(err.Error())
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
	// IO read & socket write channels (user input will be written to socket)
	ioReader := make(chan string)
	//	socketWriter := make(chan bool)

	// IO write & socket read channels (messages from socket will be printed to stdout)
	//ioWriter := make(chan bool)
	socketReader := make(chan []byte)

	go readUserInput(ioReader)
	//go writeToSocket(ioReader, socketWriter, thisPeer.connection)

	go readFromSocket(socketReader, thisPeer.connection)
	//	go writeToStdout(ioWriter, socketReader)

	for {
		select {
		case userInput := <-ioReader:
			_, err := thisPeer.connection.Write([]byte(userInput))
			if err != nil {
				// Quit here?
				fmt.Printf("Error in userInput select: %v\n", err)
			}
		case socketIncoming := <-socketReader:
			_, err := os.Stdout.Write(socketIncoming)
			if err != nil {
				// Quit here?
				fmt.Printf("Error in writing to stdout: %v\n", err)
			}
		default:
			time.Sleep(300 * time.Millisecond)
			//fmt.Printf("Default: slept for 300 ms\n")
		}
	}

	//		select {
	//
	//		}

	// Check for success writing to socket and stdout:
	//	socketWriteSuccess := <-socketWriter
	//	if !socketWriteSuccess {
	//		fmt.Printf("Error writing to socket! \n")
	//	}

	//	stdOutWriteSuccess := <-ioWriter
	//	if !stdOutWriteSuccess {
	//		fmt.Printf("Error writing to stdout! \n")
	//	}
	//}

	/*
		if -l is on,
			net.Listen('tcp', PORT)
			set connection address for socket to any
		if not
			connection address = host flag


		struct/ objsect thing (this peer)
			- connect back (executed on the target)
				- start the subprocess
			- offense (exe on hacker machine)
				- keeep taking user input


			- method:
				func make connection(){
					if this.type = offense:
						connection = net.Listener
						(needs Accept() to actually become a connection)
					else:
						connection = net.Dial
				}
	*/

	// Try to create connection:
	return
}

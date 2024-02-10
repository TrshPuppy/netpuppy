package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Socket interface { // dummy net.Conn
	Read() ([]byte, error)
	Write([]byte) (int, error)
	Close() error
}

type ConnectionGetter interface {
	GetConnectionFromListener(int, string) Socket
	GetConnectionFromClient(int, string) Socket
}

// FAKE
type FakeConnectionGetter struct {
}

func (c FakeConnectionGetter) GetConnectionFromClient(rPort int, address string) Socket {
	clientConnection := ImaginaryConnection{Port: rPort, Address: address}
	return clientConnection
}

func (c FakeConnectionGetter) GetConnectionFromListener(rPort int, address string) Socket {
	listenerConnection := ImaginaryConnection{Port: rPort, Address: address}
	return listenerConnection
}

type ImaginaryConnection struct {
	Port    int
	Address string
}

func (i ImaginaryConnection) Read() ([]byte, error) {
	var imaginaryByteArr []byte
	var imaginaryErr error

	return imaginaryByteArr, imaginaryErr
}

func (i ImaginaryConnection) Write([]byte) (int, error) {
	return 22, nil
}

func (i ImaginaryConnection) Close() error {
	var imaginaryCloseErr error
	return imaginaryCloseErr
}

// REAL
type RealSocket struct {
	realSocket net.Conn
}

type RealConnectionGetter struct {
}

func (s RealSocket) Read() ([]byte, error) {
	dataBytes, err := bufio.NewReader(s.realSocket).ReadBytes('\n')
	return dataBytes, err
}

func (s RealSocket) Write(userInput []byte) (int, error) {
	w, err := s.realSocket.Write(userInput)
	return w, err
}

func (s RealSocket) Close() error {
	err := s.realSocket.Close()
	return err
}

func (r RealConnectionGetter) GetConnectionFromClient(rPort int, address string) Socket {
	var clientConnection net.Conn
	var err error

	remoteHost := fmt.Sprintf("%v:%v", address, rPort)
	clientConnection, err = net.Dial("tcp", remoteHost)
	if err != nil {
		fmt.Printf("Error creating connection: %v\n", err)
		os.Stderr.WriteString(" " + err.Error() + "\n")
		os.Exit(1)
	}

	clientSocket := RealSocket{realSocket: clientConnection}
	return clientSocket
}

func (r RealConnectionGetter) GetConnectionFromListener(rPort int, address string) Socket {
	var listenerConnection net.Conn
	var err error

	listener, err1 := net.Listen("tcp", fmt.Sprintf(":%v", rPort))
	if err1 != nil {
		fmt.Printf("Error when creating listener connection: %v\n", err1)
		os.Stderr.WriteString(" " + err1.Error() + "\n")
		os.Exit(1)
	}

	defer listener.Close()

	listenerConnection, err = listener.Accept()
	if err != nil {
		fmt.Printf("Error when creating listener connection: %v\n", err)
		os.Stderr.WriteString(" " + err.Error() + "\n")
		os.Exit(1)
	}

	listenerSocket := RealSocket{realSocket: listenerConnection}
	return listenerSocket
}

func GetConnection(cType string, rPort int, address string) net.Conn {
	var connection net.Conn
	var err error

	// Make connection based on peer type:
	if cType == "offense" {
		connection, err = createListenerConnection(rPort)
	} else {
		connection, err = createClientConnection(address, rPort)
	}

	if err != nil {
		fmt.Printf("Error creating connection: %v\n", err)
		os.Stderr.WriteString(" " + err.Error() + "\n")
		os.Exit(1)
	}

	return connection
}

func createListenerConnection(rPort int) (net.Conn, error) {
	var connection net.Conn
	var err error

	snake_case, err1 := net.Listen("tcp", fmt.Sprintf(":%v", rPort)) // @himselfe 'listener'
	if err1 != nil {
		fmt.Printf("Error when creating listener connection: %v\n", err1)
		os.Stderr.WriteString(" " + err.Error() + "\n")
		os.Exit(1)
	}

	defer snake_case.Close()

	connection, err = snake_case.Accept()
	return connection, err
}

func createClientConnection(address string, rPort int) (net.Conn, error) {
	var yuleLog net.Conn // @Trauma_X_Sella 'connection'
	var err error

	remoteHost := fmt.Sprintf("%v:%v", address, rPort)
	yuleLog, err = net.Dial("tcp", remoteHost)

	// If there is an err, try the host address as ipv6 (need to add [] around string):
	return yuleLog, err
}

/*

search for word with /foobar, hit enter go back to normal mode, type cgn to replace next search match, then you can hit n to find the next instance and hit . to replace the next instance or hit n to go to the next
*/

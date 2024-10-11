/*
	This code is intentionally redundant!
	In order to implement unit testing on this code
	we want to avoid creating REAL socket connections
	while running tests.

	There are two shared interfaces (Socket & ConnectionGetter)
	which allow	us to isolate the actual creation of a socket
	& still write tests for this code. RealSocket & TestSocket
	are checked against the Socket interface. RealConnectionGetter
	& TestConnectionGetter are checked against the ConnectionGetter
	interface.

	TEST code:
	...includes the TestSocket & TestConnectionGetter structs
	which implement a fake socket using fake connection methods.
	The	methods will be called for unit testing and they return the
	same value-types as the REAL code.

	REAL code:
	The only code which can be used to interact w/ the REAL socket
	connection is RealSocket struct. Both the Real and Test sockets
	are verified by the Socket interface.

	The REAL socket is created and attached to the RealSocket struct
	in the GetConnectionFromClient GetConnectionFromListener methods
	attached to the RealConnectionGetter struct (being checked against
	ConnectionGetterInterface).

	Once created and attached to RealSocket, the socket can be read from/ written
	to and closed by using the Read() Write() & Close() methods. The real socket
	(net.Conn) object shouldn't be handled anywhere else (or testing will break).

	TESTING:
	If something changes about the code in a way that deviates from the
	blueprinting done by the interfaces, tests will fail.
*/

package conn

import (
	"fmt"
	"io"
	"net"
)

// SHARED Code:
// ..... The Socket and ConnectionGetter interfaces are used by both real & test code:
type SocketInterface interface {
	// Used to check real (RealSocket) & test (TestSocket) structs
	Read() ([]byte, error)
	WriteShit([]byte) (int, error)
	Close() error
	GetReader() io.Reader
	GetWriter() io.Writer
}

type ConnectionGetter interface {
	// Used to check the real (RealConnectionGetter) & test (TestConnectionGetter) structs:
	GetConnectionFromListener(int, string, bool) (SocketInterface, error)
	GetConnectionFromClient(int, string, bool) (SocketInterface, error)
}

// REAL code:
type RealSocket struct { // This is the only code which holds the ACTUAL net connection:
	Socket net.Conn
}

type RealConnectionGetter struct {
	// Leave empty
}

// Read from the ACTUAL socket on RealSocket struct:
func (s *RealSocket) Read() ([]byte, error) {
	var buffer []byte = make([]byte, 1024)
	var numberOfBytesSent int = 0
	var err error

	numberOfBytesSent, err = s.Socket.Read(buffer)
	return buffer[:numberOfBytesSent], err
}

// Write to the ACTUAL socket on RealSocket struct:
func (s *RealSocket) WriteShit(input []byte) (int, error) {
	var writeSuccess int
	var err error

	writeSuccess, err = s.Socket.Write(input)
	return writeSuccess, err
}

// Close the ACTUAL socket:
func (s *RealSocket) Close() error {
	var err error = s.Socket.Close()
	return err
}

func (s *RealSocket) GetReader() io.Reader {
	reader := s.Socket
	return reader
}

func (s *RealSocket) GetWriter() io.Writer {
	writer := s.Socket
	return writer
}

// These next 2 function create the ACTUAL socket and attach the connection to RealSocket
// ..... Create client-type socket & attach to RealSocket:
func (r RealConnectionGetter) GetConnectionFromClient(rPort int, address string, shell bool) (SocketInterface, error) {
	var clientConnection net.Conn
	var err error
	var pointerToRealSocket *RealSocket

	remoteHost := net.JoinHostPort(address, fmt.Sprintf("%v", rPort)) // Re-writing for DialTCP

	// Get client connectiokjn:
	clientConnection, err = net.Dial("tcp", remoteHost)
	if err != nil {
		fmt.Printf("ERROR ON REV SHELL DIAL CONNECT: %v\n", err)
		return nil, err
	}

	// Attach connection to RealSocket & get the pointer to the instance:
	pointerToRealSocket = &RealSocket{Socket: clientConnection}

	return pointerToRealSocket, nil
}

// Creat listener-type socket & attach to RealSocket:
func (r RealConnectionGetter) GetConnectionFromListener(rPort int, address string, shell bool) (SocketInterface, error) {
	var listenerConnection net.Conn
	var err error
	var localPort string = fmt.Sprintf(":%v", rPort)

	// Listener created first:
	listener, err1 := net.Listen("tcp", localPort)
	if err1 != nil {
		fmt.Printf("ERROR IN LISTENER: %v\n", err1)
		// listener.Close()
		return nil, err
	}

	// This ensures the listener closes before the function returns:
	defer listener.Close()

	// Create the connection using listener.Accept():
	listenerConnection, err = listener.Accept()
	if err != nil {
		return nil, err
	}

	return &RealSocket{Socket: listenerConnection}, nil
}

// TEST Code:
type TestSocket struct {
	Port    int
	Address string
}

type TestConnectionGetter struct {
	// Leave empty
}

func (i TestSocket) Read() ([]byte, error) {
	testByteArr := []byte("tiddies")
	var testErr error

	return testByteArr, testErr
}

func (i TestSocket) WriteShit([]byte) (int, error) {
	var testInt int = 22
	var testErrorNil error = nil
	return testInt, testErrorNil
}

func (i TestSocket) Close() error {
	var testCloseErr error
	return testCloseErr
}

func (s TestSocket) GetReader() io.Reader {
	var testReader io.Reader
	return testReader
}

func (s TestSocket) GetWriter() io.Writer {
	var testWriter io.Writer
	return testWriter
}

func (c TestConnectionGetter) GetConnectionFromClient(rPort int, address string, shell bool) SocketInterface {
	testClientConnection := TestSocket{Port: rPort, Address: address}
	return testClientConnection
}

func (c TestConnectionGetter) GetConnectionFromListener(rPort int, address string, shell bool) SocketInterface {
	testListenerConnection := TestSocket{Port: rPort, Address: address}
	return testListenerConnection
}

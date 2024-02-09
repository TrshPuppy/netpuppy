package utils

import (
	"fmt"
	"net"
	"os"
)

type Socket interface { // dummy net.Conn
	Read() (int, error)
	Write() (int, error)
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

func (i ImaginaryConnection) Read() (int, error) {
	return 69, nil
}

func (i ImaginaryConnection) Write() (int, error) {

	return 22, nil
}

func (i ImaginaryConnection) Listen() (string, string) {

	return "string1", "string2"
}

func (i ImaginaryConnection) Accept() (string, error) {

	return "string from accept", nil
}

// REAL
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

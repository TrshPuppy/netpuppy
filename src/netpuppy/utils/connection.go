package utils

import (
	"fmt"
	"net"
	"os"
)

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
	var connection net.Conn
	var err error

	remoteHost := fmt.Sprintf("%v:%v", address, rPort)
	connection, err = net.Dial("tcp", remoteHost)

	// If there is an err, try the host address as ipv6 (need to add [] around string):
	return connection, err
}

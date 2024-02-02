package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("Tiddies makde it to the chat!")

	// Set flag values based on input:
	listenFlag := flag.Bool("l", false, "put NetPuppy in listen mode")
	hostFlag := flag.String("H", "0.0.0.0", "target host IP address to connect to")
	portFlag := flag.Int("p", 40404, "target port")

	// Parse command line arguments:
	flag.Parse()

	// Depending on input, create this peer's type:
	type peer struct {
		connection_type string
		port            int
		address         string
		hostname        string
	}

	// Initiate peer struct
	thisPeer := peer{port: *portFlag, address: *hostFlag}

	// If -l was given, create an offense peer:
	if *listenFlag {
		thisPeer.connection_type = "offense"
		thisPeer.address = "0.0.0.0"
	} else {
		thisPeer.connection_type = "connect_back"
	}

	fmt.Printf("The connection type is: %s\n", thisPeer.connection_type)
	fmt.Printf("The host is %s\n", thisPeer.address)
	fmt.Printf("The port is %v\n", thisPeer.port)

	/*
		if -l is on,
			net.Listen('tcp', PORT)
			set connection address for socket to any
		if not
			connection address = host flag


		struct/ objsect thing (this peer)
			- connect back (executed on the target)
				- start the subprocess
			- offense (exe on hacker machine)
				- keeep taking user input





	*/

	// Try to create connection:

}

package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("Tiddies makde it to the chat!")	

	// Set flag values based on input:
	listenFlag := flag.Bool("l", false, "put NetPuppy in listen mode")
	hostFlag := flag.String("H", "127.0.0.1", "target host IP address to connect to")
	portFlag := flag.Int("p", 44444, "target port")

	// Parse command line arguments:
	flag.Parse()
	fmt.Printf("The listen flag is %v\n", *listenFlag)
	fmt.Printf("The host flag is %v\n", *hostFlag)
	fmt.Printf("The port is %v\n", *portFlag)

	/*
		if -l is on, 
			set connection address for socket to any
		if not
			connection address = host flag

	*/

	// Try to create connection:
	if listenFlag* {
		
	}
}

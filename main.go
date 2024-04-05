package main

import (
	// NetPuppy pkgs:
	"github.com/trshpuppy/netpuppy/cmd"
	"github.com/trshpuppy/netpuppy/cmd/conn"
)

func main() {
	// In order to test the connection code w/o creating REAL sockets, Run() handles most of the logic.
	//....... Define a connection getter and hand it to Run(). This will become the socket:
	var connection conn.RealConnectionGetter
	cmd.Run(connection)
}

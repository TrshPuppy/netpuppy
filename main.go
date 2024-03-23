package main

import (
	// NetPuppy pkgs:
	"netpuppy/cmd"
	"netpuppy/utils"
)

func main() {
	// In order to test the connection code w/o creating REAL sockets, Run() handles most of the logic.
	//....... Define a connection getter and hand it to Run(). This will become the socket:
	var connection utils.RealConnectionGetter
	cmd.Run(connection)
}

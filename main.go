package main

import (
	// NetPuppy pkgs:
	"github.com/trshpuppy/netpuppy/cmd"
	"github.com/trshpuppy/netpuppy/cmd/conn"
)

//...........
// TP YOU ARE HERE
//			We need to:
//				use ioctl w/ TCGETS to get the current attributes of netpuppy stdin ? [terminal?]
//				try putting netpuppy into raw mode but try turning off echo first
// 				If we get all of that done, we  need to set netpuppy stdin for the user (not the rev shell peer)
//				into raw mode, so we can use a trie to parse ansi control codes, and then send those
//					(raw/ real time) down the socket
// ..........

func main() {
	// In order to test the connection code w/o creating REAL sockets, Run() handles most of the logic.
	//....... Define a connection getter and hand it to Run(). This will become the socket:

	var connection conn.RealConnectionGetter
	cmd.Run(connection)
}

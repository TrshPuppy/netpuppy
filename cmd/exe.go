package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	// NetPuppy pkgs:
	"github.com/trshpuppy/netpuppy/cmd/conn"
	"github.com/trshpuppy/netpuppy/cmd/hosts"
	"github.com/trshpuppy/netpuppy/cmd/shell"
	"github.com/trshpuppy/netpuppy/utils"
)

// Start the madness:
func Run(c conn.ConnectionGetter) {
	// Make a parent context for the main (Run()) routine:
	parentCtx, pCancel := context.WithCancel(context.Background())
	defer pCancel()

	// Start SIGINT routine before we block Run() with child contexts:
	go func() {
		// If SIGINT: close connection, exit w/ code 2
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		defer signal.Stop(signalChan)

		<-signalChan
		pCancel()
	}()

	// Parse flags from user, attach to struct:
	flagStruct := utils.GetFlags()

	// Create peer instance based on user input:
	var thisPeer *conn.Peer = conn.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen, flagStruct.Shell)
	fmt.Printf("PEER: %v\n", thisPeer)

	// Print banner, but don't print if we are the peer running the shell (ooh sneaky!):
	if !thisPeer.Shell {
		fmt.Printf("%s", utils.Banner())

		// Update user:
		var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
		fmt.Println(updateUserBanner)
	}

	// Make the Host type based on the peer struct:
	var host hosts.Host
	host, err := hosts.NewHost(thisPeer, c)
	if err != nil {
		fmt.Printf("Error trying to create new host: %v\n", err)
	}

	// Once we get the host, call Host.Start():
	err, errCount := host.Start(parentCtx)
	if err != nil {
		fmt.Printf("Error starting host: %v\n", err)
		fmt.Printf("Error Count: %v\n", errCount)
		os.Exit(1337)
	}
}

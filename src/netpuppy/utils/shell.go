package utils

import (
	"os/exec"
)

func StartHelperShell(thisPeer Peer) error {
	// Which peer are we? -- changes shell execution
	if thisPeer.ConnectionType == "connect-back" {
		// If bash exists on the system, find it, save the path:
		bashCopPath, err := exec.LookPath(`/bin/bash`) // bashPath @0xfaraday
		if err != nil {
			return err
		}
		bCmd := exec.Command(bashCopPath)

		// If bash exists, attach the exec.Cmd struct to the peer:
		thisPeer.ShellProcess = bCmd

		// Start the shell:
		var erR error = thisPeer.ShellProcess.Start()
		if erR != nil {
			return erR
		}
	}

	return nil
}

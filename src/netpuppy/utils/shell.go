package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type BashShell interface {
	Start() error
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
}

type ShellGetter interface {
	// Used to check the real (RealShellGetter) & test (TestShellGetter) structs:
	// GetOffenseInitiatedShell() BashShell // Return RealShell OR TestShell, blueprinted against BashShell interface:
	GetConnectBackInitiatedShell(*Peer) BashShell
}

type RealShellGetter struct {
	// Empty
}

// Holds the REAL shell process/ Cmd struct (from exec pkg):
type RealShell struct {
	realShell exec.Cmd
}

// Get shell for Offense-initiated peer:
// func (g RealShellGetter) GetOffenseInitiatedShell() BashShell {
//
// }

// Get shell for CB-initiated peer:
func (g RealShellGetter) GetConnectBackInitiatedShell(thisPeer *Peer) BashShell {
	// If bash exists on the system, find it, save the path:
	var bShell RealShell

	bashCopPath, err := exec.LookPath(`/bin/bash`) // bashPath @0xfaraday
	if err != nil {
		fmt.Printf("Error finding bash shell path: %v\n", err)
		os.Stderr.WriteString(" " + err.Error() + "\n")
		os.Exit(1)
	}

	bShell = RealShell{realShell: *exec.Command(bashCopPath)}

	// If bash exists, attach the ADDRESS to exec.Cmd to the peer struct:
	//thisPeer.ShellProcess = bShell

	//	// Establish pipe to bash stdin, stdout, & stderr:
	//	bashIn, eRr := bShell.StdinPipe()
	//	if eRr != nil {
	//		fmt.Printf("Error creating shell STDIN pipe: %v\n", eRr)
	//		os.Stderr.WriteString(" " + eRr.Error() + "\n")
	//		os.Exit(1)
	//	}

	//	bashOut, erro := bShell.StdoutPipe()
	//	if erro != nil {
	//		fmt.Printf("Error creating shell STDOUT pipe: %v\n", erro)
	//		os.Stderr.WriteString(" " + erro.Error() + "\n")
	//		os.Exit(1)
	//	}

	// 	bashErr, eRro := bShell.StderrPipe()
	// 	if eRro != nil {
	// 		fmt.Printf("Error creating shell STDERR pipe: %v\n", eRro)
	// 		os.Stderr.WriteString(" " + eRro.Error() + "\n")
	// 		os.Exit(1)
	// 	}

	// Start the shell:
	// var erR error = thisPeer.ShellProcess.Start()
	// if erR != nil {
	// 	fmt.Printf("Error starting shell process: %v\n", erR)
	// 	os.Stderr.WriteString(" " + erR.Error() + "\n")
	// 	os.Exit(1)
	// }

	// Test pipe into stdin
	// 	bashIn, err := bShell.realShell.StdinPipe()
	// 	go func() {
	// 		defer bashIn.Close()
	// 		io.WriteString(bashIn, "/usr/bin/date")
	// 	}()
	//
	// 	out, eRR := bShell.realShell.StdoutPipe()
	// 	if eRR != nil {
	// 		fmt.Printf("Error Getting shell combined output: %v\n", eRR)
	// 		os.Stderr.WriteString(" " + eRR.Error() + "\n")
	// 		os.Exit(1)
	// 	}
	// 	fmt.Printf("Test stdout: %v\n", out)

	// return bShell
	return bShell
}

func (s RealShell) Start() error { // Start the shell, return the error if there is one?

	// Start the shell:
	var erR error = s.realShell.Start()
	if erR != nil {
		fmt.Printf("Error starting shell process: %v\n", erR)
		os.Stderr.WriteString(" " + erR.Error() + "\n")
		os.Exit(1)
	}

	return erR
}

func (s RealShell) StdinPipe() (io.WriteCloser, error) {
	// Establish pipe to bash stdin:
	bashIn, eRr := s.realShell.StdinPipe()
	if eRr != nil {
		fmt.Printf("Error creating shell STDIN pipe: %v\n", eRr)
		os.Stderr.WriteString(" " + eRr.Error() + "\n")
		os.Exit(1)
	}

	return bashIn, eRr
}

func (s RealShell) StdoutPipe() (io.ReadCloser, error) {
	// Establish pipe to bash stdout:
	bashOut, erro := s.realShell.StdoutPipe()
	if erro != nil {
		fmt.Printf("Error creating shell STDOUT pipe: %v\n", erro)
		os.Stderr.WriteString(" " + erro.Error() + "\n")
		os.Exit(1)
	}

	return bashOut, erro
}

func (s RealShell) StderrPipe() (io.ReadCloser, error) {
	// Establish pipe to bash stderr:
	bashErr, eRro := s.realShell.StderrPipe()
	if eRro != nil {
		fmt.Printf("Error creating shell STDERR pipe: %v\n", eRro)
		os.Stderr.WriteString(" " + eRro.Error() + "\n")
		os.Exit(1)
	}

	return bashErr, eRro
}

// Handle own errors, return shell instance blueprinted by BashShell interface:
// func StartAndReturnHelperShell(thisPeer *Peer) BashShell {
// 	// Which peer are we? -- changes shell execution
// 	var g ShellGetter
// 	var s BashShell

// 	if thisPeer.ConnectionType == "connect-back" {
// 		s = g.GetConnectBackInitiatedShell(thisPeer)
// 	}

// 	return s
// }

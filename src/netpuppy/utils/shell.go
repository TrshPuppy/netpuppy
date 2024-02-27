package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type BashShell interface {
	StartShell() error
	PipeStdin() *io.WriteCloser
	PipeStdout() *io.ReadCloser
	PipeStderr() *io.ReadCloser
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
	realShell *exec.Cmd
	// stdinPipeAddress  *io.WriteCloser
	// stdoutPipeAddress *io.ReadCloser
	// stderrPipeAddress *io.ReadCloser
}

// Get shell for Offense-initiated peer:
// func (g RealShellGetter) GetOffenseInitiatedShell() BashShell {
//
// }

// Get shell for CB-initiated peer:
func (g RealShellGetter) GetConnectBackInitiatedShell() BashShell {
	// If bash exists on the system, find it, save the path:
	var bShell RealShell
	var pointerToShell *RealShell

	bashCopPath, err := exec.LookPath(`/usr/bin/bash`) // bashPath @0xfaraday
	if err != nil {
		fmt.Printf("Error finding bash shell path: %v\n", err)
		os.Stderr.WriteString(" " + err.Error() + "\n")
		os.Exit(1)
	}

	bShell = RealShell{realShell: exec.Command(bashCopPath)}
	pointerToShell = &bShell

	fmt.Printf("Address of shell process in shell.go: %v\n", pointerToShell)
	return pointerToShell
}

func (s *RealShell) StartShell() error { // Start the shell, return the error if there is one?
	// Start the shell:
	var erR error = s.realShell.Start()
	// if erR != nil {
	// 	fmt.Printf("Error starting shell process: %v\n", erR)
	// 	os.Stderr.WriteString(" " + erR.Error() + "\n")
	// 	os.Exit(1)
	// }

	return erR
}

func (s *RealShell) PipeStdin() *io.WriteCloser {
	// Establish pipe to bash stdin:
	bashIn, eRr := s.realShell.StdinPipe()
	if eRr != nil {
		fmt.Printf("Error creating shell STDIN pipe: %v\n", eRr)
		os.Stderr.WriteString(" " + eRr.Error() + "\n")
		os.Exit(1)
	}

	defer bashIn.Close()

	// Attach memory address of pipe to struct:
	var pointerToBashIn *io.WriteCloser = &bashIn
	// s.stdinPipeAddress = pointerToBashIn

	fmt.Printf("Address of stdin pipe in shell.go: %v\n", &bashIn)
	return pointerToBashIn
}

func (s *RealShell) PipeStdout() *io.ReadCloser {
	// Establish pipe to bash stdout:
	bashOut, erro := s.realShell.StdoutPipe()
	if erro != nil {
		fmt.Printf("Error creating shell STDOUT pipe: %v\n", erro)
		os.Stderr.WriteString(" " + erro.Error() + "\n")
		os.Exit(1)
	}

	// Attach memory addresses of pipe to struct:
	var pointerToBashOut *io.ReadCloser = &bashOut
	//s.stdoutPipeAddress = pointerToBashOut

	return pointerToBashOut
}

func (s *RealShell) PipeStderr() *io.ReadCloser {
	// Establish pipe to bash stderr:
	bashErr, eRro := s.realShell.StderrPipe()
	if eRro != nil {
		fmt.Printf("Error creating shell STDERR pipe: %v\n", eRro)
		os.Stderr.WriteString(" " + eRro.Error() + "\n")
		os.Exit(1)
	}

	// Attach memory address of pipe to struct:
	var pointerToBashErr *io.ReadCloser = &bashErr
	//	s.stderrPipeAddress = pointerToBashErr

	return pointerToBashErr
}

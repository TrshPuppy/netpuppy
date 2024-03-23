/*
	This code is intentionally redundant.
	I'm trying to re-implement a similar logic as connection.go
	so I can do unit testing WITHOUT creating a real shell process.
	This is unfinished right now, but expect test methods/ structs
	which match the ones implemented below in a fasion similar to
	connection.go.
*/

package shell

import (
	"io"
	"os/exec"
)

// Interface used to blueprint the RealShell struct & eventually TestShell struct:
type ShellInterface interface {
	StartShell() error
	GetStdoutReader() (*io.ReadCloser, error)
	GetStderrReader() (*io.ReadCloser, error)
	GetStdinWriter() (*io.WriteCloser, error)
}

type ShellGetter interface {
	// Used to check the real (RealShellGetter) & test (TestShellGetter) structs:
	// GetOffenseInitiatedShell() ShellInterface // Return RealShell OR TestShell, blueprinted against BashShell interface: <-------- eventually this will be a thing
	GetConnectBackInitiatedShell() (ShellInterface, error)
}

type RealShellGetter struct {
	// Leave empty
}

// Holds the REAL shell process/ Cmd struct (from exec pkg):
type RealShell struct {
	Shell *exec.Cmd
}

// Get shell for Offense-initiated peer:
// func (g RealShellGetter) GetOffenseInitiatedShell() ShellInterface {
//
// }

// Get shell for CB-initiated peer:
func (g RealShellGetter) GetConnectBackInitiatedShell() (ShellInterface, error) {
	// If bash exists on the system, find it, save the path:
	var pointerToShell *RealShell

	bashPath, err := exec.LookPath(`/bin/bash`) // bashPath @0xfaraday
	if err != nil {
		return nil, err
	}

	// Initiate bShell with the struct & process created by exec.Command:
	pointerToShell = &RealShell{Shell: exec.Command(bashPath, "--noprofile", "--norc", "-i", "-s")}

	cmd := *pointerToShell
	cmd.Shell.Env = append(cmd.Shell.Environ(), "PS1=tiddies")
	// Get the pointer to the shell process and & return it:
	return pointerToShell, nil
}

func (s *RealShell) GetStdoutReader() (*io.ReadCloser, error) {
	readCloser, err := s.Shell.StdoutPipe()
	return &readCloser, err
}

func (s *RealShell) GetStderrReader() (*io.ReadCloser, error) {
	readCloser, err := s.Shell.StderrPipe()
	return &readCloser, err
}

func (s *RealShell) GetStdinWriter() (*io.WriteCloser, error) {
	writeCloser, err := s.Shell.StdinPipe()
	return &writeCloser, err
}

// This essentially wraps the actual exec.Cmd.Start() method:
func (s *RealShell) StartShell() error {
	// Start the shell:
	var erR error = s.Shell.Start()
	return erR
}

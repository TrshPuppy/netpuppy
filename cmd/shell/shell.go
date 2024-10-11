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
	"syscall"
)

// REAL code:
type RealShellGetter struct {
	// Leave empty
}

// Holds the REAL shell process/ Cmd struct (from exec pkg):
type RealShell struct {
	Shell *exec.Cmd
}

/*....... TO DO .......
// Get shell for Offense-initiated peer:
// func (g RealShellGetter) GetOffenseInitiatedShell() ShellInterface {
//
// }
.................... */

// Get shell for CB-initiated peer:
func (g RealShellGetter) GetConnectBackInitiatedShell() (*RealShell, error) {
	// If bash exists on the system, find it, save the path:
	var pointerToShell *RealShell

	bashPath, err := exec.LookPath(`/bin/bash`) // bashPath @0xfaraday
	if err != nil {
		return nil, err
	}

	// Initiate bShell with the struct & process created by exec.Command:
	pointerToShell = &RealShell{Shell: exec.Command(bashPath, "--norc", "-s")}
	prompt := `PS1=\[\e]0;\u@\h: \w\a\]\[\033[;94m\]┌──${debian_chroot:+($debian_chroot)──}${VIRTUAL_ENV:+(\[\033[0;1m\]$(basename $VIRTUAL_ENV)\[\033[;94m\])}(\[\033[1;31m\]\u㉿\h\[\033[;94m\])-[\[\033[0;1m\]\w\[\033[;94m\]]\n\[\033[;94m\]└─\[\033[1;31m\]\$\[\033[0m\]`

	// Get cmd and set prompt and the setsid properties on it:
	cmd := pointerToShell
	cmd.Shell.Env = append(cmd.Shell.Environ(), prompt)
	cmd.Shell.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

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

// TEST code:
// type TestShellGetter struct {
// 	//
// }

// type TestShell struct {
// 	Path   string
// 	Stdin  io.Reader
// 	Stdout io.Writer
// 	Stderr io.Writer
// }

// func (g TestShellGetter) GetConnectBackInitiatedShell() (ShellInterface, error) {
// 	var fakePath string = "/bin/tiddies"
// 	var fakeStdin io.Reader
// 	var fakeStdout io.Writer
// 	var fakeStderr io.Writer
// 	var fakeErr error

// 	var fakeShell TestShell = TestShell{Path: fakePath, Stdin: fakeStdin, Stdout: fakeStdout, Stderr: fakeStderr}
// 	var fakePointerToShell *TestShell = &fakeShell

// 	return fakePointerToShell, fakeErr
// }

// func (s *TestShell) GetStdoutReader() (*io.ReadCloser, error) {
// 	var fakeReadCloser io.ReadCloser
// 	var fakeError error

// 	return &fakeReadCloser, fakeError
// }

// func (s *TestShell) GetStderrReader() (*io.ReadCloser, error) {
// 	var fakeReadCloser io.ReadCloser
// 	var fakeError error

// 	return &fakeReadCloser, fakeError
// }

// func (s *TestShell) GetStdinWriter() (*io.WriteCloser, error) {
// 	var fakeWriteCloser io.WriteCloser
// 	var fakeError error

// 	return &fakeWriteCloser, fakeError
// }

// func (s *TestShell) StartShell() error {
// 	var fakeErr error
// 	fakeErr = nil

// 	return fakeErr
// }

package ioctl

import (
	"syscall"
	"unsafe"
)

// Entrance function (from main.go), returns the original terminal attributes or error:
func EnableRawMode(stdinFd int) (*syscall.Termios, syscall.Errno) {
	// Get attributes of current terminal and save them (we have to re-implement them later):
	currentTerm, errno := tcGetAttr(stdinFd)
	if errno != 0 {
		return nil, errno
	}

	// Once we get the current terminal attributes & save them for later,
	//... set the Termios flags for raw mode, then give the new termios to tcsetattr() wrapper
	ogTerm := *currentTerm

	currentTerm.Iflag &^= syscall.IXON | syscall.ICRNL                 // input processing
	currentTerm.Oflag &^= syscall.OPOST                                // output processing
	currentTerm.Lflag &^= syscall.ICANON | syscall.ECHO | syscall.ISIG // disable canonical mode, echo, and signals

	// Now that we've set the properties we want, pass the altered termios structure to the
	//... tcsetattr() wrapper:
	errno = tcSetAttr(stdinFd, currentTerm)
	if errno != 0 {
		return nil, errno
	}

	return &ogTerm, errno
}

// Used to set the terminal back to its OG state:
func DisableRawMode(stdinFd int, termios *syscall.Termios) syscall.Errno {
	// Use the tcsetattr() wrapper to reset the terminal to it's original state
	//... (using the attributes we saved in EnableRawMode)
	errno := tcSetAttr(stdinFd, termios)

	return errno
}

// Using Linux tcgetattr() to get the terminal attributes:
func tcGetAttr(stdinFd int) (*syscall.Termios, syscall.Errno) {
	var currentTerm syscall.Termios

	// Use IOCTL syscall w/ TCGETS op code to get current terminal attributes:
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(stdinFd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&currentTerm)), 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}

	return &currentTerm, errno
}

// Using Linux tcsetattr() to set specific attributes to the termios struct:
func tcSetAttr(stdinFd int, termios *syscall.Termios) syscall.Errno {
	// Use tcsetattr() wrapper to set the terminal to its original attributes
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(stdinFd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(termios)), 0, 0, 0)

	return errno
}

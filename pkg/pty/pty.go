package pty

import (
	"fmt"
	"os"

	"github.com/creack/pty"
)

func GetPseudoterminalDevices() (*os.File, *os.File, error) {
	return Start()
}

func Start() (*os.File, *os.File, error) {
	// Open a new pty
	mptr, sptr, err := pty.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open pty: %w", err)
	}

	return mptr, sptr, nil
}

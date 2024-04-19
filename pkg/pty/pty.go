package pty

import (
	"fmt"
	"os"
	"os/exec"
)

// func Start(c *exec.Cmd) error {
func GetPseudoterminalDevices() (*os.File, *os.File, error) {
	c := exec.Command("/bin/bash")

	// Given the *exec.Cmd, we should start and return the PTYs and PTYmx:
	mDevice, sDevice, err := Start(c)

	return mDevice, sDevice, err
}

func Start(c *exec.Cmd) (*os.File, *os.File, error) {
	var mptr *os.File
	var sname string
	var err error

	// Get pseudoterminal master from /dev/ptmx
	mptr, err = os.OpenFile("/dev/ptmx", os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, nil, err
	}

	//defer mptr.Close()

	// Get the name of the slave device:
	sname, err = GetPTSName(mptr)
	if err != nil {
		return nil, nil, err
	}

	// Now that we have the name, we have to call grantpt() & unlockpt():
	err = GrantPT(mptr)
	if err != nil {
		return nil, nil, err
	}

	err = UnlockPt(mptr)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("PTS name: %s\n", sname)

	// Now that permission is granted, and the slave is unlocked, we can open the pts device file:
	sptr, err := os.OpenFile(sname, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, nil, err
	}
	return mptr, sptr, nil
}

/*
	handle all the pty things
	- start /dev/ptmx
	- return slave
	- return master
	- handle subprocesses spawned from original bash process

	look into:
		- setting stdin/stdout/stderr
		- handle sizing
			- sizes dynamically?
				- ex: were on the other end of a socket




	two struct/ class things
		- slave
		- master
*/

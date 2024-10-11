package pty

import (
	"os"
	"syscall"
)

func GetPseudoterminalDevices() (*os.File, *os.File, error) {
	mDevice, sDevice, err := Start()

	return mDevice, sDevice, err
}

func Start() (*os.File, *os.File, error) {
	var mptr *os.File
	var sname string
	var err error

	// Get pseudoterminal master from /dev/ptmx & set to non-blocking:
	mptr, err = os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NONBLOCK, os.ModeDevice)
	if err != nil {
		return nil, nil, err
	}

	// Get the name of the slave device:
	sname, err = GetPTSName(mptr)
	if err != nil {
		return mptr, nil, err
	}

	// Now that we have the name, we have to call grantpt() & unlockpt():
	err = GrantPT(mptr)
	if err != nil {
		return mptr, nil, err
	}

	err = UnlockPt(mptr)
	if err != nil {
		return mptr, nil, err
	}

	// Now that permission is granted, and the pts is unlocked, we can open the pts device file:
	sptr, err := os.OpenFile(sname, os.O_RDWR|syscall.O_NONBLOCK, os.ModeDevice)
	if err != nil {
		return mptr, nil, err
	}
	return mptr, sptr, nil
}

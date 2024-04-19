package pty

// #define _XOPEN_SOURCE 600
// #include <stdlib.h>
// #include <stdint.h>
//
// int unlock(uint64_t fd) {
//     int unlockStatus = unlockpt(fd);
//     return unlockStatus;
// }
import "C"
import (
	"fmt"
	"os"
)

// Use cgo to call unlockpt() which unlocks our pseudoterminal master/ slave pair:
// ....... https://linux.die.net/man/3/unlockpt
func UnlockPt(masterDevice *os.File) error {
	var err error
	ifd := masterDevice.Fd()

	success := C.unlock(C.ulong(ifd))
	if success != 0 {
		err = fmt.Errorf("Error unlocking pts using unlockpt()")
	}

	return err
}

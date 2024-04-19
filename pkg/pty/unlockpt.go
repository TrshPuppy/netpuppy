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

func UnlockPt(mfd *os.File) error {
	var err error
	ifd := mfd.Fd()

	success := C.unlock(C.ulong(ifd))
	if success != 0 {
		err = fmt.Errorf("Error unlocking pts using unlockpt()")
	}

	return err
}

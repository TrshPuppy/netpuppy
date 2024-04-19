package pty

// #define _XOPEN_SOURCE 600
// #include <stdlib.h>
// #include <stdint.h>
//
// int callGrant(uint64_t fd) {
//     int grantStatus;
//     grantStatus = grantpt(fd);
//
//     return grantStatus;
// }
import "C"
import (
	"fmt"
	"os"
)

func GrantPT(mfd *os.File) error {
	var err error
	ifd := mfd.Fd()

	success := C.callGrant(C.ulong(ifd))
	if success != 0 {
		err = fmt.Errorf("Error granting permissions using grantpt()")
	}

	return err
}

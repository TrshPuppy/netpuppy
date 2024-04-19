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

// Use cgo to call grantpt() which grants our master device access to the slave device:
// ....... https://linux.die.net/man/3/grantpt
func GrantPT(masterDevice *os.File) error {
	var err error
	ifd := masterDevice.Fd()

	success := C.callGrant(C.ulong(ifd))
	if success != 0 {
		err = fmt.Errorf("Error granting permissions using grantpt()")
	}

	return err
}

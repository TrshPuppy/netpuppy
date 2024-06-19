package pty

// #ifdef __linux__
// #define _GNU_SOURCE
// #else
// #define _XOPEN_SOURCE 600
// #endif
// #include <stdlib.h>
// #include <stdint.h>
// #include <stdio.h>
//
// char *getPTSn(uint64_t fd) {
//     size_t len = 1;
//     char *buf;
//     int err;
//     do {
//         buf = malloc(len);
//         if(!buf){
// 	           return NULL;
//         };
//
//         err = ptsname_r(fd, buf, len);
//         if(err != 0) {
//             free(buf);
//             len++;
//             continue;
//         }
//    } while(err);
//
//    return buf;
// }
//
// int freeBuffer(char *buf) {
//     int err = 0;
//     if(buf == NULL) {
//     	err = 1;
//     }
//     free(buf);
//
//     return err;
// }
import "C"
import (
	"fmt"
	"os"
)

// Using CGo to call ptsname_r() which is part of stdlib.h. We need this to get the slave device file name:
// ....... https://linux.die.net/man/3/ptsname
func GetPTSName(mdf *os.File) (string, error) {
	var nullError error
	i := mdf.Fd()
	buf := C.getPTSn(C.ulong(i))
	name := C.GoString(buf)

	null := C.freeBuffer(buf)
	if null == 1 {
		nullError = fmt.Errorf("Error getting PTS name: ptsname_r() unable to find name, name = NULL")
		return name, nullError
	}

	return name, nil
}

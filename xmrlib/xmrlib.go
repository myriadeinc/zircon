package xmrlib

// #cgo CFLAGS: -std=c11 -D_GNU_SOURCE
// #cgo LDFLAGS: -L${SRCDIR} -lxmrlib -Wl,-rpath ${SRCDIR} -lstdc++
// #include <stdlib.h>
// #include <stdint.h>
// #include "src/xmrlib.h"
import (
	"C"
)
import "fmt"

func Hello() bool {
	result := C.check_num((C.uint32_t)(32))
	fmt.Println(result)
	return (bool)(result)
}

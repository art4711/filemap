package mmap

/*
#include <sys/types.h>
#include <sys/mman.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

type Map struct {
	data uintptr
	size uintptr
}

func mmap(len uintptr, prot, flags, fd int, offset uint64) (*Map, error) {
	var v Map

	x, err := C.mmap(unsafe.Pointer(uintptr(0)), C.size_t(len), C.int(prot), C.int(flags), C.int(fd), C.off_t(offset))
	if err != nil {
		return nil, err
	}
	v.data = uintptr(x)
	v.size = len
	return &v, nil
}
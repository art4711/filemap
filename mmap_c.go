package filemap

/*
#include <sys/types.h>
#include <sys/mman.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

const prot_read = C.PROT_READ
const map_shared = C.MAP_SHARED

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

func (m *Map) munmap() {
	C.munmap(unsafe.Pointer(m.data), C.size_t(m.size))
}

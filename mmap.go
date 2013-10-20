package mmap

import (
	"os"
	"syscall"
	"unsafe"
	"errors"
)

func Mapall(file *os.File) (*Map, error) {
	s, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return mmap(uintptr(s.Size()), syscall.PROT_READ, syscall.MAP_SHARED, int(file.Fd()), 0)
}

const maxint = int(^uint(0) >> 1)

func (m Map) GetArray(elem_len uintptr, off, sz uint64, sz_in_bytes bool) (unsafe.Pointer, error) {
	if sz_in_bytes {
		sz /= uint64(elem_len)
	}
	if sz > uint64(maxint) {
		return nil, errors.New("size overflow")
	}
	var sl = struct {
		addr uintptr
		len  int
		cap  int
        }{m.data, int(sz), int(sz)}
	return unsafe.Pointer(&sl), nil
//	b := *(*[]byte)(unsafe.Pointer(&sl))
}

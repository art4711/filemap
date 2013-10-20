package filemap

import (
	"os"
	"syscall"
	"unsafe"
	"errors"
)

type Map struct {
	data uintptr
	size uintptr
}

func New(file *os.File) (*Map, error) {
	s, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return mmap(uintptr(s.Size()), syscall.PROT_READ, syscall.MAP_SHARED, int(file.Fd()), 0)
}

func (m *Map)Close() {
	m.munmap()
}

const maxint = int((^uint(0) >> 1) - 1)
const maxuintptr = ^uintptr(0)

// returns an unsafe.Pointer to a fake slice that it's up to the caller to
// cast to the right type.
func (m Map) Slice(elem_len uintptr, off, sz uint64) (unsafe.Pointer, error) {
	if sz > uint64(maxint) {
		return nil, errors.New("size overflow")
	}
	nptr := uint64(m.data) + off
	if nptr + (sz * uint64(elem_len)) > uint64(maxuintptr) {
		return nil, errors.New("mapping overflow")
	}
	var sl = struct {
		addr uintptr
		len  int
		cap  int
        }{uintptr(nptr), int(sz), int(sz)}
	return unsafe.Pointer(&sl), nil
//	b := *(*[]byte)(unsafe.Pointer(&sl))
}
package filemap

import (
	"os"
	"unsafe"
	"errors"
)

type Map struct {
	data uintptr
	size uintptr
}

// Create a read-only Map of a file.
func NewReader(file *os.File) (*Map, error) {
	s, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return mmap(uintptr(s.Size()), prot_read, map_shared, int(file.Fd()), 0)
}

// Close a Map and delete all mappings.
func (m *Map)Close() {
	m.munmap()
}

const maxint = int((^uint(0) >> 1) - 1)
const maxuintptr = ^uintptr(0)

// Slice creates a slice of a Map doing all necessary error checks to make
// sure that we don't overflow the comically undersized types in a slice.
// returns an unsafe.Pointer to a fake slice that the caller needs to cast to
// a slice of the right type.
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
}

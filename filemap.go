package filemap

import (
	"os"
	"unsafe"
	"errors"
	"reflect"
	"bytes"
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

	var sl reflect.SliceHeader
	sl.Data = uintptr(nptr)
	sl.Len = int(sz)
	sl.Cap = sl.Len
	return unsafe.Pointer(&sl), nil
}

func (m Map) Bytes(off, sz uint64) ([]byte, error) {
	if sz > uint64(maxint) {
		return nil, errors.New("size overflow")
	}
	b, err := m.Slice(1, off, sz)
	if err != nil {
		return nil, err
	}
	return *(*[]byte)(b), nil
}

const string_limit = 64*1024*1024

// CString finds a 
func (m Map) CString(off uint64) ([]byte, error) {
	l := uint64(m.size) - off
	if l > string_limit {
		l = string_limit
	}
	s, err := m.Bytes(off, l)
	if err != nil {
		return nil, err
	}
	e := bytes.IndexByte(s, 0)
	if e == -1 {
		return nil, errors.New("string overflow")
	}
	return s[:e], nil
}

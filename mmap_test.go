package mmap_test

import(
	"mmap"
	"testing"
	"os"	
)

func TestTrivial(t *testing.T) {
	f, err := os.Create("testfile")
	if err != nil {
		t.Fatalf("Create: %v\n", err)
	}
	defer f.Close()
	data := make([]byte, 100)
	for i := 0; i < len(data); i++ {
		data[i] = byte(i)
	}

	f.Write(data)

	m, err := mmap.Mapall(f)
	if err != nil {
		t.Fatalf("Mapall: %v\n", err)
	}

	x, _ := m.GetArray(1, 0, 100, false)
	d := *(*[]byte)(x)

	for i, v := range(d) {
		if v != byte(i) {
			t.Fatalf("%v != %v\n", v, i);
		}
	}
}
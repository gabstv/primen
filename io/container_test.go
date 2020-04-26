package io

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testfs struct {
	okfiles map[string]bool
}

func (fs *testfs) Open(name string) (File, error) {
	if !fs.okfiles[name] {
		return nil, os.ErrNotExist
	}
	return &testfile{
		b: make([]byte, 4096),
	}, nil
}

func (fs *testfs) Stat(name string) (FileInfo, error) {
	if !fs.okfiles[name] {
		return nil, os.ErrNotExist
	}
	return &teststat{
		size: 4096,
	}, nil
}

type testfile struct {
	b []byte
	p int
}

func (f *testfile) Read(p []byte) (n int, err error) {
	if f.p >= len(f.b) {
		return 0, io.EOF
	}
	n = copy(p, f.b[f.p:])
	f.p += n
	return
}

func (f *testfile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.p = int(offset)
		return int64(f.p), nil
	case io.SeekCurrent:
		f.p += int(offset)
		return int64(f.p), nil
	}
	f.p = int(int64(len(f.b)-1) + offset)
	return int64(f.p), nil
}

func (f *testfile) Close() error {
	return nil
}

func (f *testfile) Stat() (FileInfo, error) {
	return &teststat{
		size: int64(len(f.b)),
	}, nil
}

type teststat struct {
	size int64
}

func (s *teststat) Size() int64 {
	return s.size
}

func (s *teststat) IsDir() bool {
	return false
}

func TestContainer(t *testing.T) {
	fsfs := &testfs{
		okfiles: make(map[string]bool),
	}
	fsfs.okfiles["a.txt"] = true
	fsfs.okfiles["b.txt"] = true
	fsfs.okfiles["c.txt"] = true

	c := NewContainer(context.Background(), fsfs)
	ch0 := c.Load("a.txt")
	<-ch0
	fb, err := c.Get("b.txt")
	assert.NoError(t, err)
	assert.Equal(t, 4096, len(fb))
	assert.Equal(t, int64(8192), c.Len())
	c.UnloadAll()
	//
	progch, donech := c.LoadAll([]string{"a.txt", "b.txt", "c.txt"})
	f1 := <-progch
	f2 := <-progch
	f3 := <-progch
	assert.InEpsilon(t, 0.3333333, f1, 0.01)
	assert.InEpsilon(t, 0.6666666, f2, 0.01)
	assert.InEpsilon(t, 0.9999999, f3, 0.01)
	<-donech
}

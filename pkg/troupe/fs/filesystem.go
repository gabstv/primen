package fs

import (
	"bytes"
	"io"
)

// Filesystem interface
type Filesystem interface {
	Open(name string) (io.ReadCloser, error)
	Stat(name string) (Stat, error)
}

// Stat interface
type Stat interface {
	Size() int64
	IsDir() bool
}

func ReadFile(name string, s Filesystem) ([]byte, error) {
	r, err := s.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

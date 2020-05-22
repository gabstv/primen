package io

import (
	"bytes"
	"io"
)

// File interface is almost a copy of http.File
type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Stat() (FileInfo, error)
}

// Filesystem interface
type Filesystem interface {
	Open(name string) (File, error)
	Stat(name string) (FileInfo, error)
}

// FileInfo interface
type FileInfo interface {
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

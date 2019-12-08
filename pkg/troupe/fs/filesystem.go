package fs

import (
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

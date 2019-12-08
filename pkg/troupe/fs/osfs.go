package fs

import (
	"io"
	"os"
	"path"
)

type osfs struct {
	basepath string
}

// NewOS returns a new OS filesystem with the base path as root.
func NewOS(base string) Filesystem {
	return &osfs{base}
}

func (fs *osfs) Open(name string) (io.ReadCloser, error) {
	return os.Open(path.Join(fs.basepath, name))
}

func (fs *osfs) Stat(name string) (Stat, error) {
	finfo, err := os.Stat(path.Join(fs.basepath, name))
	return finfo, err
}

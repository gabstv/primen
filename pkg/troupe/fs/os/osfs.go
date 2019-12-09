package os

import (
	"io"
	"os"
	"path"

	"github.com/gabstv/troupe/pkg/troupe/fs"
)

type osfs struct {
	basepath string
}

// New returns a new OS filesystem with the base path as root.
func New(base string) fs.Filesystem {
	return &osfs{base}
}

func (fs *osfs) Open(name string) (io.ReadCloser, error) {
	return os.Open(path.Join(fs.basepath, name))
}

func (fs *osfs) Stat(name string) (fs.Stat, error) {
	finfo, err := os.Stat(path.Join(fs.basepath, name))
	return finfo, err
}

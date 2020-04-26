package os

import (
	"os"
	"path"

	"github.com/gabstv/tau/io"
)

type osfs struct {
	basepath string
}

type ffile struct {
	*os.File
}

func (f *ffile) Stat() (io.FileInfo, error) {
	fi, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return fi, nil
}

// New returns a new OS filesystem with the base path as root.
func New(base string) io.Filesystem {
	return &osfs{base}
}

func (fs *osfs) Open(name string) (io.File, error) {
	ff, err := os.Open(path.Join(fs.basepath, name))
	if err != nil {
		return nil, err
	}
	return &ffile{
		File: ff,
	}, nil
}

func (fs *osfs) Stat(name string) (io.FileInfo, error) {
	finfo, err := os.Stat(path.Join(fs.basepath, name))
	return finfo, err
}

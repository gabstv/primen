package broccolifs

import (
	"aletheia.icu/broccoli/fs"
	"github.com/gabstv/primen/io"
)

type ffile struct {
	*fs.File
}

func (f *ffile) Stat() (io.FileInfo, error) {
	fi, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return fi, nil
}

type wrapper struct {
	*fs.Broccoli
}

func New(b *fs.Broccoli) io.Filesystem {
	return &wrapper{
		Broccoli: b,
	}
}

func (w *wrapper) Open(name string) (io.File, error) {
	ff, err := w.Broccoli.Open(name)
	if err != nil {
		return nil, err
	}
	return &ffile{
		File: ff,
	}, nil
}

func (w *wrapper) Stat(name string) (io.FileInfo, error) {
	fi, err := w.Broccoli.Stat(name)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

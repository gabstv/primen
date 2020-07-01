package audio

import "bytes"

type Wrapper struct {
	*bytes.Reader
}

func (w *Wrapper) Close() error {
	return nil
}

func NewWrapper(src []byte) *Wrapper {
	return &Wrapper{
		Reader: bytes.NewReader(src),
	}
}

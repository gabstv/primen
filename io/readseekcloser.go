package io

import (
	"bytes"
)

type AudioBuffer struct {
	*bytes.Reader
}

func (b *AudioBuffer) Close() error {
	return nil
}

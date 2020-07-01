package io

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrUnsupportedAudioType Error = "unsupported audio type and/or extension"
)

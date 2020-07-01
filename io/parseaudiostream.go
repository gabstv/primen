package io

import (
	"bytes"
	"path"

	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/audio/wav"

	paudio "github.com/gabstv/primen/audio"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
)

type AudioType int

const (
	AudioTypeWav AudioType = 1
	AudioTypeOgg AudioType = 2
	AudioTypeMP3 AudioType = 3
)

var AudioExtensions = map[string]AudioType{
	".wav": AudioTypeWav,
	".ogg": AudioTypeOgg,
	".mp3": AudioTypeMP3,
}

type AudioStream struct {
	Data   audio.ReadSeekCloser
	Length int64
}

func ParseAudioStream(name string, b []byte) (*AudioStream, error) {
	buf := &AudioBuffer{
		Reader: bytes.NewReader(b),
	}
	ext := path.Ext(name)
	ae := AudioExtensions[ext]
	switch ae {
	case AudioTypeOgg:
		stream, err := vorbis.Decode(paudio.Context(), buf)
		if err != nil {
			return nil, err
		}
		return &AudioStream{
			Data:   stream,
			Length: stream.Length(),
		}, nil
	case AudioTypeMP3:
		stream, err := mp3.Decode(paudio.Context(), buf)
		if err != nil {
			return nil, err
		}
		return &AudioStream{
			Data:   stream,
			Length: stream.Length(),
		}, nil
	case AudioTypeWav:
		stream, err := wav.Decode(paudio.Context(), buf)
		if err != nil {
			return nil, err
		}
		return &AudioStream{
			Data:   stream,
			Length: stream.Length(),
		}, nil
	}
	return nil, ErrUnsupportedAudioType
}

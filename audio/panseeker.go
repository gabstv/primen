package audio

import (
	"bytes"
	"io/ioutil"
	"math"

	"github.com/hajimehoshi/ebiten/audio"
)

const sqrt2div2 = math.Sqrt2 / 2 // math.sqrt(2)/2.0
const rad45 = math.Pi / 4        // 45º

// StereoPanStream is an audio buffer that changes the stereo channel's signal
// based on the Panning.
type StereoPanStream struct {
	audio.ReadSeekCloser
	pan float64 // -1: left; 0: center; 1: right
}

func (s *StereoPanStream) Read(p []byte) (n int, err error) {
	n, err = s.ReadSeekCloser.Read(p)
	if err != nil {
		return
	}
	ls := math.Min(s.pan*-1+1, 1)
	rs := math.Min(s.pan+1, 1)
	for i := 0; i < len(p); i += 4 {
		lc := int16(float64(int16(p[i])|int16(p[i+1])<<8) * ls)
		rc := int16(float64(int16(p[i+2])|int16(p[i+3])<<8) * rs)

		p[i] = byte(lc)
		p[i+1] = byte(lc >> 8)
		p[i+2] = byte(rc)
		p[i+3] = byte(rc >> 8)
	}
	return
}

// func (s *StereoPanStream) Close() error {
// 	return nil
// }

func (s *StereoPanStream) SetPan(pan float64) {
	s.pan = min(max(-1, pan), 1)
}

func (s *StereoPanStream) Pan() float64 {
	return s.pan
}

// NewStereoPanStream returns a new StereoPanStream with a shared buffer src.
// The src's format must be linear PCM (16bits little endian, 2 channel stereo)
// without a header (e.g. RIFF header). The sample rate must be same as that
// of the audio context.
//
// The src can be shared by multiple buffers.
func NewStereoPanStream(src []byte) *StereoPanStream {
	return &StereoPanStream{
		ReadSeekCloser: ioutil.NopCloser(bytes.NewReader(src)).(audio.ReadSeekCloser),
	}
}

// NewStereoPanStreamFromReader returns a new StereoPanStream with buffer src.
//
// The src's format must be linear PCM (16bits little endian, 2 channel stereo)
// without a header (e.g. RIFF header). The sample rate must be same as that
// of the audio context.
func NewStereoPanStreamFromReader(src audio.ReadSeekCloser) *StereoPanStream {
	return &StereoPanStream{
		ReadSeekCloser: src,
	}
}

// test that it fulfills Ebiten's ReadSeekCloser
var _ audio.ReadSeekCloser = &StereoPanStream{}

// MonoPanStream is an audio buffer that changes the stereo channel's signal
// based on the Panning.
type MonoPanStream struct {
	audio.ReadSeekCloser
	pan float64 // -1: left; 0: center; 1: right
}

func (s *MonoPanStream) Read(p []byte) (n int, err error) {
	n, err = s.ReadSeekCloser.Read(p)
	if err != nil {
		return
	}
	angle := rad45 * s.pan
	acos := math.Cos(angle)
	asin := math.Sin(angle)
	ls := sqrt2div2 * (acos - asin)
	rs := sqrt2div2 * (acos + asin)
	for i := 0; i < len(p); i += 4 {
		l := float64(int16(p[i]) | int16(p[i+1])<<8)
		r := float64(int16(p[i+2]) | int16(p[i+3])<<8)
		mixed := (l + r) / 2
		lc := int16(mixed * ls)
		rc := int16(mixed * rs)

		p[i] = byte(lc)
		p[i+1] = byte(lc >> 8)
		p[i+2] = byte(rc)
		p[i+3] = byte(rc >> 8)
	}
	return
}

// func (s *MonoPanStream) Close() error {
// 	return nil
// }

func (s *MonoPanStream) SetPan(pan float64) {
	s.pan = min(max(-1, pan), 1)
}

func (s *MonoPanStream) Pan() float64 {
	return s.pan
}

// NewMonoPanStream returns a new MonoPanStream with a shared buffer src.
// The src's format must be linear PCM (16bits little endian, 2 channel stereo)
// without a header (e.g. RIFF header). The sample rate must be same as that
// of the audio context.
//
// The src can be shared by multiple buffers.
func NewMonoPanStream(src []byte) *MonoPanStream {
	return &MonoPanStream{
		ReadSeekCloser: ioutil.NopCloser(bytes.NewReader(src)).(audio.ReadSeekCloser),
	}
}

// NewMonoPanStreamFromReader returns a new MonoPanStream with buffer src.
//
// The src's format must be linear PCM (16bits little endian, 2 channel stereo)
// without a header (e.g. RIFF header). The sample rate must be same as that
// of the audio context.
func NewMonoPanStreamFromReader(src audio.ReadSeekCloser) *MonoPanStream {
	return &MonoPanStream{
		ReadSeekCloser: src,
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

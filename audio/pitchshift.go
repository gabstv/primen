package audio

import (
	"io"
	"math"

	"github.com/gabstv/ebiten/audio"
)

// http://blogs.zynaptiq.com/bernsee/pitch-shifting-using-the-ft/

// PitchShiftStream is an audio buffer that can dynamically alter the pitch oif the
// original audio as it is being played
type PitchShiftStream struct {
	audio.ReadSeekCloser
	pitchm1 float64 // pitch -1
	fb      []byte  // buffer size * 2
}

// SetPitch sets a pitch between 0.5 and 2
//
// Set the pitch to 1.0 to use the original playback
func (s *PitchShiftStream) SetPitch(pitch float64) {
	s.pitchm1 = math.Max(.5, math.Min(2, pitch)) - 1
}

// Pitch returns the current pitch
func (s *PitchShiftStream) Pitch() float64 {
	return s.pitchm1 + 1
}

//TODO: use FFT for smooth transforms
// https://github.com/takatoh/fft/blob/master/fft.go
// http://blog.bjornroche.com/2012/07/frequency-detection-using-fft-aka-pitch.html
// https://en.wikipedia.org/wiki/Pitch_correction
// https://en.wikipedia.org/wiki/Window_function
// https://en.wikipedia.org/wiki/Pitch_detection_algorithm
// https://stackoverflow.com/questions/8288547/fft-pitch-detection-melody-extraction

func (s *PitchShiftStream) Read(p []byte) (n int, err error) {
	pitch := s.pitchm1 + 1
	if pitch == 1.0 {
		return s.ReadSeekCloser.Read(p)
	}
	if len(s.fb) < len(p)*2 {
		s.fb = make([]byte, len(p)*2)
	}

	lcpos, _ := s.ReadSeekCloser.Seek(0, io.SeekCurrent)
	limit := int(math.Ceil(float64(len(p))*pitch/4) * 4)
	if limit < 0 {
		return s.ReadSeekCloser.Read(p)
	}
	xn, err := s.ReadSeekCloser.Read(s.fb[:limit])
	if err != nil {
		return 0, err
	}
	cursor := 0.0
	targeti := 0
	plen := len(p)
	for {
		i := int(math.Floor(cursor)) * 4
		p[targeti] = s.fb[i]
		p[targeti+1] = s.fb[i+1]
		p[targeti+2] = s.fb[i+2]
		p[targeti+3] = s.fb[i+3]
		cursor += pitch
		targeti += 4
		if targeti >= plen {
			//TODO: debug line below
			s.ReadSeekCloser.Seek(lcpos+int64(i)+4, io.SeekStart)
			return plen, nil
		}
		if int(math.Floor(cursor*4))+4 >= xn {
			//TODO: debug line below
			s.ReadSeekCloser.Seek(lcpos+int64(i)+4, io.SeekStart)
			break
		}
	}
	return targeti - 4, nil
}

// NewPitchShiftStreamFromReader returns a new PitchShiftStream with buffer src.
//
// The src's format must be linear PCM (16bits little endian, 2 channel stereo)
// without a header (e.g. RIFF header). The sample rate must be same as that
// of the audio context.
func NewPitchShiftStreamFromReader(src audio.ReadSeekCloser) *PitchShiftStream {
	return &PitchShiftStream{
		ReadSeekCloser: src,
	}
}

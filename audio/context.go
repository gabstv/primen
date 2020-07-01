package audio

import (
	"sync"

	"github.com/hajimehoshi/ebiten/audio"
)

// SampleRate is the rate (in Hz) at which all the audio will be played.
// Ebiten recommends a sample rate of 44100Hz or 48000Hz. 22050Hz is reported to
// have playback issues on Safari.
type SampleRate int

const (
	Rate44100Hz SampleRate = 44100
	Rate48000Hz SampleRate = 48000
	Rate22050Hz SampleRate = 22050
)

// DefaultRate is the rate that all the audio will pe played on. If you need
// to change it, do it before loading any audio from Primen, as the audio context
// cannot be changed once it is created.
var DefaultRate = Rate44100Hz

var (
	actx *audio.Context
	am   sync.Mutex
)

// Context returns an audio context with the DefaultRate
func Context() *audio.Context {
	am.Lock()
	defer am.Unlock()
	if actx != nil {
		return actx
	}
	actx, _ = audio.NewContext(int(DefaultRate))
	return actx
}

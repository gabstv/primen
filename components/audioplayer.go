package components

import (
	"io"
	"time"

	paudio "github.com/gabstv/primen/audio"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type AudioPlayer struct {
	ebiplayer    *audio.Player
	panctrl      paudio.PanStream
	pitchshifter *paudio.PitchShiftStream
	//channels
}

type NewAudioPlayerInput struct {
	RawAudio      []byte        // use RawAudio (for shared buffers) and sfx
	Buffer        io.ReadSeeker // use Buffer for large files
	Panning       bool          // use audio Panning feature
	StereoPanning bool
	PitchShift    bool
	Infinite      bool
	IntroLength   int64
	LoopLength    int64
}

func mustEbiPlayer(p *audio.Player, e error) *audio.Player {
	if e != nil {
		panic(e)
	}
	return p
}

func NewAudioPlayer(input NewAudioPlayerInput) AudioPlayer {
	if input.Buffer == nil && input.RawAudio == nil {
		panic("input.Buffer and input.RawAudio are nil at the same time")
	}
	var lsrk io.ReadSeeker
	if input.RawAudio != nil {
		lsrk = paudio.NewWrapper(input.RawAudio)
	} else {
		lsrk = input.Buffer
	}
	var pan paudio.PanStream
	if input.Panning {
		if input.StereoPanning {
			pan = paudio.NewStereoPanStreamFromReader(lsrk)
		} else {
			pan = paudio.NewMonoPanStreamFromReader(lsrk)
		}
		lsrk = pan
	}
	var pshift *paudio.PitchShiftStream
	if input.PitchShift {
		pshift = paudio.NewPitchShiftStreamFromReader(lsrk)
		lsrk = pshift
	}
	if input.Infinite {
		lsrk = audio.NewInfiniteLoopWithIntro(lsrk, input.IntroLength, input.LoopLength)
	}
	return AudioPlayer{
		ebiplayer:    mustEbiPlayer(audio.NewPlayer(paudio.Context(), lsrk)),
		panctrl:      pan,
		pitchshifter: pshift,
	}
}

func (p *AudioPlayer) Play() {
	p.ebiplayer.Play()
}

func (p *AudioPlayer) Pause() {
	p.ebiplayer.Pause()
}

func (p *AudioPlayer) Rewind() error {
	return p.ebiplayer.Rewind()
}

func (p *AudioPlayer) Current() time.Duration {
	return p.ebiplayer.Current()
}

func (p *AudioPlayer) IsPlaying() bool {
	return p.ebiplayer.IsPlaying()
}

func (p *AudioPlayer) Seek(offset time.Duration) error {
	return p.ebiplayer.Seek(offset)
}

func (p *AudioPlayer) SetVolume(volume float64) {
	p.ebiplayer.SetVolume(volume)
}

func (p *AudioPlayer) Volume() float64 {
	return p.ebiplayer.Volume()
}

func (p *AudioPlayer) SetPan(pan float64) {
	if p.panctrl != nil {
		p.panctrl.SetPan(pan)
	}
}

func (p *AudioPlayer) Pan() float64 {
	if p.panctrl != nil {
		return p.panctrl.Pan()
	}
	return 0
}

func (p *AudioPlayer) SetPitch(pan float64) {
	if p.pitchshifter != nil {
		p.pitchshifter.SetPitch(pan)
	}
}

func (p *AudioPlayer) Pitch() float64 {
	if p.pitchshifter != nil {
		return p.pitchshifter.Pitch()
	}
	return 1
}

//go:generate ecsgen -n AudioPlayer -p components -o audioplayer_component.go --component-tpl --vars "UUID=9C7DB259-6A3E-4DD3-B277-4B35DA5709AF"

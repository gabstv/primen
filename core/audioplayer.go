package core

import (
	"time"

	paudio "github.com/gabstv/primen/audio"
	"github.com/hajimehoshi/ebiten/audio"
)

type AudioPlayer struct {
	ebiplayer *audio.Player
	panctrl   *paudio.StereoPanStream
	//channels
}

type NewAudioPlayerInput struct {
	RawAudio    []byte               // use RawAudio (for shared buffers) and sfx
	Buffer      audio.ReadSeekCloser // use Buffer for large files
	Panning     bool                 // use audio Panning feature
	Infinite    bool
	IntroLength int64
	LoopLength  int64
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
	var lsrk audio.ReadSeekCloser
	if input.RawAudio != nil {
		lsrk = paudio.NewWrapper(input.RawAudio)
	} else {
		lsrk = input.Buffer
	}
	var pan *paudio.StereoPanStream
	if input.Panning {
		pan = paudio.NewStereoPanStreamFromReader(lsrk)
		lsrk = pan
	}
	if input.Infinite {
		lsrk = audio.NewInfiniteLoopWithIntro(lsrk, input.IntroLength, input.LoopLength)
	}
	return AudioPlayer{
		ebiplayer: mustEbiPlayer(audio.NewPlayer(paudio.Context(), lsrk)),
		panctrl:   pan,
	}
}

func (p *AudioPlayer) Play() {
	_ = p.ebiplayer.Play()
}

func (p *AudioPlayer) Pause() {
	_ = p.ebiplayer.Pause()
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

//go:generate ecsgen -n AudioPlayer -p core -o audioplayer_component.go --component-tpl --vars "UUID=9C7DB259-6A3E-4DD3-B277-4B35DA5709AF"

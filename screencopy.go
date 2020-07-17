package primen

import (
	"sync"

	"github.com/hajimehoshi/ebiten"
)

type ScreenCopyRequest interface {
	Done() <-chan struct{}
	ScreenCopy() *ebiten.Image
}

type screenCopyRequest struct {
	sync.Mutex
	ch  chan struct{}
	img *ebiten.Image
}

func (r *screenCopyRequest) Done() <-chan struct{} {
	return r.ch
}

func (r *screenCopyRequest) ScreenCopy() *ebiten.Image {
	r.Lock()
	defer r.Unlock()
	return r.img
}

func (e *engine) WaitAndGrabScreenImage() ScreenCopyRequest {
	i := &screenCopyRequest{
		ch: make(chan struct{}),
	}
	e.screencopych <- i
	return i
}

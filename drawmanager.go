package primen

import (
	"sort"

	"github.com/gabstv/ebiten"
	"github.com/gabstv/primen/core"
)

type EngineDrawManager interface {
	core.DrawManager
	DrawTargets()
	PrepareTargets()
}

type soloDrawManager struct {
	screen *ebiten.Image
}

var _ EngineDrawManager = (*soloDrawManager)(nil)

func (m *soloDrawManager) DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask core.DrawMask) {
	if mask == 0 {
		return
	}
	m.screen.DrawImage(image, opt)
}

func (m *soloDrawManager) Screen() *ebiten.Image {
	return m.screen
}

func (m *soloDrawManager) PrepareTargets() {
	// solo draw manager doesn't have draw targets
}

func (m *soloDrawManager) DrawTargets() {
	// solo draw manager doesn't have draw targets
}

func (m *soloDrawManager) DrawTarget(id core.DrawTargetID) core.DrawTarget {
	return nil
}

type dtDrawManager struct {
	screen *ebiten.Image
	dts    []EngineDrawTarget
}

var _ EngineDrawManager = (*dtDrawManager)(nil)

func (m *dtDrawManager) DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask core.DrawMask) {
	if mask == 0 {
		return
	}
	for _, dt := range m.dts {
		dt.DrawImage(image, opt, mask)
	}
	//m.screen.DrawImage(image, opt)
}

func (m *dtDrawManager) Screen() *ebiten.Image {
	return m.screen
}

func (m *dtDrawManager) PrepareTargets() {
	for _, dt := range m.dts {
		dt.PrepareFrame(m.screen)
	}
}

func (m *dtDrawManager) DrawTargets() {
	for _, dt := range m.dts {
		dt.DrawFrame(m.screen)
	}
}

func (m *dtDrawManager) DrawTarget(id core.DrawTargetID) core.DrawTarget {
	i := sort.Search(len(m.dts), func(i int) bool {
		return m.dts[i].ID() >= id
	})
	if i < len(m.dts) && m.dts[i].ID() == id {
		return m.dts[i]
	}
	return nil
}

func (e *engine) newDrawManager(screen *ebiten.Image) EngineDrawManager {
	e.drawTargetLock.Lock()
	l := len(e.drawTargets)
	e.drawTargetLock.Unlock()
	if l < 1 {
		return &soloDrawManager{
			screen: screen,
		}
	}
	e.drawTargetLock.Lock()
	defer e.drawTargetLock.Unlock()
	cp := make([]EngineDrawTarget, len(e.drawTargets))
	copy(cp, e.drawTargets)
	return &dtDrawManager{
		screen: screen,
		dts:    cp,
	}
}

package primen

import (
	"sort"

	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/geom"
	"github.com/hajimehoshi/ebiten"
)

type EngineDrawTarget interface {
	core.DrawTarget
	PrepareFrame(screen *ebiten.Image)
	DrawFrame(screen *ebiten.Image)
}

/////
type baseDrawTarget struct {
	id   core.DrawTargetID
	mask core.DrawMask
	m    ebiten.GeoM
	mset bool
	size geom.Vec
}

func (d *baseDrawTarget) ID() core.DrawTargetID {
	return d.id
}

func (d *baseDrawTarget) DrawMask() core.DrawMask {
	return d.mask
}

func (d *baseDrawTarget) Translate(v geom.Vec) {
	m := d.m
	m.Translate(v.X, v.Y)
	d.m = m
	d.mset = true
}

func (d *baseDrawTarget) Scale(v geom.Vec) {
	m := d.m
	m.Scale(v.X, v.Y)
	d.m = m
	d.mset = true
}

func (d *baseDrawTarget) Rotate(rad float64) {
	m := d.m
	m.Rotate(rad)
	d.m = m
	d.mset = true
}

func (d *baseDrawTarget) ResetTransform() {
	d.mset = false
	d.m = ebiten.GeoM{}
}

func (d *baseDrawTarget) Size() geom.Vec {
	return d.size
}

// drawImage draws image into dst (if mask test succeeds)
func (d *baseDrawTarget) drawImage(dst, image *ebiten.Image, opt *ebiten.DrawImageOptions, mask core.DrawMask) {
	if (mask & d.mask) != d.mask {
		return
	}
	if d.mset {
		// m := d.m
		// pm := opt.GeoM
		// m.Concat(pm)
		// opt.GeoM = m
		// _ = dst.DrawImage(image, opt)
		// opt.GeoM = pm
		m := opt.GeoM
		pm := opt.GeoM
		m.Concat(d.m)
		opt.GeoM = m
		_ = dst.DrawImage(image, opt)
		opt.GeoM = pm
		return
	}
	dst.DrawImage(image, opt)
}

/////

type drawTarget struct {
	*baseDrawTarget
	image  *ebiten.Image
	filter ebiten.Filter
	bounds geom.Rect
}

var _ EngineDrawTarget = (*drawTarget)(nil)

func (d *drawTarget) DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask core.DrawMask) {
	d.drawImage(d.image, image, opt, mask)
}

func (d *drawTarget) Image() *ebiten.Image {
	return d.image
}

func (d *drawTarget) DrawFrame(screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	if d.bounds.IsZero() {
		_ = screen.DrawImage(d.image, opt)
		return
	}
	w, h := screen.Size()
	opt.GeoM.Translate(float64(w)*d.bounds.Min.X, float64(h)*d.bounds.Min.Y)
	_ = screen.DrawImage(d.image, opt)
}

func (d *drawTarget) PrepareFrame(screen *ebiten.Image) {
	d.setSize(screen)
	d.image.Clear()
}

func (d *drawTarget) Size() geom.Vec {
	return d.size
}

func (d *drawTarget) setSize(screen *ebiten.Image) {
	w, h := screen.Size()
	tsize := geom.Vec{float64(w), float64(h)}
	if !d.bounds.IsZero() {
		tsize = tsize.Mul(d.bounds.Size())
	}
	if !tsize.EqualsEpsilon(d.size) {
		// create image
		if d.image != nil {
			d.image.Dispose()
			d.image = nil
		}
		d.image, _ = ebiten.NewImage(int(tsize.X), int(tsize.Y), d.filter)
		d.size = tsize
	}
}

//

type screenDrawTarget struct {
	*baseDrawTarget
	screen *ebiten.Image
}

var _ EngineDrawTarget = (*screenDrawTarget)(nil)

func (d *screenDrawTarget) DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask core.DrawMask) {
	d.drawImage(d.screen, image, opt, mask)
}

func (d *screenDrawTarget) Image() *ebiten.Image {
	return d.screen
}

func (d *screenDrawTarget) DrawFrame(screen *ebiten.Image) {
	// noop because this drawTarget is the screen
}

func (d *screenDrawTarget) PrepareFrame(screen *ebiten.Image) {
	d.screen = screen
}

func (d *screenDrawTarget) Size() geom.Vec {
	return d.size
}

//

type programmableDrawTarget struct {
	*baseDrawTarget
	drawImageFn    func(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask core.DrawMask, camG ebiten.GeoM)
	drawFrameFn    func(screen *ebiten.Image)
	prepareFrameFn func(screen *ebiten.Image)
	imageFn        func() *ebiten.Image
	sizeFn         func() geom.Vec
}

var _ EngineDrawTarget = (*programmableDrawTarget)(nil)

func (d *programmableDrawTarget) DrawImage(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask core.DrawMask) {
	d.drawImageFn(image, opt, mask, d.m)
}

func (d *programmableDrawTarget) Image() *ebiten.Image {
	return d.imageFn()
}

func (d *programmableDrawTarget) DrawFrame(screen *ebiten.Image) {
	d.drawFrameFn(screen)
}

func (d *programmableDrawTarget) PrepareFrame(screen *ebiten.Image) {
	d.prepareFrameFn(screen)
}

func (d *programmableDrawTarget) Size() geom.Vec {
	return d.sizeFn()
}

//

func (e *engine) NewDrawTarget(mask core.DrawMask, bounds geom.Rect, filter ebiten.Filter) core.DrawTargetID {
	if bounds.Min.X > bounds.Max.X || bounds.Min.Y > bounds.Min.Y {
		panic("invalid drawTarget bounds")
	}
	e.drawTargetLock.Lock()
	defer e.drawTargetLock.Unlock()
	e.lastDrawTargetID++
	id := e.lastDrawTargetID
	t := &drawTarget{
		baseDrawTarget: &baseDrawTarget{
			id:   id,
			mask: mask,
		},
		bounds: bounds,
		filter: filter,
	}
	e.drawTargets = append(e.drawTargets, t)
	return id
}

func (e *engine) NewScreenOffsetDrawTarget(mask core.DrawMask) core.DrawTargetID {
	e.drawTargetLock.Lock()
	defer e.drawTargetLock.Unlock()
	e.lastDrawTargetID++
	id := e.lastDrawTargetID
	t := &screenDrawTarget{
		baseDrawTarget: &baseDrawTarget{
			id:   id,
			mask: mask,
			size: e.SizeVec(),
		},
	}
	e.drawTargets = append(e.drawTargets, t)
	return id
}

type ProgrammableDrawTargetInput struct {
	DrawImage    func(image *ebiten.Image, opt *ebiten.DrawImageOptions, mask core.DrawMask, camG ebiten.GeoM)
	DrawFrame    func(screen *ebiten.Image)
	PrepareFrame func(screen *ebiten.Image)
	Image        func() *ebiten.Image
	Size         func() geom.Vec
}

func (e *engine) NewProgrammableDrawTarget(input ProgrammableDrawTargetInput) core.DrawTargetID {
	e.drawTargetLock.Lock()
	defer e.drawTargetLock.Unlock()
	e.lastDrawTargetID++
	id := e.lastDrawTargetID
	t := &programmableDrawTarget{
		baseDrawTarget: &baseDrawTarget{
			id: id,
		},
		drawFrameFn:    input.DrawFrame,
		drawImageFn:    input.DrawImage,
		imageFn:        input.Image,
		prepareFrameFn: input.PrepareFrame,
		sizeFn:         input.Size,
	}
	e.drawTargets = append(e.drawTargets, t)
	return id
}

func (e *engine) DrawTarget(id core.DrawTargetID) core.DrawTarget {
	e.drawTargetLock.Lock()
	defer e.drawTargetLock.Unlock()
	i := sort.Search(len(e.drawTargets), func(i int) bool {
		return e.drawTargets[i].ID() >= id
	})
	if i < len(e.drawTargets) && e.drawTargets[i].ID() == id {
		return e.drawTargets[i]
	}
	return nil
}

func (e *engine) RemoveDrawTarget(id core.DrawTargetID) bool {
	e.drawTargetLock.Lock()
	defer e.drawTargetLock.Unlock()
	i := sort.Search(len(e.drawTargets), func(i int) bool {
		return e.drawTargets[i].ID() >= id
	})
	if i < len(e.drawTargets) && e.drawTargets[i].ID() == id {
		e.drawTargets = e.drawTargets[:i+copy(e.drawTargets[i:], e.drawTargets[i+1:])]
		return true
	}
	return false
}

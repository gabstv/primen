package core

import (
	"image/color"
	"math/rand"

	"github.com/gabstv/ecs/v2"
	"github.com/hajimehoshi/ebiten"
)

type SpawnStrategy int

const (
	SpawnReplace SpawnStrategy = iota
	SpawnPause   SpawnStrategy = iota
)

var (
	emitterRNG = rand.New(rand.NewSource(1337))
)

type ParticleEmitter struct {
	// if there is no transform attached, the props
	// below are used

	x      float64 // logical X position
	y      float64 // logical Y position
	angle  float64 // radians
	scaleX float64 // logical X scale (1 = 100%)
	scaleY float64 // logical Y scale (1 = 100%)

	// exclusive props

	max           int // max active particles
	props         ParticleProps
	parenttre     ecs.Entity // this is used to calculate GeoM
	parenttr      *Transform // this is used to calculate GeoM
	parenttrevtid int64      // used to deregister transform destroy event
	particles     []Particle
	strategy      SpawnStrategy
	disabled      bool
	opt           ebiten.DrawImageOptions
	eid           ecs.Entity // cached entity id used to obtains its own transform
	ew            ecs.BaseWorld
	rand          *rand.Rand
}

func (e *ParticleEmitter) SetEmissionParentToSelf() {
	e.SetEmissionParent(e.eid)
}

func (e *ParticleEmitter) SetEmissionParent(en ecs.Entity) {
	if e.parenttrevtid != 0 && e.parenttr != nil && e.parenttre != 0 {
		// unset
		e.parenttr.RemoveDestroyListener(e.parenttrevtid)
		e.parenttrevtid = 0
	}
	e.parenttr = nil
	e.parenttre = 0
	//
	e.parenttre = en
	if tr := e.fetchTransform(); tr != nil {
		ww := e.ew
		e.parenttrevtid = tr.AddDestroyListener(func() {
			if x := GetParticleEmitterComponentData(ww, en); x != nil {
				x.transformWillBeDestroyed()
			}
		})
	} else {
		e.parenttre = 0
	}
}

func (e *ParticleEmitter) transformWillBeDestroyed() {
	e.parenttrevtid = 0
	e.parenttr = nil
	e.parenttre = 0
}

func (e *ParticleEmitter) Emit() bool {
	rng := e.rand
	if rng == nil {
		rng = emitterRNG
	}
	particle := e.props.NewParticle(rng, e)
	if len(e.particles) >= e.max && e.max != 0 {
		switch e.strategy {
		case SpawnReplace:
			last := len(e.particles) - 1
			copy(e.particles, e.particles[1:])
			e.particles[last] = particle
			return true
		case SpawnPause:
			return false
		}
	}
	e.particles = append(e.particles, particle)
	return true
}

func (e *ParticleEmitter) Draw(ctx DrawCtx, d *Drawable) {
	if d == nil {
		return
	}
	if e.disabled {
		return
	}
	//g := d.G(t.scaleX, t.scaleY, t.angle, t.x, t.y)
	//o := &t.opt
	opt := &ebiten.DrawImageOptions{}
	for _, p := range e.particles {
		//
		opt.GeoM.Reset()
		opt.GeoM.Scale(p.sx, p.sy)
		opt.GeoM.Rotate(p.r)
		opt.GeoM.Translate(p.pox+p.px, p.poy+p.py)
		if p.parenttr != nil {
			opt.GeoM.Concat(p.parenttr.m)
		}
		opt.ColorM.Reset()
		opt.ColorM.Scale(p.clr, p.clb, p.clg, p.cla)
		ctx.Renderer().DrawImageRaw(p.img, opt)
	}
}

func (e *ParticleEmitter) fetchTransform() *Transform {
	if e.parenttre != 0 {
		e.parenttr = GetTransformComponentData(e.ew, e.parenttre)
		if e.parenttr == nil {
			e.parenttre = 0
			return nil
		}
		return e.parenttr
	}
	e.parenttr = nil
	return nil
}

//go:generate ecsgen -n ParticleEmitter -p core -o particleemitter_component.go --component-tpl --vars "UUID=19A70DF9-0B1A-4A85-B23E-7BCA8E0857D7" --vars "BeforeRemove=c.beforeRemove(e)" --vars "OnAdd=c.onAdd(e)"

func (c *ParticleEmitterComponent) beforeRemove(e ecs.Entity) {
	if d := GetDrawableComponentData(c.world, e); d != nil {
		d.drawer = nil
	}
	emt := GetParticleEmitterComponentData(c.world, e)
	emt.eid = 0
	if emt.parenttrevtid != 0 && emt.parenttr != nil {
		emt.parenttr.RemoveDestroyListener(emt.parenttrevtid)
		emt.parenttrevtid = 0
		emt.parenttre = 0
		emt.parenttr = nil
	}
}

func (c *ParticleEmitterComponent) onAdd(e ecs.Entity) {
	if d := GetDrawableComponentData(c.world, e); d != nil {
		d.drawer = c.Data(e)
	} else {
		SetDrawableComponentData(c.world, e, Drawable{
			drawer: c.Data(e),
		})
	}
	GetParticleEmitterComponentData(c.world, e).eid = e
}

type ParticleProps struct {
	Px, Py                     float64         // initial position
	Vpx0, Vpy0, Vpx1, Vpy1     float64         // randomized position range
	Vx, Vy                     float64         // initial velocity
	Vvx0, Vvy0, Vvx1, Vvy1     float64         // randomized velocity range
	Ax, Ay                     float64         // initial acceleration
	Vax0, Vay0, Vax1, Vay1     float64         // randomized acceleration range
	R                          float64         // initial rotation (radians)
	Vr0, Vr1                   float64         // randomized initial rotation (radians)
	Rvb                        float64         // initial rotation velocity (radians/second)
	Rve                        float64         // end rotation velocity (radians/second)
	Sx, Sy                     float64         // initial scale
	Vsx0, Vsy0, Vsx1, Vsy1     float64         // randomized initial scale
	Esx, Esy                   float64         // end scale
	Vesx0, Vesy0, Vesx1, Vesy1 float64         // randomized end scale
	Ox, Oy                     float64         // origin
	Dur                        float64         // duration
	Vdur0, Vdur1               float64         // randomized duration (seconds)
	Colorb                     color.RGBA      // initial color tint modifier
	Colore                     color.RGBA      // end color tint modifier
	Source                     []*ebiten.Image // particle source(s); if it's more than one, it will be randomized
}

func (pp ParticleProps) NewParticle(rng *rand.Rand, e *ParticleEmitter) Particle {
	p := Particle{
		bclr: float64(pp.Colorb.R) / 255,
		bclg: float64(pp.Colorb.G) / 255,
		bclb: float64(pp.Colorb.B) / 255,
		bcla: float64(pp.Colorb.A) / 255,
		clr:  float64(pp.Colorb.R) / 255,
		clg:  float64(pp.Colorb.G) / 255,
		clb:  float64(pp.Colorb.B) / 255,
		cla:  float64(pp.Colorb.A) / 255,
		eclr: float64(pp.Colore.R) / 255,
		eclg: float64(pp.Colore.G) / 255,
		eclb: float64(pp.Colore.B) / 255,
		ecla: float64(pp.Colore.A) / 255,
		//img: rng.Intn(len(pp.Source)),
		//pox: calc after obtaining image,
		//poy: calc after obtaining image,
		parenttr:  e.parenttr,
		parenttre: e.parenttre,
		ax:        pp.Ax + Lerpf(pp.Vax0, pp.Vax1, rng.Float64()),
		ay:        pp.Ay + Lerpf(pp.Vay0, pp.Vay1, rng.Float64()),
		dur:       pp.Dur + Lerpf(pp.Vdur0, pp.Vdur1, rng.Float64()),
		px:        pp.Px + Lerpf(pp.Vpx0, pp.Vpx1, rng.Float64()),
		py:        pp.Py + Lerpf(pp.Vpy0, pp.Vpy1, rng.Float64()),
		r:         pp.R + Lerpf(pp.Vr0, pp.Vr1, rng.Float64()),
		vx:        pp.Vx + Lerpf(pp.Vvx0, pp.Vvx1, rng.Float64()),
		vy:        pp.Vy + Lerpf(pp.Vvy0, pp.Vvy1, rng.Float64()),
		//sx: GET FROM bsx,
		//sy: GET FROM bsy,
		bsx: pp.Sx + Lerpf(pp.Vsx0, pp.Vsx1, rng.Float64()),
		bsy: pp.Sy + Lerpf(pp.Vsy0, pp.Vsy1, rng.Float64()),
		esx: pp.Esx + Lerpf(pp.Vesx0, pp.Vesx1, rng.Float64()),
		esy: pp.Esy + Lerpf(pp.Vesy0, pp.Vesy1, rng.Float64()),
		t:   0,
	}
	p.sx = p.bsx
	p.sy = p.bsy
	if len(pp.Source) == 1 {
		p.img = pp.Source[0]
	} else if len(pp.Source) > 1 {
		p.img = pp.Source[rng.Intn(len(pp.Source))]
	}
	if p.img != nil {
		xx, yy := p.img.Size()
		p.pox = applyOrigin(float64(xx), pp.Ox)
		p.poy = applyOrigin(float64(yy), pp.Oy)
	}
	return p
}

type Particle struct {
	px, py                 float64
	vx, vy                 float64
	ax, ay                 float64
	sx, sy                 float64
	bsx, bsy, esx, esy     float64
	r                      float64
	pox, poy               float64 // precomputed origin
	bclr, bclg, bclb, bcla float64 // initial color tint
	eclr, eclg, eclb, ecla float64 // end color tint
	clr, clg, clb, cla     float64 // current color tint
	t, dur                 float64
	parenttre              ecs.Entity
	parenttr               *Transform
	img                    *ebiten.Image
}

//go:generate ecsgen -n DrawableParticleEmitter -p core -o particleemitter_drawablesystem.go --system-tpl --vars "Priority=10" --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "OnResize=s.onResize()" --vars "OnWillResize=s.onWillResize()" --vars "UUID=627C4B36-EE45-40C6-91AE-617D5CFDD8FC" --components "Drawable" --components "ParticleEmitter"

var matchDrawableParticleEmitterSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	return eflag.Contains(GetDrawableComponent(w).Flag().Or(GetParticleEmitterComponent(w).Flag()))
}

var resizematchDrawableParticleEmitterSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	if eflag.Contains(GetDrawableComponent(w).Flag()) {
		return true
	}
	if eflag.Contains(GetParticleEmitterComponent(w).Flag()) {
		return true
	}
	return false
}

func (s *DrawableParticleEmitterSystem) onEntityAdded(e ecs.Entity) {

}

func (s *DrawableParticleEmitterSystem) onEntityRemoved(e ecs.Entity) {
	//if x := GetP
}

func (s *DrawableParticleEmitterSystem) onWillResize() {
	v := s.V()
	for i := range v.entities {
		v.entities[i].ParticleEmitter.parenttr = nil
		v.entities[i].Drawable.drawer = nil
	}
}

func (s *DrawableParticleEmitterSystem) onResize() {
	for _, match := range s.V().Matches() {
		match.ParticleEmitter.fetchTransform()
		match.Drawable.drawer = match.ParticleEmitter
	}
}

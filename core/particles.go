package core

import (
	"image/color"
	"math"
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
	opt           ebiten.DrawImageOptions
	eid           ecs.Entity // cached entity id used to obtains its own transform
	ew            ecs.BaseWorld
	rand          *rand.Rand
	emission      EmissionProp
	emissiont     float64
	emissiontnext float64
	compositeMode ebiten.CompositeMode
	disabled      bool
}

func NewParticleEmitter(w ecs.BaseWorld) ParticleEmitter {
	return ParticleEmitter{
		particles: make([]Particle, 0, 64),
		max:       64,
		scaleX:    1,
		scaleY:    1,
		strategy:  SpawnPause,
		props: ParticleProps{
			Vy:     -100,
			Colorb: color.RGBA{255, 255, 255, 255},
			Colore: color.RGBA{255, 255, 255, 0},
			Vvx0:   -3,
			Vvx1:   3,
			Ox:     .5,
			Oy:     .5,
			Vr0:    -math.Pi,
			Vr1:    math.Pi,
			Dur:    1,
			Sx:     1,
			Sy:     1,
			Esx:    .5,
			Esy:    .5,
		},
		emission: EmissionProp{
			Enabled: true,
			N0:      1,
			N1:      1,
			T0:      .1,
			T1:      .2,
		},
		ew: w,
	}
}

func (e *ParticleEmitter) X() float64 {
	return e.x
}

func (e *ParticleEmitter) SetX(x float64) *ParticleEmitter {
	e.x = x
	return e
}

func (e *ParticleEmitter) Y() float64 {
	return e.y
}

func (e *ParticleEmitter) SetY(y float64) *ParticleEmitter {
	e.y = y
	return e
}

func (e *ParticleEmitter) Angle() float64 {
	return e.angle
}

func (e *ParticleEmitter) SetAngle(angle float64) *ParticleEmitter {
	e.angle = angle
	return e
}

func (e *ParticleEmitter) ScaleX() float64 {
	return e.scaleX
}

func (e *ParticleEmitter) SetScaleX(sx float64) *ParticleEmitter {
	e.scaleX = sx
	return e
}

func (e *ParticleEmitter) ScaleY() float64 {
	return e.scaleY
}

func (e *ParticleEmitter) SetScaleY(sy float64) *ParticleEmitter {
	e.scaleY = sy
	return e
}

func (e *ParticleEmitter) MaxParticles() int {
	return e.max
}

func (e *ParticleEmitter) SetMaxParticles(max int) *ParticleEmitter {
	e.max = max
	return e
}

func (e *ParticleEmitter) SetStrategy(strategy SpawnStrategy) *ParticleEmitter {
	e.strategy = strategy
	return e
}

func (e *ParticleEmitter) Enabled() bool {
	return !e.disabled
}

func (e *ParticleEmitter) SetEnabled(enabled bool) *ParticleEmitter {
	e.disabled = !enabled
	return e
}

func (e *ParticleEmitter) Props() ParticleProps {
	return e.props
}

func (e *ParticleEmitter) SetProps(props ParticleProps) *ParticleEmitter {
	e.props = props
	return e
}

func (e *ParticleEmitter) EmissionProp() EmissionProp {
	return e.emission
}

func (e *ParticleEmitter) SetEmissionProp(prop EmissionProp) *ParticleEmitter {
	e.emission = prop
	return e
}

func (e *ParticleEmitter) CompositeMode() ebiten.CompositeMode {
	return e.compositeMode
}

func (e *ParticleEmitter) SetCompositeMode(m ebiten.CompositeMode) *ParticleEmitter {
	e.compositeMode = m
	return e
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
	if len(e.props.Source) < 1 {
		return false
	}
	rng := e.rand
	if rng == nil {
		rng = emitterRNG
	}
	particle := e.props.NewParticle(rng, e)
	//TODO: link position with transform!
	particle.px += e.x
	particle.py += e.y
	particle.r += e.angle
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
	opt := &ebiten.DrawImageOptions{}
	opt.CompositeMode = e.compositeMode
	for _, p := range e.particles {
		//
		opt.GeoM.Reset()
		opt.GeoM.Translate(p.pox, p.poy)
		opt.GeoM.Scale(p.sx, p.sy)
		opt.GeoM.Rotate(p.r)
		if p.parenttr != nil {
			opt.GeoM.Concat(p.parenttr.m)
		}
		opt.GeoM.Translate(p.px, p.py)
		opt.ColorM.Reset()
		opt.ColorM.Scale(p.clr, p.clb, p.clg, p.cla)
		if p.hue != 0 {
			opt.ColorM.RotateHue(p.hue)
		}
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
	c.Data(e).eid = e
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
	Vvr0, Vvr1                 float64         // randomized initial rotation velocity (radians/second)
	Vrab0, Vrab1               float64         // initial rotation acceleration (radians/second)
	Vrae0, Vrae1               float64         // end rotation acceleration (radians/second)
	Sx, Sy                     float64         // initial scale
	Vsx0, Vsy0, Vsx1, Vsy1     float64         // randomized initial scale
	Esx, Esy                   float64         // end scale
	Vesx0, Vesy0, Vesx1, Vesy1 float64         // randomized end scale
	Ox, Oy                     float64         // origin
	Dur                        float64         // duration
	Vdur0, Vdur1               float64         // randomized duration (seconds)
	Hueshift                   float64         // hue shift (/second)
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
		rv:        Lerpf(pp.Vvr0, pp.Vvr1, rng.Float64()),
		rab:       Lerpf(pp.Vrab0, pp.Vrab1, rng.Float64()),
		rae:       Lerpf(pp.Vrae0, pp.Vrae1, rng.Float64()),
		vx:        pp.Vx + Lerpf(pp.Vvx0, pp.Vvx1, rng.Float64()),
		vy:        pp.Vy + Lerpf(pp.Vvy0, pp.Vvy1, rng.Float64()),
		//sx: GET FROM bsx,
		//sy: GET FROM bsy,
		bsx:      pp.Sx + Lerpf(pp.Vsx0, pp.Vsx1, rng.Float64()),
		bsy:      pp.Sy + Lerpf(pp.Vsy0, pp.Vsy1, rng.Float64()),
		esx:      pp.Esx + Lerpf(pp.Vesx0, pp.Vesx1, rng.Float64()),
		esy:      pp.Esy + Lerpf(pp.Vesy0, pp.Vesy1, rng.Float64()),
		t:        0,
		hueshift: pp.Hueshift,
		hue:      0,
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
	p.parenttr = e.parenttr
	p.parenttre = e.parenttre
	return p
}

type EmissionProp struct {
	Fn      func(ctx UpdateCtx, e *ParticleEmitter) bool // use a custom function to emit particles
	T0, T1  float64                                      // emit a particle every Lerp(T0,T1,rand())
	N0, N1  int                                          // emit (N0 >= p > N1) particle(s) every Lerp(T0,T1,rand())
	Enabled bool
}

type Particle struct {
	px, py                 float64
	vx, vy                 float64
	ax, ay                 float64
	sx, sy                 float64
	bsx, bsy, esx, esy     float64
	r                      float64
	rv                     float64 // rotation velocity
	rab, rae               float64 // rotation acceleration
	pox, poy               float64 // precomputed origin
	bclr, bclg, bclb, bcla float64 // initial color tint
	eclr, eclg, eclb, ecla float64 // end color tint
	clr, clg, clb, cla     float64 // current color tint
	t, dur                 float64
	hueshift               float64
	hue                    float64
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

// DrawPriority noop
func (s *DrawableParticleEmitterSystem) DrawPriority(ctx DrawCtx) {

}

// Draw noop (drawing is controlled by *Drawable)
func (s *DrawableParticleEmitterSystem) Draw(ctx DrawCtx) {}

// UpdatePriority noop
func (s *DrawableParticleEmitterSystem) UpdatePriority(ctx UpdateCtx) {}

// Update computes labes if dirty
func (s *DrawableParticleEmitterSystem) Update(ctx UpdateCtx) {
	dt := ctx.DT()
	for _, v := range s.V().Matches() {
		if v.ParticleEmitter.disabled {
			continue
		}
		var del []int
		e := v.ParticleEmitter
		// if auto emission is enabled
		if e.emission.Enabled {
			if e.emission.Fn != nil {
				if e.emission.Fn(ctx, e) {
					_ = e.Emit()
				}
			} else {
				rng := e.rand
				if rng == nil {
					rng = emitterRNG
				}
				if e.emissiontnext == 0 {
					e.emissiontnext = Lerpf(e.emission.T0, e.emission.T1, rng.Float64())
				}
				e.emissiont += dt
				if e.emissiont >= e.emissiontnext {
					// emit N
					nn := e.emission.N0
					if e.emission.N0 < e.emission.N1 {
						nn = rng.Intn(e.emission.N1-e.emission.N0) + e.emission.N0
					}
					for i := 0; i < nn; i++ {
						e.Emit()
					}
					e.emissiont = 0
					e.emissiontnext = Lerpf(e.emission.T0, e.emission.T1, rng.Float64())
				}
			}
		}
		//
		for i := range e.particles {
			// upd time
			e.particles[i].t += dt
			t := e.particles[i].t
			pp := e.particles[i]
			ct := 0.0
			if pp.dur > 0 {
				ct = t / pp.dur
			}
			// velocity
			e.particles[i].vx += pp.ax * dt
			e.particles[i].vy += pp.ay * dt
			// size
			e.particles[i].sx = Lerpf(pp.bsx, pp.esx, ct)
			e.particles[i].sy = Lerpf(pp.bsy, pp.esy, ct)
			// color
			e.particles[i].clr = Lerpf(pp.bclr, pp.eclr, ct)
			e.particles[i].clg = Lerpf(pp.bclg, pp.eclg, ct)
			e.particles[i].clb = Lerpf(pp.bclb, pp.eclb, ct)
			e.particles[i].cla = Lerpf(pp.bcla, pp.ecla, ct)
			// position
			e.particles[i].px += e.particles[i].vx * dt
			e.particles[i].py += e.particles[i].vy * dt
			// rotation
			e.particles[i].rv += Lerpf(pp.rab, pp.rae, ct) * dt
			e.particles[i].r += e.particles[i].rv * dt
			// hue shifting
			e.particles[i].hue += pp.hueshift * dt
			if t > pp.dur {
				del = append(del, i)
			}
		}
		if len(del) > 0 {
			for i := len(del) - 1; i >= 0; i-- {
				x := del[i]
				e.particles = e.particles[:x+copy(e.particles[x:], e.particles[x+1:])]
			}
		}
	}
}

package graphics

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten/v2"
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
	max             int // max active particles
	props           ParticleProps
	particles       []Particle
	strategy        SpawnStrategy
	opt             ebiten.DrawImageOptions
	eid             ecs.Entity // cached entity id used to obtains its own transform
	ew              ecs.BaseWorld
	rand            *rand.Rand
	emission        EmissionProp
	emissiont       float64
	emissiontnext   float64
	compositeMode   ebiten.CompositeMode
	disabled        bool
	lockedparticles bool
	parentlevel     uint

	drawMask core.DrawMask
}

func NewParticleEmitter(w ecs.BaseWorld) ParticleEmitter {
	return ParticleEmitter{
		drawMask:  core.DrawMaskDefault,
		particles: make([]Particle, 0, 64),
		max:       64,
		strategy:  SpawnPause,
		props: ParticleProps{
			YVelocity:     -100,
			InitColor:     color.RGBA{255, 255, 255, 255},
			EndColor:      color.RGBA{255, 255, 255, 0},
			XVelocityVar0: -3,
			XVelocityVar1: 3,
			OriginX:       .5,
			OriginY:       .5,
			RotationVar0:  -math.Pi,
			RotationVar1:  math.Pi,
			Duration:      1,
			InitScale:     1,
			EndScale:      .5,
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

func (e *ParticleEmitter) DrawMask() core.DrawMask {
	return e.drawMask
}

func (e *ParticleEmitter) SetDrawMask(mask core.DrawMask) {
	e.drawMask = mask
}

func (e *ParticleEmitter) SetParentLevel(level uint) *ParticleEmitter {
	e.parentlevel = level
	return e
}

func (e *ParticleEmitter) SetLockedParticles(locked bool) *ParticleEmitter {
	e.lockedparticles = locked
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

func (e *ParticleEmitter) Emit(tr *components.Transform) bool {
	if len(e.props.Source) < 1 {
		return false
	}
	rng := e.rand
	if rng == nil {
		rng = emitterRNG
	}
	particle := e.props.NewParticle(rng)
	//TODO: link position with transform!
	if !e.lockedparticles {
		sys := components.GetTransformSystem(e.ew)
		gx, gy := sys.LocalToGlobalTr(particle.px, particle.py, tr)
		//a := particle.r + tr.angle
		// particle.px += tr.x
		// particle.py += tr.y
		// particle.r += tr.angle

		ttp, ttpe := tr.ParentTransform(), tr.Parent()
		for i := uint(0); i < e.parentlevel; i++ {
			if ttp == nil {
				break
			}
			ttpe = ttp.Parent()
			ttp = ttp.ParentTransform()
		}
		particle.parent = ttp
		particle.parente = ttpe
		if particle.parent != nil {
			lx, ly, _ := sys.GlobalToLocal(gx, gy, particle.parente)
			particle.px, particle.py = lx, ly
		} else {
			particle.px, particle.py = gx, gy
		}
		particle.sx *= tr.ScaleX()
		particle.sy *= tr.ScaleY()
	} else {
		particle.parent = tr
		particle.parente = e.eid
	}
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

func (e *ParticleEmitter) Update(ctx core.UpdateCtx, tr *components.Transform) {}

func (e *ParticleEmitter) Draw(ctx core.DrawCtx, tr *components.Transform) {
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
		opt.GeoM.Translate(p.px, p.py)
		if p.parent != nil {
			opt.GeoM.Concat(p.parent.GeoM())
		}
		opt.ColorM.Reset()
		opt.ColorM.Scale(p.clr, p.clg, p.clb, p.cla)
		if p.hue != 0 {
			opt.ColorM.RotateHue(p.hue)
		}
		ctx.Renderer().DrawImage(p.img, opt, e.drawMask)
	}
}

//go:generate ecsgen -n ParticleEmitter -p graphics -o particleemitter_component.go --component-tpl --vars "UUID=19A70DF9-0B1A-4A85-B23E-7BCA8E0857D7" --vars "BeforeRemove=c.beforeRemove(e)" --vars "OnAdd=c.onAdd(e)" --vars "Setup=c.onCompSetup()"

func (c *ParticleEmitterComponent) beforeRemove(e ecs.Entity) {
	c.Data(e).eid = 0
}

func (c *ParticleEmitterComponent) onAdd(e ecs.Entity) {
	c.Data(e).eid = e
}

func (c *ParticleEmitterComponent) onCompSetup() {
	RegisterDrawableComponent(c.world, c.flag, func(w ecs.BaseWorld, e ecs.Entity) Drawable {
		return GetParticleEmitterComponentData(w, e)
	})
}

type ParticleProps struct {
	X, Y                                       float64         // initial position
	XVar0, XVar1                               float64         // randomized position range
	YVar0, YVar1                               float64         // randomized position range
	XVelocity, YVelocity                       float64         // initial velocity
	XVelocityVar0, XVelocityVar1               float64         // randomized velocity range
	YVelocityVar0, YVelocityVar1               float64         // randomized velocity range
	XAccel, YAccel                             float64         // initial acceleration
	XAccelVar0, XAccelVar1                     float64         // randomized acceleration range
	YAccelVar0, YAccelVar1                     float64         // randomized acceleration range
	Rotation                                   float64         // initial rotation (radians)
	RotationVar0, RotationVar1                 float64         // randomized initial rotation (radians)
	RotationVelocityVar0, RotationVelocityVar1 float64         // randomized initial rotation velocity (radians/second)
	RotationAccelVar0, RotationAccelVar1       float64         // initial rotation acceleration (radians/second)
	EndRotationAccelVar0, EndRotationAccelVar1 float64         // end rotation acceleration (radians/second)
	InitScale                                  float64         // initial scale
	InitScaleVar0, InitScaleVar1               float64         // randomized initial scale
	EndScale                                   float64         // end scale
	EndScaleVar0, EndScaleVar1                 float64         // randomized end scale
	OriginX, OriginY                           float64         // origin
	Duration                                   float64         // duration
	DurationVar0, DurationVar1                 float64         // randomized duration (seconds)
	HueRotationSpeed                           float64         // hue shift (/second)
	InitColor                                  color.RGBA      // initial color tint modifier
	EndColor                                   color.RGBA      // end color tint modifier
	Source                                     []*ebiten.Image // particle source(s); if it's more than one, it will be randomized
}

func (pp *ParticleProps) SetPositionRange(xmin, xmax, ymin, ymax float64) {
	pp.XVar0, pp.XVar1 = xmin, xmax
	pp.YVar0, pp.YVar1 = ymin, ymax
}

func (pp *ParticleProps) SetVelocityRange(xmin, xmax, ymin, ymax float64) {
	pp.XVelocityVar0, pp.XVelocityVar1 = xmin, xmax
	pp.YVelocityVar0, pp.YVelocityVar1 = ymin, ymax
}

func (pp ParticleProps) NewParticle(rng *rand.Rand) Particle {
	initscale := pp.InitScale + core.Lerpf(pp.InitScaleVar0, pp.InitScaleVar1, rng.Float64())
	endscale := pp.EndScale + core.Lerpf(pp.EndScaleVar0, pp.EndScaleVar1, rng.Float64())
	p := Particle{
		bclr: float64(pp.InitColor.R) / 255,
		bclg: float64(pp.InitColor.G) / 255,
		bclb: float64(pp.InitColor.B) / 255,
		bcla: float64(pp.InitColor.A) / 255,
		clr:  float64(pp.InitColor.R) / 255,
		clg:  float64(pp.InitColor.G) / 255,
		clb:  float64(pp.InitColor.B) / 255,
		cla:  float64(pp.InitColor.A) / 255,
		eclr: float64(pp.EndColor.R) / 255,
		eclg: float64(pp.EndColor.G) / 255,
		eclb: float64(pp.EndColor.B) / 255,
		ecla: float64(pp.EndColor.A) / 255,
		//img: rng.Intn(len(pp.Source)),
		//pox: calc after obtaining image,
		//poy: calc after obtaining image,
		ax:       pp.XAccel + core.Lerpf(pp.XAccelVar0, pp.XAccelVar1, rng.Float64()),
		ay:       pp.YAccel + core.Lerpf(pp.YAccelVar0, pp.YAccelVar1, rng.Float64()),
		dur:      pp.Duration + core.Lerpf(pp.DurationVar0, pp.DurationVar1, rng.Float64()),
		px:       pp.X + core.Lerpf(pp.XVar0, pp.XVar1, rng.Float64()),
		py:       pp.Y + core.Lerpf(pp.YVar0, pp.YVar1, rng.Float64()),
		r:        pp.Rotation + core.Lerpf(pp.RotationVar0, pp.RotationVar1, rng.Float64()),
		rv:       core.Lerpf(pp.RotationVelocityVar0, pp.RotationVelocityVar1, rng.Float64()),
		rab:      core.Lerpf(pp.RotationAccelVar0, pp.RotationAccelVar1, rng.Float64()),
		rae:      core.Lerpf(pp.EndRotationAccelVar0, pp.EndRotationAccelVar1, rng.Float64()),
		vx:       pp.XVelocity + core.Lerpf(pp.XVelocityVar0, pp.XVelocityVar1, rng.Float64()),
		vy:       pp.YVelocity + core.Lerpf(pp.YVelocityVar0, pp.YVelocityVar1, rng.Float64()),
		sx:       initscale,
		sy:       initscale,
		bsx:      initscale,
		bsy:      initscale,
		esx:      endscale,
		esy:      endscale,
		t:        0,
		hueshift: pp.HueRotationSpeed,
		hue:      0,
	}
	if len(pp.Source) == 1 {
		p.img = pp.Source[0]
	} else if len(pp.Source) > 1 {
		p.img = pp.Source[rng.Intn(len(pp.Source))]
	}
	if p.img != nil {
		xx, yy := p.img.Size()
		p.pox = core.ApplyOrigin(float64(xx), pp.OriginX)
		p.poy = core.ApplyOrigin(float64(yy), pp.OriginY)
	}
	return p
}

type EmissionProp struct {
	Fn      func(ctx core.UpdateCtx, e *ParticleEmitter) bool // use a custom function to emit particles
	T0, T1  float64                                           // emit a particle every Lerp(T0,T1,rand())
	N0, N1  int                                               // emit (N0 >= p > N1) particle(s) every Lerp(T0,T1,rand())
	Enabled bool
}

type Particle struct {
	px, py                 float64 // position
	r                      float64 // rotation
	sx, sy                 float64 // scale
	vx, vy                 float64
	ax, ay                 float64
	bsx, bsy, esx, esy     float64
	rv                     float64 // rotation velocity
	rab, rae               float64 // rotation acceleration
	pox, poy               float64 // precomputed origin
	bclr, bclg, bclb, bcla float64 // initial color tint
	eclr, eclg, eclb, ecla float64 // end color tint
	clr, clg, clb, cla     float64 // current color tint
	t, dur                 float64
	hueshift               float64
	hue                    float64
	parente                ecs.Entity
	parent                 *components.Transform
	img                    *ebiten.Image
}

//go:generate ecsgen -n ParticleEmitter -p graphics -o particleemitter_drawablesystem.go --system-tpl --vars "Priority=10" --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "OnResize=s.onResize()" --vars "OnWillResize=s.onWillResize()" --vars "UUID=627C4B36-EE45-40C6-91AE-617D5CFDD8FC" --components "ParticleEmitter" --components "Transform;*components.Transform;components.GetTransformComponentData(v.world, e)" --go-import "\"github.com/gabstv/primen/components\""

var matchParticleEmitterSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	return eflag.Contains(components.GetTransformComponent(w).Flag().Or(GetParticleEmitterComponent(w).Flag()))
}

var resizematchParticleEmitterSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	if eflag.Contains(components.GetTransformComponent(w).Flag()) {
		return true
	}
	if eflag.Contains(GetParticleEmitterComponent(w).Flag()) {
		return true
	}
	return false
}

func (s *ParticleEmitterSystem) onEntityAdded(e ecs.Entity) {

}

func (s *ParticleEmitterSystem) onEntityRemoved(e ecs.Entity) {
	//if x := GetP
}

func (s *ParticleEmitterSystem) onWillResize() {
	v := s.V()
	for i := range v.entities {
		for k := range v.entities[i].ParticleEmitter.particles {
			v.entities[i].ParticleEmitter.particles[k].parent = nil
		}
	}
}

func (s *ParticleEmitterSystem) onResize() {
	etrs := make(map[ecs.Entity]*components.Transform)
	for _, match := range s.V().Matches() {
		for k := range match.ParticleEmitter.particles {
			if v := etrs[match.ParticleEmitter.particles[k].parente]; v != nil {
				match.ParticleEmitter.particles[k].parent = v
			} else {
				if ce := match.ParticleEmitter.particles[k].parente; ce != 0 {
					tr := components.GetTransformComponentData(s.world, ce)
					if tr != nil {
						etrs[ce] = tr
						match.ParticleEmitter.particles[k].parent = tr
					}
				}
			}
		}
	}
}

// DrawPriority noop
func (s *ParticleEmitterSystem) DrawPriority(ctx core.DrawCtx) {

}

// Draw noop (drawing is controlled by *Drawable)
func (s *ParticleEmitterSystem) Draw(ctx core.DrawCtx) {}

// UpdatePriority noop
func (s *ParticleEmitterSystem) UpdatePriority(ctx core.UpdateCtx) {}

// Update computes labes if dirty
func (s *ParticleEmitterSystem) Update(ctx core.UpdateCtx) {
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
					_ = e.Emit(v.Transform)
				}
			} else {
				rng := e.rand
				if rng == nil {
					rng = emitterRNG
				}
				if e.emissiontnext == 0 {
					e.emissiontnext = core.Lerpf(e.emission.T0, e.emission.T1, rng.Float64())
				}
				e.emissiont += dt
				if e.emissiont >= e.emissiontnext {
					// emit N
					nn := e.emission.N0
					if e.emission.N0 < e.emission.N1 {
						nn = rng.Intn(e.emission.N1-e.emission.N0) + e.emission.N0
					}
					for i := 0; i < nn; i++ {
						e.Emit(v.Transform)
					}
					e.emissiont = 0
					e.emissiontnext = core.Lerpf(e.emission.T0, e.emission.T1, rng.Float64())
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
			e.particles[i].sx = core.Lerpf(pp.bsx, pp.esx, ct)
			e.particles[i].sy = core.Lerpf(pp.bsy, pp.esy, ct)
			// color
			e.particles[i].clr = core.Lerpf(pp.bclr, pp.eclr, ct)
			e.particles[i].clg = core.Lerpf(pp.bclg, pp.eclg, ct)
			e.particles[i].clb = core.Lerpf(pp.bclb, pp.eclb, ct)
			e.particles[i].cla = core.Lerpf(pp.bcla, pp.ecla, ct)
			// position
			e.particles[i].px += e.particles[i].vx * dt
			e.particles[i].py += e.particles[i].vy * dt
			// rotation
			e.particles[i].rv += core.Lerpf(pp.rab, pp.rae, ct) * dt
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

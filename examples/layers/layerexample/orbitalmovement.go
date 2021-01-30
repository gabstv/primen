package layerexample

import (
	"math"
	"math/rand"

	"github.com/gabstv/ecs/v2"
	"github.com/gabstv/primen"
	"github.com/gabstv/primen/components"
	"github.com/gabstv/primen/components/graphics"
	"github.com/gabstv/primen/core"
	"github.com/hajimehoshi/ebiten/v2"
)

type OrbitalMovement struct {
	Speed       float64
	Dx          float64
	Dy          float64
	Ox          float64
	Oy          float64
	R           float64
	AngleR      float64
	ChildSprite ecs.Entity
	HueShift    bool
}

//go:generate ecsgen -n OrbitalMovement -p layerexample -o orbitalmovement_component.go --component-tpl --vars "UUID=DAD60C25-6B0D-4D3D-BF8E-5EB424FD8F1B"

//go:generate ecsgen -n OrbitalMovement -p layerexample -o orbitalmovement_system.go --system-tpl --vars "Priority=0" --vars "UUID=826684C9-E190-4BF2-93D7-2FA61A5BCEEC" --vars "Setup=s.setupVars()" --components "OrbitalMovement" --components "Sprite;*graphics.Sprite;graphics.GetSpriteComponent(v.world).Data(e)" --components "DrawLayer;*graphics.DrawLayer;graphics.GetDrawLayerComponent(v.world).Data(e)" --components "Transform;*components.Transform;components.GetTransformComponent(v.world).Data(e)" --go-import "\"github.com/gabstv/primen/components\"" --go-import "\"github.com/gabstv/primen/components/graphics\"" --go-import "\"github.com/hajimehoshi/ebiten\"" --members "paused bool" --members "globalScale float64" --members "radiusScale float64" --members "xframes chan struct{}" --members "wave1 float64" --members "waver float64" --members "fgs []*ebiten.Image" --members "bgs []*ebiten.Image"

var matchOrbitalMovementSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	// must contain
	f := graphics.GetDrawLayerComponent(w).Flag()
	f = f.Or(components.GetTransformComponent(w).Flag())
	f = f.Or(graphics.GetSpriteComponent(w).Flag())
	f = f.Or(GetOrbitalMovementComponent(w).Flag())
	return eflag.Contains(f)
}

var resizematchOrbitalMovementSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	if eflag.Contains(graphics.GetDrawLayerComponent(w).Flag()) {
		return true
	}
	if eflag.Contains(components.GetTransformComponent(w).Flag()) {
		return true
	}
	if eflag.Contains(graphics.GetSpriteComponent(w).Flag()) {
		return true
	}
	if eflag.Contains(GetOrbitalMovementComponent(w).Flag()) {
		return true
	}
	return false
}

// DrawPriority noop
func (s *OrbitalMovementSystem) DrawPriority(ctx core.DrawCtx) {}

// Draw noop
func (s *OrbitalMovementSystem) Draw(ctx core.DrawCtx) {}

// UpdatePriority noop
func (s *OrbitalMovementSystem) UpdatePriority(ctx core.UpdateCtx) {}

// Update positions
func (s *OrbitalMovementSystem) Update(ctx core.UpdateCtx) {
	if s.paused {
		select {
		case <-s.xframes:
			// run a single frame
		default:
			return
		}
	}
	dt := ctx.DT()
	s.waver += dt
	s.wave1 = 1 + math.Cos(s.waver)
	if s.waver > 2*math.Pi {
		s.waver -= 2 * math.Pi
	}
	for _, match := range s.V().Matches() {
		movecomp := match.OrbitalMovement
		movecomp.R += movecomp.Speed * dt * s.globalScale * s.wave1
		xx := math.Cos(movecomp.R) * movecomp.Dx * s.radiusScale
		yy := math.Sin(movecomp.R) * movecomp.Dy * s.radiusScale
		match.Transform.SetX(movecomp.Ox + xx)
		match.Transform.SetY(movecomp.Oy + yy)
		match.Transform.SetAngle(match.Transform.Angle() + (dt * (math.Pi / 4) * movecomp.AngleR))
		if rand.Float64() < 0.0001 {
			newlayer := rand.Intn(4)
			match.DrawLayer.Layer = graphics.LayerIndex(newlayer)
			match.Sprite.SetImage(s.bgs[newlayer])
			cspr := graphics.GetSpriteComponent(s.world).Data(movecomp.ChildSprite)
			cdl := graphics.GetDrawLayerComponent(s.world).Data(movecomp.ChildSprite)
			cspr.SetImage(s.fgs[newlayer])
			cdl.Layer = primen.Layer(newlayer)
			ctx.Engine().DispatchEvent("act_of_nature", match.Entity)
		}
		if movecomp.HueShift {
			// match.Sprite.SetColorHue(s.waver * 1.2)
			match.Sprite.RotateHue(dt * 2)
		}
	}
}

func (s *OrbitalMovementSystem) setupVars() {
	s.paused = false
	s.xframes = make(chan struct{}, 30)
	s.globalScale = 1
	s.radiusScale = 1
	s.wave1 = 1
	s.waver = 0
}

func (s *OrbitalMovementSystem) SetBgs(bgs []*ebiten.Image) {
	s.bgs = bgs
}

func (s *OrbitalMovementSystem) SetFgs(fgs []*ebiten.Image) {
	s.fgs = fgs
}

func (s *OrbitalMovementSystem) TogglePause() {
	s.paused = !s.paused
}

func (s *OrbitalMovementSystem) Paused() bool {
	return s.paused
}

func (s *OrbitalMovementSystem) PushFrame() {
	s.xframes <- struct{}{}
}

func (s *OrbitalMovementSystem) AddScale() {
	s.radiusScale += .1
}

func (s *OrbitalMovementSystem) SubScale() {
	s.radiusScale -= .1
}

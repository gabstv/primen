package gcs

import (
	"image"

	"github.com/gabstv/ecs"
	"github.com/gabstv/groove/pkg/groove"
	"github.com/gabstv/groove/pkg/groove/common"
	"github.com/hajimehoshi/ebiten"
)

type AnimClipMode byte

const (
	AnimClamp AnimClipMode = 1
	AnimLoop  AnimClipMode = 2
)

const (
	SpriteAnimationPriority int = -6
)

var (
	spriteanimationWC = &common.WorldComponents{}
)

func init() {
	groove.DefaultComp(func(e *groove.Engine, w *ecs.World) {
		//SpriteComponent(w)
	})
	groove.DefaultSys(func(e *groove.Engine, w *ecs.World) {
		//SpriteSystem(w)
	})
}

type SpriteAnimation struct {
	Enabled     bool
	Play        bool
	ActiveClip  int
	ActiveFrame int
	Clips       []SpriteAnimationClip
	T           float64
	// Default fps for clips with no fps specified
	Fps int

	// caches
	lastClipsLen int
	clipMap      map[string]int
}

type SpriteAnimationClip struct {
	// The name of an animation is not allowed to be changed during runtime
	// but since this is part of a component (and components shouldn't have logic),
	// it is a public member.
	Name   string
	Frames []image.Rectangle
	Fps    int
}

// SpriteAnimationComponent will get the registered sprite anim component of the world.
// If a component is not present, it will create a new component
// using world.NewComponent
func SpriteAnimationComponent(w *ecs.World) *ecs.Component {
	c := spriteanimationWC.Get(w)
	if c == nil {
		var err error
		c, err = w.NewComponent(ecs.NewComponentInput{
			Name: "groove.gcs.SpriteAnimation",
			ValidateDataFn: func(data interface{}) bool {
				_, ok := data.(*SpriteAnimation)
				return ok
			},
			DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
				sd := data.(*SpriteAnimation)
				sd.Clips = nil
			},
		})
		if err != nil {
			panic(err)
		}
		spriteanimationWC.Set(w, c)
	}
	return c
}

// SpriteAnimationSystem creates the sprite system
func SpriteAnimationSystem(w *ecs.World) *ecs.System {
	sys := w.NewSystem(SpriteAnimationPriority, SpriteAnimationSystemExec, spriteanimationWC.Get(w))
	if w.Get(DefaultImageOptions) == nil {
		opt := &ebiten.DrawImageOptions{}
		w.Set(DefaultImageOptions, opt)
	}
	sys.AddTag(groove.WorldTagUpdate)
	return sys
}

// SpriteAnimationSystemExec is the main function of the SpriteSystem
func SpriteAnimationSystemExec(dt float64, v *ecs.View, s *ecs.System) {
	world := v.World()
	matches := v.Matches()
	spriteanimcomp := spriteanimationWC.Get(world)
	//engine := world.Get(groove.EngineKey).(*groove.Engine)
	for _, m := range matches {
		spranim := m.Components[spriteanimcomp].(*SpriteAnimation)
		if !spranim.Enabled || !spranim.Play {
			continue
		}

	}
}

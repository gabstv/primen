package gcs

import (
	"image"

	"github.com/gabstv/ecs"
	"github.com/gabstv/groove/pkg/groove"
	"github.com/gabstv/groove/pkg/groove/common"
	"github.com/hajimehoshi/ebiten"
)

// AnimClipMode determines how time is treated outside of the keyframed range of an animation clip.
type AnimClipMode byte

const (
	// AnimOnce - When time reaches the end of the animation clip, the clip will automatically
	// stop playing and time will be reset to beginning of the clip.
	AnimOnce AnimClipMode = 0
	// AnimLoop - When time reaches the end of the animation clip, time will continue at the beginning.
	AnimLoop AnimClipMode = 1
	// AnimPingPong - When time reaches the end of the animation clip, time will ping pong back between beginning
	// and end.
	AnimPingPong AnimClipMode = 2 //TODO: implement ping pong
	// AnimClampForever - Plays back the animation. When it reaches the end, it will keep playing the last frame
	// and never stop playing.
	//
	// When playing backwards it will reach the first frame and will keep playing that. This is useful for additive
	// animations, which should never be stopped when they reach the maximum.
	AnimClampForever AnimClipMode = 3 //TODO: implement AnimClampForever
)

const (
	SpriteAnimationPriority     int = -6
	SpriteAnimationLinkPriority int = -5
)

var (
	spriteanimationWC = &common.WorldComponents{}
)

func init() {
	groove.DefaultComp(func(e *groove.Engine, w *ecs.World) {
		SpriteAnimationComponent(w)
	})
	groove.DefaultSys(func(e *groove.Engine, w *ecs.World) {
		SpriteAnimationSystem(w)
		SpriteAnimationLinkSystem(w)
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
	Fps float64

	// caches
	lastClip int
	//lastClipsLen int
	lastPlay bool
	clipMap  map[string]int
}

type SpriteAnimationClip struct {
	// The name of an animation is not allowed to be changed during runtime
	// but since this is part of a component (and components shouldn't have logic),
	// it is a public member.
	Name     string
	Frames   []image.Rectangle
	Fps      float64
	ClipMode AnimClipMode
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
	sys.AddTag(groove.WorldTagUpdate)
	return sys
}

// SpriteAnimationSystemExec is the main function of the SpriteSystem
func SpriteAnimationSystemExec(dt float64, v *ecs.View, s *ecs.System) {
	world := v.World()
	matches := v.Matches()
	spriteanimcomp := spriteanimationWC.Get(world)
	globalfps := nonzeroval(ebiten.CurrentFPS(), 60)
	//engine := world.Get(groove.EngineKey).(*groove.Engine)
	for _, m := range matches {
		spranim := m.Components[spriteanimcomp].(*SpriteAnimation)
		if !spranim.Enabled || !spranim.Play {
			if !spranim.Play && spranim.lastPlay {
				spranim.lastPlay = false
			}
			continue
		}
		clip := spranim.Clips[spranim.ActiveClip]
		localfps := nonzeroval(clip.Fps, spranim.Fps, globalfps)
		localdt := (dt * localfps) / globalfps
		if !spranim.lastPlay {
			// the animation was stopped on the last iteration
			if spranim.lastClip == spranim.ActiveClip {
				// since it is the same clip, this can be affected by
				// the clip AnimClipMode
				//TODO: handle AnimClipMode behavior
			}
		}
		spranim.lastClip = spranim.ActiveClip
		spranim.lastPlay = true
		spranim.T += localdt * localfps
		if spranim.T >= 1 {
			// next frame
			nextframe := spranim.ActiveFrame + 1
			//spranim.T -= 1
			if nextframe >= len(clip.Frames) {
				// animation ended
				switch clip.ClipMode {
				case AnimOnce:
					spranim.T = 0
					spranim.Play = false
				case AnimLoop:
					spranim.T = Clamp(spranim.T-1, 0, 1)
					spranim.ActiveFrame = 0
					//TODO: other clip modes
				}
			} else {
				spranim.T = Clamp(spranim.T-1, 0, 1)
				spranim.ActiveFrame = nextframe
			}
		}
	}
}

// SpriteAnimationLinkSystem creates the sprite system
func SpriteAnimationLinkSystem(w *ecs.World) *ecs.System {
	sys := w.NewSystem(SpriteAnimationLinkPriority, SpriteAnimationLinkSystemExec, spriteanimationWC.Get(w), spriteWC.Get(w))
	sys.AddTag(groove.WorldTagDraw)
	return sys
}

// SpriteAnimationLinkSystemExec is what glues the animation and sprite together
func SpriteAnimationLinkSystemExec(dt float64, v *ecs.View, s *ecs.System) {
	world := v.World()
	matches := v.Matches()
	spriteanimcomp := spriteanimationWC.Get(world)
	spritecomp := spriteWC.Get(world)
	//globalfps := ebiten.CurrentFPS()
	//engine := world.Get(groove.EngineKey).(*groove.Engine)
	for _, m := range matches {
		spranim := m.Components[spriteanimcomp].(*SpriteAnimation)
		spr := m.Components[spritecomp].(*Sprite)
		if !spranim.Enabled {
			continue
		}
		spr.Bounds = spranim.Clips[spranim.ActiveClip].Frames[spranim.ActiveFrame]
	}
}

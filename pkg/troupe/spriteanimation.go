package troupe

import (
	"image"

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
	AnimPingPong AnimClipMode = 2
	// AnimClampForever - Plays back the animation. When it reaches the end, it will keep playing the last frame
	// and never stop playing.
	//
	// When playing backwards it will reach the first frame and will keep playing that. This is useful for code
	// that uses the Playing bool to avoid activating said trigger.
	AnimClampForever AnimClipMode = 3
)

const (
	// SpriteAnimationPriority - default -6
	SpriteAnimationPriority int = -6
	// SpriteAnimationLinkPriority - default -5
	SpriteAnimationLinkPriority int = -5
)

var (
	spriteanimationWC = &WorldComponents{}
)

func init() {
	DefaultComp(func(e *Engine, w *World) {
		SpriteAnimationComponent(w)
	})
	DefaultSys(func(e *Engine, w *World) {
		SpriteAnimationSystem(w)
		SpriteAnimationLinkSystem(w)
	})
}

// SpriteAnimation holds the data of a sprite animation (and clips)
type SpriteAnimation struct {
	Enabled     bool
	Playing     bool
	ActiveClip  int
	ActiveFrame int
	Clips       []SpriteAnimationClip
	T           float64
	// Default fps for clips with no fps specified
	Fps float64

	// caches
	lastClip          int
	lastPlaying       bool
	clipMap           map[string]int
	clipMapLen        int
	nextAnimationName string
	nextAnimationSet  bool
	reversed          bool
}

// PlayClip sets the animation to play a clip by name
func (a *SpriteAnimation) PlayClip(name string) {
	a.nextAnimationName = name
	a.nextAnimationSet = true
}

// SpriteAnimationClip is an animation clip, like a character walk cycle.
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
func SpriteAnimationComponent(w Worlder) *Component {
	c := spriteanimationWC.Get(w)
	if c == nil {
		var err error
		c, err = w.NewComponent(NewComponentInput{
			Name: "troupe.SpriteAnimation",
			ValidateDataFn: func(data interface{}) bool {
				_, ok := data.(*SpriteAnimation)
				return ok
			},
			DestructorFn: func(_ Worlder, entity Entity, data interface{}) {
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
func SpriteAnimationSystem(w *World) *System {
	sys := w.NewSystem(SpriteAnimationPriority, SpriteAnimationSystemExec, spriteanimationWC.Get(w))
	sys.AddTag(WorldTagUpdate)
	return sys
}

// SpriteAnimationSystemExec is the main function of the SpriteSystem
func SpriteAnimationSystemExec(ctx Context, screen *ebiten.Image) {
	//dt float64, v *ecs.View, s *ecs.System
	dt := ctx.DT()
	v := ctx.System().View()
	world := v.World()
	matches := v.Matches()
	spriteanimcomp := spriteanimationWC.Get(world)
	globalfps := nonzeroval(ebiten.CurrentFPS(), 60)
	for _, m := range matches {
		spranim := m.Components[spriteanimcomp].(*SpriteAnimation)
		if !spranim.Enabled || !spranim.Playing {
			if !spranim.Playing && spranim.lastPlaying {
				spranim.lastPlaying = false
			}
			continue
		}
		spriteAnimResolveClipMap(spranim)
		spriteAnimResolvePlayClip(spranim)
		spriteAnimResolvePlayback(globalfps, dt, spranim)
	}
}

func spriteAnimResolvePlayClip(spranim *SpriteAnimation) {
	if !spranim.nextAnimationSet {
		return
	}
	spranim.nextAnimationSet = false
	index, ok := spranim.clipMap[spranim.nextAnimationName]
	spranim.nextAnimationName = ""
	if !ok {
		return
	}
	spranim.T = 0
	spranim.Playing = true
	spranim.ActiveFrame = 0
	spranim.ActiveClip = index
	spranim.reversed = false
}

func spriteAnimResolveClipMap(spranim *SpriteAnimation) {
	if spranim.clipMapLen == len(spranim.Clips) {
		return
	}
	// rebuild cache
	spranim.clipMap = make(map[string]int)
	for k, v := range spranim.Clips {
		spranim.clipMap[v.Name] = k
	}
	spranim.clipMapLen = len(spranim.Clips)
}

func spriteAnimResolvePlayback(globalfps, dt float64, spranim *SpriteAnimation) {
	clip := spranim.Clips[spranim.ActiveClip]
	localfps := nonzeroval(clip.Fps, spranim.Fps, globalfps)
	localdt := (dt * localfps) / globalfps
	if !spranim.lastPlaying {
		// the animation was stopped on the last iteration
		if spranim.lastClip == spranim.ActiveClip {
			// since it is the same clip, this can be affected by
			// the clip AnimClipMode
			//TODO: handle AnimClipMode behavior
			//TODO: maybe triggers
		}
	}
	if spranim.lastClip != spranim.ActiveClip {
		// reset the reversed state
		spranim.reversed = false
	}
	spranim.lastClip = spranim.ActiveClip
	spranim.lastPlaying = true
	spranim.T += localdt * localfps
	if spranim.T >= 1 {
		// next frame
		nextframe := spranim.ActiveFrame + 1
		if spranim.reversed {
			nextframe = spranim.ActiveFrame - 1
		}
		if nextframe >= len(clip.Frames) {
			// animation ended
			switch clip.ClipMode {
			case AnimOnce:
				spranim.T = 0
				spranim.Playing = false
			case AnimLoop:
				spranim.T = Clamp(spranim.T-1, 0, 1)
				spranim.ActiveFrame = 0
			case AnimPingPong:
				spranim.T = Clamp(spranim.T-1, 0, 1)
				spranim.reversed = true
			case AnimClampForever:
				// the last frame will keep on playing
				spranim.T = Clamp(spranim.T-1, 0, 1)
			}
		} else if nextframe < 0 {
			// the animation is reversed and reached the beginning
			spranim.T = Clamp(spranim.T-1, 0, 1)
			spranim.reversed = false
		} else {
			spranim.T = Clamp(spranim.T-1, 0, 1)
			spranim.ActiveFrame = nextframe
		}
	}
}

// SpriteAnimationLinkSystem creates the sprite system
func SpriteAnimationLinkSystem(w *World) *System {
	sys := w.NewSystem(SpriteAnimationLinkPriority, SpriteAnimationLinkSystemExec, spriteanimationWC.Get(w), w.Component(spriteComponentName))
	sys.AddTag(WorldTagDraw)
	return sys
}

// SpriteAnimationLinkSystemExec is what glues the animation and sprite together
func SpriteAnimationLinkSystemExec(ctx Context, screen *ebiten.Image) {
	v := ctx.System().View()
	world := ctx.World()
	matches := v.Matches()
	spriteanimcomp := spriteanimationWC.Get(world)
	spritecomp := world.Component(spriteComponentName)
	for _, m := range matches {
		spranim := m.Components[spriteanimcomp].(*SpriteAnimation)
		spr := m.Components[spritecomp].(*Sprite)
		if !spranim.Enabled {
			continue
		}
		spr.Bounds = spranim.Clips[spranim.ActiveClip].Frames[spranim.ActiveFrame]
	}
}

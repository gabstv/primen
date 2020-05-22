package core

import (
	"image"
	"sync"

	"github.com/gabstv/ecs"
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
	SNSpriteAnimation     = "primen.SpriteAnimationSystem"
	SNSpriteAnimationLink = "primen.SpriteAnimationLinkSystem"
	CNSpriteAnimation     = "primen.SpriteAnimationComponent"
)

type SpriteAnimationComponentSystem struct {
	BaseComponentSystem
}

func (cs *SpriteAnimationComponentSystem) SystemName() string {
	return SNSpriteAnimation
}

func (cs *SpriteAnimationComponentSystem) SystemPriority() int {
	return -6
}

func (cs *SpriteAnimationComponentSystem) SystemExec() SystemExecFn {
	return SpriteAnimationSystemExec
}

func (cs *SpriteAnimationComponentSystem) SystemTags() []string {
	return []string{
		"draw",
	}
}

func (cs *SpriteAnimationComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		spriteAnimationComponentDef(w),
	}
}

func (cs *SpriteAnimationComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
	return emptyCompSlice
}

func spriteAnimationComponentDef(w *ecs.World) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNSpriteAnimation,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*SpriteAnimation)
			return ok
		},
		DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
			sd := data.(*SpriteAnimation)
			sd.Clips = nil
		},
	})
}

type AnimationEventID int64

type AnimationEventListeners struct {
	nextid int64
	l      sync.Mutex
	m      map[string][]AnimationEventID
	m2     map[AnimationEventID]AnimationEventFn
}

func (ae *AnimationEventListeners) ensure(group string) {
	ae.l.Lock()
	defer ae.l.Unlock()
	if ae.m == nil {
		ae.m = make(map[string][]AnimationEventID)
	}
	if ae.m2 == nil {
		ae.m2 = make(map[AnimationEventID]AnimationEventFn)
	}
	if group != "" {
		if ae.m[group] == nil {
			ae.m[group] = make([]AnimationEventID, 0, 2)
		}
	}
}

func (ae *AnimationEventListeners) Add(name string, fn AnimationEventFn) AnimationEventID {
	ae.ensure(name)
	ae.l.Lock()
	defer ae.l.Unlock()
	ae.nextid++
	id := AnimationEventID(ae.nextid)
	ae.m[name] = append(ae.m[name], id)
	ae.m2[id] = fn
	return id
}

func (ae *AnimationEventListeners) AddCatchAll(fn AnimationEventFn) AnimationEventID {
	return ae.Add("*", fn)
}

func (ae *AnimationEventListeners) Remove(id AnimationEventID) bool {
	ae.ensure("")
	ae.l.Lock()
	defer ae.l.Unlock()
	if _, ok := ae.m2[id]; !ok {
		return false
	}
	// maybe switch to a more performat approach if there's a use case for adding/deleting
	// a large volume of event observers on a single animation controller
	fnd := false
	for k, v := range ae.m {
		if v != nil {
			vix := -1
			for vi, vid := range v {
				if vid == id {
					vix = vi
					break
				}
			}
			if vix > -1 {
				v = append(v[:vix], v[vix+1:]...)
				ae.m[k] = v
				fnd = true
				delete(ae.m2, id)
				break
			}
		}
	}
	return fnd
}

func (ae *AnimationEventListeners) Clear() {
	ae.l.Lock()
	defer ae.l.Unlock()
	ae.m = nil
	for k := range ae.m2 {
		ae.m2[k] = nil
	}
	ae.m = nil
	ae.m2 = nil
}

func (ae *AnimationEventListeners) Dispatch(name, value string) {
	ae.ensure("")
	ae.l.Lock()
	defer ae.l.Unlock()
	x := ae.m[name]
	if x != nil {
		for _, id := range x {
			if ae.m2[id] != nil {
				go ae.m2[id](name, value)
			}
		}
	}
	x = ae.m["*"]
	if x != nil {
		for _, id := range x {
			if ae.m2[id] != nil {
				go ae.m2[id](name, value)
			}
		}
	}
}

type AnimationEventFn func(name, value string)

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
	lastImage         *ebiten.Image
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

func (a *SpriteAnimation) AnimEvent(name, value string) {
	//FIXME: implement this
}

// SpriteAnimationClip is an animation clip, like a character walk cycle.
type SpriteAnimationClip struct {
	// The name of an animation is not allowed to be changed during runtime
	// but since this is part of a component (and components shouldn't have logic),
	// it is a public member.
	Name       string
	Image      *ebiten.Image
	Frames     []image.Rectangle
	Events     []*SpriteAnimationEvent //TODO: link
	Fps        float64
	ClipMode   AnimClipMode
	EndedEvent *SpriteAnimationEvent //TODO: link
}

type SpriteAnimationEvent struct {
	Name  string
	Value string
}

// SpriteAnimationSystemExec is the main function of the SpriteSystem
func SpriteAnimationSystemExec(ctx Context) {
	//screen := ctx.Screen()
	//dt float64, v *ecs.View, s *ecs.System
	dt := ctx.DT()
	v := ctx.System().View()
	matches := v.Matches()
	spriteanimcomp := ctx.World().Component(CNSpriteAnimation)
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
	if evs := spranim.Clips[index].Events; evs != nil && len(evs) > 0 && evs[0] != nil {
		spranim.AnimEvent(evs[0].Name, evs[0].Value)
	}
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
	spranim.lastImage = spranim.Clips[spranim.lastClip].Image
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
				if clip.EndedEvent != nil {
					spranim.AnimEvent(clip.EndedEvent.Name, clip.EndedEvent.Value)
				}
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
			// dispatch event!
			if evs := clip.Events; evs != nil && len(evs) > nextframe && evs[nextframe] != nil {
				spranim.AnimEvent(evs[nextframe].Name, evs[nextframe].Value)
			}
		}
	}
}

type SpriteAnimationLinkComponentSystem struct {
	BaseComponentSystem
}

func (cs *SpriteAnimationLinkComponentSystem) SystemName() string {
	return SNSpriteAnimationLink
}

func (cs *SpriteAnimationLinkComponentSystem) SystemPriority() int {
	return -5
}

func (cs *SpriteAnimationLinkComponentSystem) SystemExec() SystemExecFn {
	return SpriteAnimationLinkSystemExec
}

func (cs *SpriteAnimationLinkComponentSystem) SystemTags() []string {
	return []string{"draw"}
}

func (cs *SpriteAnimationLinkComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		spriteAnimationComponentDef(w),
		drawableComponentDef(w),
	}
}

func (cs *SpriteAnimationLinkComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
	return emptyCompSlice
}

// SpriteAnimationLinkSystemExec is what glues the animation and sprite together
func SpriteAnimationLinkSystemExec(ctx Context) {
	//screen := ctx.Screen()
	v := ctx.System().View()
	world := ctx.World()
	matches := v.Matches()
	spriteanimcomp := world.Component(CNSpriteAnimation)
	spritecomp := world.Component(CNDrawable)
	for _, m := range matches {
		spranim := m.Components[spriteanimcomp].(*SpriteAnimation)
		spr := m.Components[spritecomp].(Drawable)
		if !spranim.Enabled {
			continue
		}
		// replace image if the animation clip usaes a different one
		if spranim.lastImage != nil {
			if w, ok := spr.(DrawableImager); ok {
				if w.GetImage() != spranim.lastImage {
					w.SetImage(spranim.lastImage)
				}
			}
		}
		spr.SetBounds(spranim.Clips[spranim.ActiveClip].Frames[spranim.ActiveFrame])
	}
}

func init() {
	RegisterComponentSystem(&SpriteAnimationComponentSystem{})
	RegisterComponentSystem(&SpriteAnimationLinkComponentSystem{})
}

package core

import (
	"github.com/gabstv/ecs/v2"
)

// SpriteAnimation holds the data of a sprite animation (and clips)
type SpriteAnimation struct {
	playing     bool
	activeClip  AnimationClip
	activeFrame int
	anim        Animation
	t           float64
	reversed    bool
	fps         float64                      // Default fps for clips with no fps specified
	clipMap     map[string]AnimationClip     // clip cache
	clipEvents  map[string][]*AnimationEvent // clip cache

	listeners *AnimationEventListeners
}

func NewSpriteAnimation(fps float64, anim Animation) SpriteAnimation {
	sa := SpriteAnimation{
		activeFrame: -1,
		fps:         fps,
		listeners:   &AnimationEventListeners{},
	}
	if anim != nil {
		sa.SetAnimation(anim)
	}
	return sa
}

func (a *SpriteAnimation) Animation() Animation {
	return a.anim
}

// SetAnimation parses the clip map and animations
func (a *SpriteAnimation) SetAnimation(anim Animation) {
	if anim == nil {
		panic("anim cannot be nil")
	}
	n := anim.Count()
	a.anim = anim
	a.clipMap = make(map[string]AnimationClip)
	a.clipEvents = make(map[string][]*AnimationEvent)
	for i := 0; i < n; i++ {
		clip := anim.GetClip(i)
		a.clipMap[clip.GetName()] = clip
		a.clipEvents[clip.GetName()] = a.anim.GetClipEvents(i)
	}
	a.reset()
}

func (a *SpriteAnimation) Playing() bool {
	return a.playing
}

func (a *SpriteAnimation) Reversed() bool {
	return a.reversed
}

// SetReversed sets "reversed". It does nothing if it's not playing.
func (a *SpriteAnimation) SetReversed(reversed bool) {
	a.reversed = reversed
}

// PlayClip sets the animation to play a clip by name
func (a *SpriteAnimation) PlayClip(name string) bool {
	return a.PlayClipFrame(name, 0)
}

func (a *SpriteAnimation) AddEventListener(name string, fn AnimationEventFn) AnimationEventID {
	return a.listeners.Add(name, fn)
}

func (a *SpriteAnimation) AddEventListenerW(fn AnimationEventFn) AnimationEventID {
	return a.listeners.AddCatchAll(fn)
}

func (a *SpriteAnimation) RemoveEventListener(id AnimationEventID) bool {
	return a.listeners.Remove(id)
}

// PlayClip sets the animation to play a clip by name and at a specific frame
func (a *SpriteAnimation) PlayClipFrame(name string, frame int) bool {
	clip := a.clipMap[name]
	if clip == nil {
		return false
	}
	return a.play(clip, frame)
}

// PlayClipIndex plays a clip at index i
func (a *SpriteAnimation) PlayClipIndex(i int) bool {
	return a.PlayClipIndexFrame(i, 0)
}

// PlayClipIndexFrame plays a clip at index i
func (a *SpriteAnimation) PlayClipIndexFrame(i, frame int) bool {
	if i < 0 {
		return false
	}
	if a.anim.Count() <= i {
		return false
	}
	anim := a.anim.GetClip(i)
	return a.play(anim, frame)
}

// AnimEvent dispatches an animation event to listeners
func (a *SpriteAnimation) AnimEvent(name, value string) {
	a.listeners.Dispatch(name, value)
}

func (a *SpriteAnimation) tryDispatchClipEvent(clip AnimationClip, frame int) {
	if a.clipEvents[clip.GetName()] == nil {
		return
	}
	if len(a.clipEvents[clip.GetName()]) <= frame {
		return
	}
	if evt := a.clipEvents[clip.GetName()][frame]; evt != nil {
		a.AnimEvent(evt.Name, evt.Value)
	}
}

func (a *SpriteAnimation) trySwapImage(clip AnimationClip, frame int, sprite *Sprite) {
	img := clip.GetImage(frame)
	if img != nil && sprite.Image() != img {
		sprite.SetImage(img)
		sprite.SetOffset(clip.GetOffset(frame))
	}
}

func (a *SpriteAnimation) play(clip AnimationClip, frame int) bool {
	if frame < 0 {
		return false
	}
	if clip.GetFrameCount() <= frame {
		return false
	}
	a.t = 0
	a.playing = true
	a.activeClip = clip
	a.activeFrame = frame
	a.reversed = false
	a.tryDispatchClipEvent(clip, frame)
	return true
}

func (a *SpriteAnimation) reset() {
	//TODO: play "default" animation if set
	a.playing = false
	a.activeClip = nil
	a.activeFrame = -1
	a.t = 0
	a.reversed = false
}

//go:generate ecsgen -n SpriteAnimation -p core -o spriteanimation_component.go --component-tpl --vars "UUID=5A056275-C47D-44D2-994C-BD0AF107870C"

//go:generate ecsgen -n SpriteAnimation -p core -o spriteanimation_system.go --system-tpl --vars "Priority=12" --vars "EntityAdded=s.onEntityAdded(e)" --vars "EntityRemoved=s.onEntityRemoved(e)" --vars "UUID=FFD3127E-6066-4561-8B2A-E1B59EBE489C" --components "Sprite" --components "SpriteAnimation"

var matchSpriteAnimationSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if !f.Contains(GetSpriteComponent(w).Flag()) {
		return false
	}
	if !f.Contains(GetSpriteAnimationComponent(w).Flag()) {
		return false
	}
	return true
}

var resizematchSpriteAnimationSystem = func(f ecs.Flag, w ecs.BaseWorld) bool {
	if f.Contains(GetSpriteComponent(w).Flag()) {
		return true
	}
	if f.Contains(GetSpriteAnimationComponent(w).Flag()) {
		return true
	}
	return false
}

func (s *SpriteAnimationSystem) onEntityAdded(e ecs.Entity) {
	if spranim := GetSpriteAnimationComponentData(s.world, e); spranim.playing &&
		spranim.activeClip != nil && spranim.activeFrame > -1 &&
		spranim.activeFrame < spranim.activeClip.GetFrameCount() {
		sprite := GetSpriteComponentData(s.world, e)
		sprite.SetImage(spranim.activeClip.GetImage(spranim.activeFrame))
		sprite.SetOffset(spranim.activeClip.GetOffset(spranim.activeFrame))
	}
}

func (s *SpriteAnimationSystem) onEntityRemoved(e ecs.Entity) {

}

// DrawPriority noop
func (s *SpriteAnimationSystem) DrawPriority(ctx DrawCtx) {}

// Draw noop
func (s *SpriteAnimationSystem) Draw(ctx DrawCtx) {}

// UpdatePriority noop
func (s *SpriteAnimationSystem) UpdatePriority(ctx UpdateCtx) {}

// Update plays the animations
func (s *SpriteAnimationSystem) Update(ctx UpdateCtx) {
	dt := ctx.DT()
	globalfps := nonzeroval(ctx.TPS(), 60)
	for _, v := range s.V().Matches() {
		if !v.SpriteAnimation.playing {
			continue
		}
		frame := v.SpriteAnimation.activeFrame
		if frame < 0 {
			continue
		}
		clip := v.SpriteAnimation.activeClip
		if clip == nil {
			continue
		}
		v.SpriteAnimation.trySwapImage(clip, frame, v.Sprite)

		localfps := nonzeroval(clip.GetFPS(), v.SpriteAnimation.fps, globalfps)

		at := (localfps * dt)
		v.SpriteAnimation.t += at

		if v.SpriteAnimation.t >= 1 {
			// next frame
			nextframe := frame + 1
			if v.SpriteAnimation.reversed {
				nextframe = frame - 1
			}
			if nextframe >= clip.GetFrameCount() {
				// reached end
				switch clip.GetMode() {
				case AnimOnce:
					v.SpriteAnimation.t = 0
					v.SpriteAnimation.playing = false
					if e := clip.GetEndedEvent(); e != nil {
						v.SpriteAnimation.AnimEvent(e.Name, e.Value)
					}
				case AnimLoop:
					v.SpriteAnimation.t = Clamp(v.SpriteAnimation.t-1, 0, 1)
					v.SpriteAnimation.activeFrame = 0
					v.SpriteAnimation.trySwapImage(clip, 0, v.Sprite)
				case AnimPingPong:
					v.SpriteAnimation.t = Clamp(v.SpriteAnimation.t-1, 0, 1)
					v.SpriteAnimation.activeFrame = intmax(clip.GetFrameCount()-2, 0)
					v.SpriteAnimation.reversed = true
					v.SpriteAnimation.trySwapImage(clip, v.SpriteAnimation.activeFrame, v.Sprite)
					//TODO: option to repeat last and first frames (?)
				case AnimClampForever:
					v.SpriteAnimation.t = Clamp(v.SpriteAnimation.t-1, 0, 1)
					v.SpriteAnimation.activeFrame = clip.GetFrameCount() - 1
				}
			} else if nextframe < 0 {
				switch clip.GetMode() {
				case AnimOnce:
					v.SpriteAnimation.t = 0
					v.SpriteAnimation.playing = false
					if e := clip.GetEndedEvent(); e != nil {
						v.SpriteAnimation.AnimEvent(e.Name, e.Value)
					}
				case AnimLoop:
					v.SpriteAnimation.t = Clamp(v.SpriteAnimation.t-1, 0, 1)
					v.SpriteAnimation.activeFrame = clip.GetFrameCount() - 1
					v.SpriteAnimation.trySwapImage(clip, v.SpriteAnimation.activeFrame, v.Sprite)
				case AnimPingPong:
					v.SpriteAnimation.t = Clamp(v.SpriteAnimation.t-1, 0, 1)
					v.SpriteAnimation.activeFrame = 0
					v.SpriteAnimation.reversed = false
					v.SpriteAnimation.trySwapImage(clip, v.SpriteAnimation.activeFrame, v.Sprite)
					//TODO: option to repeat last and first frames (?)
				case AnimClampForever:
					v.SpriteAnimation.t = Clamp(v.SpriteAnimation.t-1, 0, 1)
					v.SpriteAnimation.activeFrame = 0
				}
			} else {
				// boundaries not reached
				v.SpriteAnimation.t = Clamp(v.SpriteAnimation.t-1, 0, 1)
				v.SpriteAnimation.activeFrame = nextframe
				v.SpriteAnimation.trySwapImage(clip, v.SpriteAnimation.activeFrame, v.Sprite)
				v.SpriteAnimation.tryDispatchClipEvent(clip, nextframe)
			}
		}
	}
}

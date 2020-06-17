package core

import (
	"github.com/gabstv/ecs/v2"
	"github.com/hajimehoshi/ebiten"
)

// const (
// 	SNSpriteAnimation     = "primen.SpriteAnimationSystem"
// 	SNSpriteAnimationLink = "primen.SpriteAnimationLinkSystem"
// 	CNSpriteAnimation     = "primen.SpriteAnimationComponent"
// )

// type SpriteAnimationComponentSystem struct {
// 	BaseComponentSystem
// }

// func (cs *SpriteAnimationComponentSystem) SystemName() string {
// 	return SNSpriteAnimation
// }

// func (cs *SpriteAnimationComponentSystem) SystemPriority() int {
// 	return -6
// }

// func (cs *SpriteAnimationComponentSystem) SystemExec() SystemExecFn {
// 	return SpriteAnimationSystemExec
// }

// func (cs *SpriteAnimationComponentSystem) SystemTags() []string {
// 	return []string{
// 		"draw",
// 	}
// }

// func (cs *SpriteAnimationComponentSystem) Components(w *ecs.World) []*ecs.Component {
// 	return []*ecs.Component{
// 		spriteAnimationComponentDef(w),
// 	}
// }

// func (cs *SpriteAnimationComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
// 	return emptyCompSlice
// }

// func spriteAnimationComponentDef(w *ecs.World) *ecs.Component {
// 	return UpsertComponent(w, ecs.NewComponentInput{
// 		Name: CNSpriteAnimation,
// 		ValidateDataFn: func(data interface{}) bool {
// 			_, ok := data.(*SpriteAnimation)
// 			return ok
// 		},
// 		DestructorFn: func(_ *ecs.World, entity ecs.Entity, data interface{}) {
// 			sd := data.(*SpriteAnimation)
// 			sd.Anim = nil
// 		},
// 	})
// }

// SpriteAnimation holds the data of a sprite animation (and clips)
type SpriteAnimation struct {
	Enabled     bool
	Playing     bool
	ActiveClip  int
	ActiveFrame int
	Anim        Animation
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
	a.Playing = true
}

// PlayClipIndex plays a clip at index i
func (a *SpriteAnimation) PlayClipIndex(i int) {
	if i < 0 {
		return
	}
	if a.Anim.Count() <= i {
		return
	}
	a.nextAnimationName = a.Anim.GetClip(i).GetName()
	a.nextAnimationSet = true
	a.Playing = true
}

func (a *SpriteAnimation) AnimEvent(name, value string) {
	//FIXME: implement this
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

}

func (s *SpriteAnimationSystem) onEntityRemoved(e ecs.Entity) {

}

func (s *SpriteAnimationSystem) DrawPriority(ctx DrawCtx) {

}

func (s *SpriteAnimationSystem) Draw(ctx DrawCtx) {

}

func (s *SpriteAnimationSystem) UpdatePriority(ctx UpdateCtx) {

}

func (s *SpriteAnimationSystem) Update(ctx UpdateCtx) {
	dt := ctx.DT()
	globalfps := nonzeroval(ctx.TPS(), 60)
	for _, v := range s.V().Matches() {
		if !v.SpriteAnimation.Enabled || !v.SpriteAnimation.Playing {
			if !v.SpriteAnimation.Playing && v.SpriteAnimation.lastPlaying {
				v.SpriteAnimation.lastPlaying = false
			}
			continue
		}
		spriteAnimResolveClipMap(v.SpriteAnimation)
		spriteAnimResolvePlayClip(v.SpriteAnimation)
		spriteAnimResolvePlayback(globalfps, dt, v.SpriteAnimation)
	}
	for _, v := range s.V().Matches() {
		if !v.SpriteAnimation.Enabled {
			continue
		}
		// replace image if the animation clip uses a different one
		if v.SpriteAnimation.lastImage == nil {
			continue
		}
		if v.Sprite.Image == v.SpriteAnimation.lastImage {
			continue
		}
		v.Sprite.Image = v.SpriteAnimation.lastImage
		offx, offy := v.SpriteAnimation.Anim.GetClipOffset(v.SpriteAnimation.ActiveClip, v.SpriteAnimation.ActiveFrame)
		v.Sprite.OffsetX = offx
		v.Sprite.OffsetY = offy
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
	if evs := spranim.Anim.GetClipEvents(index); evs != nil && len(evs) > 0 && evs[0] != nil {
		spranim.AnimEvent(evs[0].Name, evs[0].Value)
	}
}

func spriteAnimResolveClipMap(spranim *SpriteAnimation) {
	if spranim.clipMapLen == spranim.Anim.Count() {
		return
	}
	// rebuild cache
	spranim.clipMap = make(map[string]int)
	spranim.Anim.Each(func(i int, clip AnimationClip) bool {
		spranim.clipMap[clip.GetName()] = i
		return true
	})
	spranim.clipMapLen = spranim.Anim.Count()
}

func spriteAnimResolvePlayback(globalfps, dt float64, spranim *SpriteAnimation) {
	clip := spranim.Anim.GetClip(spranim.ActiveClip)
	localfps := nonzeroval(clip.GetFPS(), spranim.Fps, globalfps)
	localdt := dt //(dt * localfps) / globalfps
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
	spranim.lastImage = spranim.Anim.GetClipImage(spranim.lastClip, spranim.ActiveFrame)
	spranim.lastPlaying = true
	spranim.T += localdt * localfps
	if spranim.T >= 1 {
		// next frame
		nextframe := spranim.ActiveFrame + 1
		if spranim.reversed {
			nextframe = spranim.ActiveFrame - 1
		}
		if nextframe >= clip.GetFrameCount() {
			// animation ended
			switch clip.GetMode() {
			case AnimOnce:
				spranim.T = 0
				spranim.Playing = false
				if e := clip.GetEndedEvent(); e != nil {
					spranim.AnimEvent(e.Name, e.Value)
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
			if evs := clip.GetEvents(); evs != nil && len(evs) > nextframe && evs[nextframe] != nil {
				spranim.AnimEvent(evs[nextframe].Name, evs[nextframe].Value)
			}
		}
	}
}

// type SpriteAnimationLinkComponentSystem struct {
// 	BaseComponentSystem
// }

// func (cs *SpriteAnimationLinkComponentSystem) SystemName() string {
// 	return SNSpriteAnimationLink
// }

// func (cs *SpriteAnimationLinkComponentSystem) SystemPriority() int {
// 	return -5
// }

// func (cs *SpriteAnimationLinkComponentSystem) SystemExec() SystemExecFn {
// 	return SpriteAnimationLinkSystemExec
// }

// func (cs *SpriteAnimationLinkComponentSystem) SystemTags() []string {
// 	return []string{"draw"}
// }

// func (cs *SpriteAnimationLinkComponentSystem) Components(w *ecs.World) []*ecs.Component {
// 	return []*ecs.Component{
// 		spriteAnimationComponentDef(w),
// 		drawableComponentDef(w),
// 	}
// }

// func (cs *SpriteAnimationLinkComponentSystem) ExcludeComponents(w *ecs.World) []*ecs.Component {
// 	return emptyCompSlice
// }

// // SpriteAnimationLinkSystemExec is what glues the animation and sprite together
// func SpriteAnimationLinkSystemExec(ctx Context) {
// 	//screen := ctx.Screen()
// 	v := ctx.System().View()
// 	world := ctx.World()
// 	matches := v.Matches()
// 	spriteanimcomp := world.Component(CNSpriteAnimation)
// 	spritecomp := world.Component(CNDrawable)
// 	for _, m := range matches {
// 		spranim := m.Components[spriteanimcomp].(*SpriteAnimation)
// 		spr := m.Components[spritecomp].(Drawable)
// 		if !spranim.Enabled {
// 			continue
// 		}
// 		// replace image if the animation clip uses a different one
// 		if spranim.lastImage != nil {
// 			if w, ok := spr.(DrawableImager); ok {
// 				if w.GetImage() != spranim.lastImage {
// 					w.SetImage(spranim.lastImage)
// 					spr.SetOffset(spranim.Anim.GetClipOffset(spranim.ActiveClip, spranim.ActiveFrame))
// 				}
// 			}
// 		}
// 	}
// }

// func init() {
// 	RegisterComponentSystem(&SpriteAnimationComponentSystem{})
// 	RegisterComponentSystem(&SpriteAnimationLinkComponentSystem{})
// }

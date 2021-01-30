package graphics

import (
	"image"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
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

//    ________         _   _____      _          _
//   |_   __  |       / |_|_   _|    (_)        / |_
//     | |_ \_|_   __`| |-' | |      __   .--. `| |-'.---.  _ .--.  .---.  _ .--.
//     |  _| _[ \ [  ]| |   | |   _ [  | ( (`\] | | / /__\\[ `.-. |/ /__\\[ `/'`\]
//    _| |__/ |\ \/ / | |, _| |__/ | | |  `'.'. | |,| \__., | | | || \__., | |
//   |________| \__/  \__/|________|[___][\__) )\__/ '.__.'[___||__]'.__.'[___]

// AnimationEventID is the id oa an animation listener instance
type AnimationEventID int64

// AnimationEventListeners groups AnimationEvent listeners
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

// Add a new listener
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

// AddCatchAll adds a wildcard listener
func (ae *AnimationEventListeners) AddCatchAll(fn AnimationEventFn) AnimationEventID {
	return ae.Add("*", fn)
}

// Remove a listener by ID
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

// Clear removes all event listeners
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

// Dispatch triggers all compatible listeners with the event
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

// AnimationEventFn is a valid AnimationEvent function
type AnimationEventFn func(name, value string)

// AnimationEvent holds a name and value to pass to listeners when
// the event is triggered
type AnimationEvent struct {
	Name  string
	Value string
}

// .-. . . .-. .-. .-. .-. .-. .-. .-. .-.
//  |  |\|  |  |-  |(  |-  |-| |   |-  `-.
// `-' ' `  '  `-' ' ' '   ` ' `-' `-' `-'

// Animation is a shared container of clips (reusable)
type Animation interface {
	Each(fn func(i int, clip AnimationClip) bool)
	GetClip(index int) AnimationClip
	GetClipEvents(index int) []*AnimationEvent
	GetClipImage(clipindex, frame int) *ebiten.Image
	GetClipRect(clipindex, frame int) image.Rectangle
	GetClipOffset(clipindex, frame int) (x, y float64)
	Count() int
}

// AnimationClip is the data os an animation clip (reusable)
type AnimationClip interface {
	GetName() string
	GetFPS() float64
	GetImage(frame int) *ebiten.Image
	GetFrameCount() int
	GetMode() AnimClipMode
	GetEndedEvent() *AnimationEvent
	GetEvents() []*AnimationEvent
	GetRect(frame int) image.Rectangle
	GetOffset(frame int) (x, y float64)
}

//████████╗██╗██╗     ███████╗██████╗
//╚══██╔══╝██║██║     ██╔════╝██╔══██╗
//   ██║   ██║██║     █████╗  ██║  ██║
//   ██║   ██║██║     ██╔══╝  ██║  ██║
//   ██║   ██║███████╗███████╗██████╔╝
//   ╚═╝   ╚═╝╚══════╝╚══════╝╚═════╝

// TiledAnimation is an animation of a single
// image source and many Rectangles that represent views (subimages)
type TiledAnimation struct {
	Clips []TiledAnimationClip
}

// GetClip returns an animation clip by index
func (a *TiledAnimation) GetClip(index int) AnimationClip {
	if a.Clips == nil {
		return nil
	}
	if len(a.Clips) <= index || index < 0 {
		return nil
	}
	return a.Clips[index]
}

// GetClipEvents returns all animation clip events by the clip index
func (a *TiledAnimation) GetClipEvents(index int) []*AnimationEvent {
	clip := a.GetClip(index)
	if clip == nil {
		return nil
	}
	return clip.GetEvents()
}

// GetClipImage returns the clip image
func (a *TiledAnimation) GetClipImage(clipi, frame int) *ebiten.Image {
	clip := a.GetClip(clipi)
	if clip == nil {
		return nil
	}
	return clip.GetImage(frame)
}

// GetClipRect returns the clip Rectangle (bounds)
func (a *TiledAnimation) GetClipRect(clipi, frame int) image.Rectangle {
	clip := a.GetClip(clipi)
	if clip == nil {
		return image.Rectangle{}
	}
	return clip.GetRect(frame)
}

// GetClipOffset returns the clip offset
func (a *TiledAnimation) GetClipOffset(clipi, frame int) (x, y float64) {
	clip := a.GetClip(clipi)
	if clip == nil {
		return 0, 0
	}
	return clip.GetOffset(frame)
}

// Count returns the total count of cnimation clips
func (a *TiledAnimation) Count() int {
	return len(a.Clips)
}

// Each iterates through all animation clips
func (a *TiledAnimation) Each(fn func(i int, clip AnimationClip) bool) {
	for i, v := range a.Clips {
		if !fn(i, v) {
			return
		}
	}
}

// TiledAnimationClip is an animation clip, like a character walk cycle.
type TiledAnimationClip struct {
	// The name of an animation is not allowed to be changed during runtime
	// but since this is part of a component (and components shouldn't have logic),
	// it is a public member.
	Name       string
	Image      *ebiten.Image
	Frames     []image.Rectangle
	Events     []*AnimationEvent //TODO: link
	Fps        float64
	ClipMode   AnimClipMode
	EndedEvent *AnimationEvent //TODO: link
}

// GetName returns the animation clip name
func (c TiledAnimationClip) GetName() string {
	return c.Name
}

func (c TiledAnimationClip) GetFPS() float64 {
	return c.Fps
}

func (c TiledAnimationClip) GetImage(frame int) *ebiten.Image {
	return c.Image.SubImage(c.Frames[frame]).(*ebiten.Image)
}

func (c TiledAnimationClip) GetFrameCount() int {
	return len(c.Frames)
}

func (c TiledAnimationClip) GetMode() AnimClipMode {
	return c.ClipMode
}

func (c TiledAnimationClip) GetEvents() []*AnimationEvent {
	return c.Events
}

func (c TiledAnimationClip) GetEndedEvent() *AnimationEvent {
	return c.EndedEvent
}

// GetRect returns the drawable bounds
func (c TiledAnimationClip) GetRect(frame int) image.Rectangle {
	if c.Frames != nil && len(c.Frames) > frame && frame >= 0 {
		return c.Frames[frame]
	}
	return image.Rectangle{}
}

// GetOffset is anoop for TiledAnimationClip
func (c TiledAnimationClip) GetOffset(frame int) (x, y float64) {
	// noop for TiledAnimationClip
	return 0, 0
}

// ██████╗ ██████╗ ███████╗ ██████╗ ██████╗ ███╗   ███╗██████╗ ██╗   ██╗████████╗███████╗██████╗
// ██╔══██╗██╔══██╗██╔════╝██╔════╝██╔═══██╗████╗ ████║██╔══██╗██║   ██║╚══██╔══╝██╔════╝██╔══██╗
// ██████╔╝██████╔╝█████╗  ██║     ██║   ██║██╔████╔██║██████╔╝██║   ██║   ██║   █████╗  ██║  ██║
// ██╔═══╝ ██╔══██╗██╔══╝  ██║     ██║   ██║██║╚██╔╝██║██╔═══╝ ██║   ██║   ██║   ██╔══╝  ██║  ██║
// ██║     ██║  ██║███████╗╚██████╗╚██████╔╝██║ ╚═╝ ██║██║     ╚██████╔╝   ██║   ███████╗██████╔╝
// ╚═╝     ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝      ╚═════╝    ╚═╝   ╚══════╝╚═════╝

// PrecomputedAnimation is an animation of a one or more image sources
// It is the default data structure of an Atlas
type PrecomputedAnimation struct {
	Clips []PcAnimClip
}

// GetClip returns an animation clip by index
func (a *PrecomputedAnimation) GetClip(index int) AnimationClip {
	if a.Clips == nil {
		return nil
	}
	if len(a.Clips) <= index || index < 0 {
		return nil
	}
	return a.Clips[index]
}

// GetClipEvents returns all animation clip events by the clip index
func (a *PrecomputedAnimation) GetClipEvents(index int) []*AnimationEvent {
	clip := a.GetClip(index)
	if clip == nil {
		return nil
	}
	return clip.GetEvents()
}

// GetClipImage returns the clip image
func (a *PrecomputedAnimation) GetClipImage(clipi, frame int) *ebiten.Image {
	clip := a.GetClip(clipi)
	if clip == nil {
		return nil
	}
	return clip.GetImage(frame)
}

// GetClipRect returns the clip Rectangle (bounds)
func (a *PrecomputedAnimation) GetClipRect(clipi, frame int) image.Rectangle {
	clip := a.GetClip(clipi)
	if clip == nil {
		return image.Rectangle{}
	}
	return clip.GetRect(frame)
}

// GetClipOffset returns the clip offset
func (a *PrecomputedAnimation) GetClipOffset(clipi, frame int) (x, y float64) {
	clip := a.GetClip(clipi)
	if clip == nil {
		return 0, 0
	}
	return clip.GetOffset(frame)
}

// Count returns the total count of cnimation clips
func (a *PrecomputedAnimation) Count() int {
	return len(a.Clips)
}

// Each iterates through all animation clips
func (a *PrecomputedAnimation) Each(fn func(i int, clip AnimationClip) bool) {
	for i, v := range a.Clips {
		if !fn(i, v) {
			return
		}
	}
}

type PcFrame struct {
	Image   *ebiten.Image
	Rect    image.Rectangle
	OffsetX float64
	OffsetY float64
}

// PcAnimClip is a pre-computed animation clip.
type PcAnimClip struct {
	// The name of an animation is not allowed to be changed during runtime
	// but since this is part of a component (and components shouldn't have logic),
	// it is a public member.
	Name       string
	Frames     []PcFrame
	Events     []*AnimationEvent //TODO: link
	Fps        float64
	ClipMode   AnimClipMode
	EndedEvent *AnimationEvent //TODO: link
}

// GetName returns the animation clip name
func (c PcAnimClip) GetName() string {
	return c.Name
}

func (c PcAnimClip) GetFPS() float64 {
	return c.Fps
}

func (c PcAnimClip) GetImage(frame int) *ebiten.Image {
	if c.Frames != nil && len(c.Frames) > frame && frame >= 0 {
		return c.Frames[frame].Image
	}
	return nil
}

func (c PcAnimClip) GetFrameCount() int {
	return len(c.Frames)
}

func (c PcAnimClip) GetMode() AnimClipMode {
	return c.ClipMode
}

func (c PcAnimClip) GetEvents() []*AnimationEvent {
	return c.Events
}

func (c PcAnimClip) GetEndedEvent() *AnimationEvent {
	return c.EndedEvent
}

// GetRect returns the drawable bounds
func (c PcAnimClip) GetRect(frame int) image.Rectangle {
	if c.Frames != nil && len(c.Frames) > frame && frame >= 0 {
		return c.Frames[frame].Rect
	}
	return image.Rectangle{}
}

// GetOffset returns the pixel offset (from the origin) of the current frame
func (c PcAnimClip) GetOffset(frame int) (x, y float64) {
	if c.Frames != nil && len(c.Frames) > frame && frame >= 0 {
		return c.Frames[frame].OffsetX, c.Frames[frame].OffsetY
	}
	return 0, 0
}

package aseprite

import (
	"encoding/json"
	"image"
	"strconv"
)

// AnimDirection represents Aseprite animation directions
type AnimDirection string

const (
	// AnimForward forward animation
	AnimForward AnimDirection = "forward"
	// AnimReverse reverse animation
	AnimReverse AnimDirection = "reverse"
	// AnimPingPong goes forward and backwards endlessly
	AnimPingPong AnimDirection = "pingpong"
)

// File represents an Aseprite sprite sheet file
type File struct {
	Frames []FrameInfo `json:"frames"`
	Meta   Metadata    `json:"meta"`
}

// Length returns the amount of frames
func (f *File) Length() int {
	return len(f.Frames)
}

// Walk iterates on every FrameInfo. It stops early if fn returns false.
func (f *File) Walk(fn func(i FrameInfo) bool) {
	for _, v := range f.Frames {
		if !fn(v) {
			return
		}
	}
}

// GetMetadata retrieves Aseprite metadata
func (f *File) GetMetadata() Metadata {
	return f.Meta
}

// GetFrameByIndex returns the FrameInfo at the index position. It returns false
// if out of bounds.
func (f *File) GetFrameByIndex(index int) (i FrameInfo, ok bool) {
	if index < 0 {
		return
	}
	if len(f.Frames) <= index {
		return
	}
	i = f.Frames[index]
	ok = true
	return
}

// GetFrameByName returns the FrameInfo with the specified name.
func (f *File) GetFrameByName(name string) (i FrameInfo, ok bool) {
	for _, v := range f.Frames {
		if v.Filename == name {
			return v, true
		}
	}
	return
}

type FrameInfo struct {
	Filename         string    `json:"filename"`
	Frame            FrameRect `json:"frame"`
	Rotated          bool      `json:"rotated"`
	Trimmed          bool      `json:"trimmed"`
	SpriteSourceSize FrameRect `json:"spriteSourceSize"`
	SourceSize       ImSize    `json:"sourceSize"`
	Duration         int       `json:"duration"`
}

// FrameRect is the frame bounds
type FrameRect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

func (r FrameRect) String() string {
	return "FrameRect{x:" + strconv.Itoa(r.X) + ",y:" + strconv.Itoa(r.Y) + ",w:" + strconv.Itoa(r.W) + ",h:" + strconv.Itoa(r.H) + "}"
}

func (r FrameRect) ToRect() image.Rectangle {
	return image.Rect(r.X, r.Y, r.X+r.W, r.Y+r.H)
}

type ImSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Vec2 struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Metadata struct {
	App       string     `json:"app"`
	Version   string     `json:"version"`
	Image     string     `json:"image"`
	Format    string     `json:"format"`
	Size      ImSize     `json:"size,omitempty"`
	Scale     string     `json:"scale"`
	FrameTags []FrameTag `json:"frameTags"`
	Layers    []Layer    `json:"layers,omitempty"`
	Slices    []Slice    `json:"slices"`
}

type Layer struct {
	Name      string  `json:"name"`
	Opacity   float64 `json:"opacity"`
	BlendMode string  `json:"blendMode"`
}

// "name": "lbar1", "from": 0, "to": 11, "direction": "forward"
type FrameTag struct {
	Name      string        `json:"name"`
	From      int           `json:"from"`
	To        int           `json:"to"`
	Direction AnimDirection `json:"direction"`
}

type Slice struct {
	Name  string          `json:"name"`
	Color string          `json:"color"`
	Keys  []SliceKeyframe `json:"keys"`
}

type SliceKeyframe struct {
	Frame  int       `json:"frame"`
	Bounds FrameRect `json:"bounds"`
	Pivot  Vec2      `json:"pivot"`
}

// Parse an Aseprite sheet JSON file. Warning: export frams as ARRAY. Do not use
// the Map option. JSON maps are not guaranteed to be ordered, and Primen is made
// in Go, which doesn't preserve the order of a Map.
// The JSON spec states that relying on key ordered maps is abad idea.
// https://github.com/golang/go/issues/27179#issuecomment-415525033
func Parse(jsonb []byte) (*File, error) {
	m := &File{}
	err := json.Unmarshal(jsonb, m)
	return m, err
}

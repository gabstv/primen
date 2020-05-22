package aseprite

import (
	"encoding/json"
	"image"
	"strconv"
)

type FileType string

const (
	FileTypeSlice FileType = "slice"
	FileTypeMap   FileType = "map"
)

type AnimDirection string

const (
	AnimForward  AnimDirection = "forward"
	AnimReverse  AnimDirection = "reverse"
	AnimPingPong AnimDirection = "pingpong"
)

type File interface {
	Walk(fn func(i FrameInfo) bool)
	Type() FileType
	GetMetadata() Metadata
	GetFrame(name string) (i FrameInfo, ok bool)
}

type FileMap struct {
	Frames map[string]FrameInfo `json:"frames"`
	Meta   Metadata             `json:"meta"`
}

func (f *FileMap) Walk(fn func(i FrameInfo) bool) {
	for k, v := range f.Frames {
		v2 := v
		v2.Filename = k
		if !fn(v2) {
			return
		}
	}
}

func (f *FileMap) Type() FileType {
	return FileTypeMap
}

func (f *FileMap) GetMetadata() Metadata {
	return f.Meta
}

func (f *FileMap) GetFrame(name string) (i FrameInfo, ok bool) {
	i, ok = f.Frames[name]
	return
}

type FileSlice struct {
	Frames []FrameInfo `json:"frames"`
	Meta   Metadata    `json:"meta"`
}

func (f *FileSlice) Walk(fn func(i FrameInfo) bool) {
	for _, v := range f.Frames {
		if !fn(v) {
			return
		}
	}
}

func (f *FileSlice) Type() FileType {
	return FileTypeSlice
}

func (f *FileSlice) GetMetadata() Metadata {
	return f.Meta
}

func (f *FileSlice) GetFrame(name string) (i FrameInfo, ok bool) {
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
	App       string      `json:"app"`
	Version   string      `json:"version"`
	Image     string      `json:"image"`
	Format    string      `json:"format"`
	Size      ImSize      `json:"size,omitempty"`
	Scale     string      `json:"scale"`
	FrameTags []FrameTag  `json:"frameTags"`
	Layers    []Layer     `json:"layers,omitempty"`
	Slices    interface{} `json:"slices"`
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

func Parse(jsonb []byte) (File, error) {
	m, merr := ParseMap(jsonb)
	s, serr := ParseSlice(jsonb)
	if merr != nil && serr != nil {
		return nil, merr
	}
	if m != nil {
		if s != nil {
			if len(s.Frames) > len(m.Frames) {
				return s, nil
			}
		}
		return m, nil
	}
	return s, nil
}

func ParseMap(jsonb []byte) (*FileMap, error) {
	m := &FileMap{}
	err := json.Unmarshal(jsonb, m)
	return m, err
}

func ParseSlice(jsonb []byte) (*FileSlice, error) {
	m := &FileSlice{}
	err := json.Unmarshal(jsonb, m)
	return m, err
}

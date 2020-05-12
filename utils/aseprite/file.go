package aseprite

import "encoding/json"

type FileType string

const (
	FileTypeSlice FileType = "slice"
	FileTypeMap   FileType = "map"
)

type File interface {
	Walk(fn func(i FrameInfo) bool)
	Type() FileType
	GetMetadata() Metadata
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

type ImSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Metadata struct {
	App       string      `json:"app"`
	Version   string      `json:"version"`
	Image     string      `json:"image"`
	Format    string      `json:"format"`
	Size      ImSize      `json:"size,omitempty"`
	Scale     string      `json:"scale"`
	FrameTags interface{} `json:"frameTags"`
	Layers    []Layer     `json:"layers,omitempty"`
	Slices    interface{} `json:"slices"`
}

type Layer struct {
	Name      string  `json:"name"`
	Opacity   float64 `json:"opacity"`
	BlendMode string  `json:"blendMode"`
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

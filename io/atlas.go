package io

import (
	"bytes"
	"image"
	"image/png"
	"sort"

	"github.com/gabstv/primen/components/graphics"
	"github.com/gabstv/primen/io/pb"
	proto "github.com/golang/protobuf/proto"
	"github.com/hajimehoshi/ebiten"
)

type Atlas struct {
	ebimg     []*ebiten.Image
	frames    map[string]*Sprite
	anims     map[string]*graphics.PrecomputedAnimation
	animClips map[string]graphics.PcAnimClip
}

type Sprite struct {
	Name  string
	Image *ebiten.Image
	Pivot image.Point
}

func (a *Atlas) GetSubImage(name string) *Sprite {
	return a.frames[name]
}

func (a *Atlas) GetSubImages() []*Sprite {
	out := make([]*Sprite, 0, len(a.frames))
	for _, v := range a.frames {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func (a *Atlas) GetAnimClip(name string) graphics.AnimationClip {
	return a.animClips[name]
}

func (a *Atlas) GetAnimClips() []ClipG {
	out := make([]ClipG, 0, len(a.animClips))
	for k, v := range a.animClips {
		out = append(out, ClipG{
			Name:     k,
			AnimClip: v,
		})
	}
	return out
}

func (a *Atlas) GetAnimation(name string) graphics.Animation {
	return a.anims[name]
}

func (a *Atlas) GetAnimations() []AnimG {
	out := make([]AnimG, 0, len(a.anims))
	for k, v := range a.anims {
		out = append(out, AnimG{
			Name: k,
			Anim: v,
		})
	}
	return out
}

func ParseAtlas(b []byte) (*Atlas, error) {
	src := &pb.AtlasFile{}
	if err := proto.Unmarshal(b, src); err != nil {
		return nil, err
	}
	imgs := make([]*ebiten.Image, len(src.Images))
	for i, v := range src.Images {
		rdr := bytes.NewReader(v)
		im, err := png.Decode(rdr)
		if err != nil {
			return nil, err
		}
		ei, err := ebiten.NewImageFromImage(im, pb.ToEbitenFilter(src.Filters[i]))
		if err != nil {
			return nil, err
		}
		imgs[i] = ei
	}
	// frames!
	frames := make(map[string]*Sprite)
	for k, v := range src.Frames {
		im := imgs[int(v.Image)]
		simg := im.SubImage(image.Rect(int(v.X), int(v.Y), int(v.X+v.W), int(v.Y+v.H))).(*ebiten.Image)
		frames[k] = &Sprite{
			Name:  k,
			Image: simg,
			Pivot: image.Point{
				X: int(v.Ox),
				Y: int(v.Oy),
			},
		}
	}
	// anim clips
	clips := make(map[string]graphics.PcAnimClip)
	for k, v := range src.Clips {
		cl := importAnimClip(k, v, frames)
		clips[cl.Name] = cl
	}
	// anim groups
	animgs := make(map[string]*graphics.PrecomputedAnimation)
	for k, v := range src.Animations {
		anim := &graphics.PrecomputedAnimation{
			Clips: make([]graphics.PcAnimClip, 0),
		}
		for _, clipv := range v.Clips {
			cc := importAnimClip(clipv.Name, clipv, frames)
			anim.Clips = append(anim.Clips, cc)
		}
		animgs[k] = anim
	}
	// all set
	return &Atlas{
		ebimg:     imgs,
		frames:    frames,
		animClips: clips,
		anims:     animgs,
	}, nil
}

func importAnimClip(name string, v *pb.AnimationClip, frames map[string]*Sprite) graphics.PcAnimClip {
	cl := graphics.PcAnimClip{
		Name:   name,
		Fps:    float64(v.Fps),
		Events: make([]*graphics.AnimationEvent, 0),
		Frames: make([]graphics.PcFrame, 0),
	}
	if v.EndedEvent != nil {
		cl.EndedEvent = &graphics.AnimationEvent{
			Name:  v.EndedEvent.Name,
			Value: v.EndedEvent.Value,
		}
	}
	switch v.ClipMode {
	case pb.AnimationClipMode_PING_PONG:
		cl.ClipMode = graphics.AnimPingPong
	case pb.AnimationClipMode_ONCE:
		cl.ClipMode = graphics.AnimOnce
	case pb.AnimationClipMode_LOOP:
		cl.ClipMode = graphics.AnimLoop
	case pb.AnimationClipMode_CLAMP_FOREVER:
		cl.ClipMode = graphics.AnimClampForever
	}
	for _, vf := range v.Frames {
		realf := frames[vf.FrameName]
		if realf == nil {
			panic("io frame not found: " + vf.FrameName)
		}
		if vf.Event == nil {
			cl.Events = append(cl.Events, nil)
		} else {
			cl.Events = append(cl.Events, &graphics.AnimationEvent{
				Name:  vf.Event.Name,
				Value: vf.Event.Value,
			})
		}
		sz := realf.Image.Bounds().Size()
		cf := graphics.PcFrame{
			OffsetX: float64(realf.Pivot.X),
			OffsetY: float64(realf.Pivot.Y),
			Rect:    image.Rect(0, 0, sz.X, sz.Y),
			Image:   realf.Image,
		}
		cl.Frames = append(cl.Frames, cf)
	}
	return cl
}

type SubImageG struct {
	Name     string
	SubImage *ebiten.Image
}

type ClipG struct {
	Name     string
	AnimClip graphics.AnimationClip
}

type AnimG struct {
	Name string
	Anim graphics.Animation
}

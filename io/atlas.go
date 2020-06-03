package io

import (
	"bytes"
	"image"
	"image/png"

	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/io/pb"
	proto "github.com/golang/protobuf/proto"
	"github.com/hajimehoshi/ebiten"
)

type Atlas struct {
	ebimg     []*ebiten.Image
	frames    map[string]*ebiten.Image
	anims     map[string]*core.PrecomputedAnimation
	animClips map[string]core.PcAnimClip
}

func (a *Atlas) GetSubImage(name string) *ebiten.Image {
	return a.frames[name]
}

func (a *Atlas) GetSubImages() []SubImageG {
	out := make([]SubImageG, 0, len(a.frames))
	for k, v := range a.frames {
		out = append(out, SubImageG{
			Name:     k,
			SubImage: v,
		})
	}
	return out
}

func (a *Atlas) GetAnimClip(name string) core.AnimationClip {
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

func (a *Atlas) GetAnimation(name string) core.Animation {
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
	frames := make(map[string]*ebiten.Image)
	for k, v := range src.Frames {
		im := imgs[int(v.Image)]
		simg := im.SubImage(image.Rect(int(v.X), int(v.Y), int(v.X+v.W), int(v.Y+v.H))).(*ebiten.Image)
		frames[k] = simg
	}
	// anim clips
	clips := make(map[string]core.PcAnimClip)
	for k, v := range src.AnimClips {
		cl := importAnimClip(k, v, imgs)
		clips[cl.Name] = cl
	}
	// anim groups
	animgs := make(map[string]*core.PrecomputedAnimation)
	for k, v := range src.AnimGroups {
		anim := &core.PrecomputedAnimation{
			Clips: make([]core.PcAnimClip, 0),
		}
		for global, local := range v.Clips {
			cc := clips[global]
			cc.Name = local
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

func importAnimClip(name string, v *pb.AnimationClip, imgs []*ebiten.Image) core.PcAnimClip {
	cl := core.PcAnimClip{
		Name:   name,
		Fps:    float64(v.Fps),
		Events: make([]*core.AnimationEvent, 0),
		Frames: make([]core.PcFrame, 0),
	}
	if v.EndedEvent != nil {
		cl.EndedEvent = &core.AnimationEvent{
			Name:  v.EndedEvent.Name,
			Value: v.EndedEvent.Value,
		}
	}
	switch v.ClipMode {
	case pb.AnimationClipMode_PING_PONG:
		println("TODO: implement pb.AnimationClipMode_PING_PONG")
		//TODO: implement pb.AnimationClipMode_PING_PONG
	case pb.AnimationClipMode_ONCE:
		cl.ClipMode = core.AnimOnce
	case pb.AnimationClipMode_LOOP:
		cl.ClipMode = core.AnimLoop
	case pb.AnimationClipMode_CLAMP_FOREVER:
		cl.ClipMode = core.AnimClampForever
	}
	for _, vf := range v.Frames {
		if vf.Event == nil {
			cl.Events = append(cl.Events, nil)
		} else {
			cl.Events = append(cl.Events, &core.AnimationEvent{
				Name:  vf.Event.Name,
				Value: vf.Event.Value,
			})
		}
		cf := core.PcFrame{
			OffsetX: float64(vf.Ox),
			OffsetY: float64(vf.Oy),
			Rect:    image.Rect(0, 0, int(vf.W), int(vf.H)),
			Image:   imgs[int(vf.Image)].SubImage(image.Rect(int(vf.X), int(vf.Y), int(vf.X+vf.W), int(vf.X+vf.H))).(*ebiten.Image),
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
	AnimClip core.AnimationClip
}

type AnimG struct {
	Name string
	Anim core.Animation
}

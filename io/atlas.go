package io

import (
	"bytes"
	"image"
	"image/png"

	"github.com/gabstv/tau/io/pb"
	proto "github.com/golang/protobuf/proto"
	"github.com/hajimehoshi/ebiten"
)

type Atlas struct {
	ebimg  []*ebiten.Image
	frames map[string]*ebiten.Image
}

func (a *Atlas) Get(name string) *ebiten.Image {
	return a.frames[name]
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
	// all set
	return &Atlas{
		ebimg:  imgs,
		frames: frames,
	}, nil
}

package aseprite

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gabstv/primen/internal/atlaspacker"
	"github.com/gabstv/primen/io/pb"
)

type ImportInput struct {
	Template *AtlasImporterGroup
	Source   []AsepriteInput
	PackerI  atlaspacker.PackerInput
}

type AsepriteInput struct {
	Filename  string
	FrameData *File
	ImageData []byte
}

func Import(ctx context.Context, input ImportInput) (*pb.AtlasFile, error) {
	outputf := &pb.AtlasFile{
		Images:     make([][]byte, 0),
		Filters:    make([]pb.ImageFilter, 0),
		Frames:     make(map[string]*pb.Frame),
		Clips:      make(map[string]*pb.AnimationClip),
		Animations: make(map[string]*pb.Animation),
	}
	impkr := newImImporter()
	pkr := &atlaspacker.BinTreeRectPacker{}
	for _, tpl := range input.Template.Templates {
		for _, ips := range input.Source {
			if ips.Filename == tpl.AsepriteSheet {
				if err := importAtlas(ctx, tpl, ips, pkr, impkr); err != nil {
					return outputf, err
				}
				break
			}
		}
	}
	file, err := buildAtlas(ctx, buildAtlasInput{
		PackerI: input.PackerI,
		Filter:  input.Template.ImageFilter,
		Im:      impkr,
		Pkr:     pkr,
		Clips:   input.Template.Clips,
		Anims:   input.Template.Animations,
	})
	if err != nil {
		return nil, err
	}
	return file, nil
}

func importAtlas(ctx context.Context, tpl AtlasImporter, ase AsepriteInput, pkr *atlaspacker.BinTreeRectPacker, imptr *imImporter) error {
	img, err := png.Decode(bytes.NewReader(ase.ImageData))
	if err != nil {
		return fmt.Errorf("error decoding png image: %w", err)
	}
	//
	rules := importByRules{
		Template:  tpl,
		Filename:  ase.Filename,
		FrameData: ase.FrameData,
		Img:       img,
		Pkr:       pkr,
		Imptr:     imptr,
	}
	return importAtlasByFrames(ctx, rules)
}

type importByRules struct {
	Template  AtlasImporter
	Filename  string
	FrameData *File
	Img       image.Image
	Pkr       *atlaspacker.BinTreeRectPacker
	Imptr     *imImporter
}

func importAtlasByFrames(ctx context.Context, r importByRules) error {
	idm := make(map[image.Image]*atlaspacker.RectPackerNode)
	//r.Template.
	imbank := &atlasIMCache{}
	r.FrameData.Walk(func(i FrameInfo) bool {
		if frame, ok := r.Template.FrameWithFilename(i.Filename); ok {
			clipim := imbank.getSubImage(r.Img, i.Frame)
			if _, ok := idm[clipim]; !ok {
				node := r.Pkr.AddRect(clipim.Bounds())
				idm[clipim] = node
				r.Imptr.addSprite(i.Filename, node, clipim, frame.Pivot)
			} else {
				r.Imptr.addSprite(i.Filename, idm[clipim], clipim, frame.Pivot)
			}
		} else if r.Template.ExportUndefinedFrames {
			clipim := imbank.getSubImage(r.Img, i.Frame)
			if _, ok := idm[clipim]; !ok {
				node := r.Pkr.AddRect(clipim.Bounds())
				idm[clipim] = node
				r.Imptr.addSprite(i.Filename, node, clipim, Vec2{})
			} else {
				r.Imptr.addSprite(i.Filename, idm[clipim], clipim, Vec2{})
			}
		}
		return true
	})
	return nil
}

func getFilter(v string) pb.ImageFilter {
	v = strings.TrimSpace(strings.ToLower(v))
	switch v {
	case "pixel", "nn", "nearest":
		return pb.ImageFilter_NEAREST
	case "linear":
		return pb.ImageFilter_LINEAR
	}
	return pb.ImageFilter_DEFAULT
}

type buildAtlasInput struct {
	PackerI atlaspacker.PackerInput
	Pkr     *atlaspacker.BinTreeRectPacker
	Im      *imImporter
	Filter  string
	Clips   []AnimationClip
	Anims   []Animation
}

func buildAtlas(ctx context.Context, input buildAtlasInput) (*pb.AtlasFile, error) {
	pki := input.PackerI
	pkr := input.Pkr
	im := input.Im
	file := &pb.AtlasFile{
		Images:     make([][]byte, 0),
		Filters:    make([]pb.ImageFilter, 0),
		Frames:     make(map[string]*pb.Frame),
		Clips:      make(map[string]*pb.AnimationClip),
		Animations: make(map[string]*pb.Animation),
	}
	if pki.MaxWidth <= 0 {
		pki.MaxWidth = 4096
	}
	if pki.MaxHeight <= 0 {
		pki.MaxHeight = 4096
	}
	atlases, err := pkr.Pack(ctx, pki)
	if err != nil {
		return nil, err
	}
	xsprites := im.Sprites()
	//anims := make(map[string]*pb.AnimationClip)
	for i, atlas := range atlases {
		nodemap := make(map[*atlaspacker.RectPackerNode]int)
		for i, v := range atlas.Nodes {
			nodemap[v] = i
		}
		outimg := image.NewRGBA(image.Rect(0, 0, atlas.Width, atlas.Height))
		frames := make(map[string]*pb.Frame)
		for _, spr := range xsprites {
			_, nodeok := nodemap[spr.Node]
			if !nodeok {
				continue
			}
			spr.FinalRect = FrameRect{
				X: spr.Node.X,
				Y: spr.Node.Y,
				W: spr.Node.Width,
				H: spr.Node.Height,
			}
			spr.FinalImageIndex = i
			//TODO: prevent from drawing twice (if nn was already done)
			if pki.Debug {
				dbuf := new(bytes.Buffer)
				_ = png.Encode(dbuf, spr.Image)
				ioutil.WriteFile("xsprites_"+spr.Name+".png", dbuf.Bytes(), 0644)
			}
			draw.Draw(outimg, spr.Node.R(), spr.Image, image.ZP, draw.Over)
			frames[spr.Name] = &pb.Frame{
				Image: uint32(i),
				X:     uint32(spr.Node.X),
				Y:     uint32(spr.Node.Y),
				W:     uint32(spr.Node.Width),
				H:     uint32(spr.Node.Height),
				Ox:    int32(spr.Pivot.X * -1),
				Oy:    int32(spr.Pivot.Y * -1),
			}
		}
		file.Filters = append(file.Filters, getFilter(input.Filter))
		for fn, frame := range frames {
			file.Frames[fn] = frame
		}
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, outimg); err != nil {
			return nil, err
		}
		file.Images = append(file.Images, buf.Bytes())
		if pki.Debug {
			ioutil.WriteFile("atlas_"+strconv.Itoa(i)+".png", buf.Bytes(), 0644)
		}
	}
	// put animations and solo clips
	for _, clip := range input.Clips {
		pbclip, err := getClip(file, clip)
		if err != nil {
			return nil, err
		}
		file.Clips[pbclip.Name] = pbclip
	}
	for _, anim := range input.Anims {
		pbanim := &pb.Animation{
			Name:  anim.Name,
			Clips: make([]*pb.AnimationClip, 0, len(anim.Clips)),
		}
		for _, clip := range anim.Clips {
			pbclip, err := getClip(file, clip)
			if err != nil {
				return nil, err
			}
			pbanim.Clips = append(pbanim.Clips, pbclip)
		}
		file.Animations[pbanim.Name] = pbanim
	}
	return file, nil
}

func getClip(file *pb.AtlasFile, clip AnimationClip) (*pb.AnimationClip, error) {
	pbclip := &pb.AnimationClip{
		Name:     clip.Name,
		ClipMode: importClipMode(clip.ClipMode),
		Fps:      float32(clip.FPS),
		Frames:   make([]*pb.AnimFrame, 0),
	}
	if clip.EndedEvent != nil {
		pbclip.EndedEvent = &pb.AnimationEvent{
			Name:  clip.EndedEvent.EventName,
			Value: clip.EndedEvent.EventValue,
		}
	}
	evmap := make(map[int]*AnimEventIO)
	for _, v := range clip.Events {
		evmap[v.Frame] = &v
	}
	for fi, fv := range clip.Frames {
		if file.Frames[fv] == nil {
			return nil, errors.New("frame not found: " + fv)
		}
		f := &pb.AnimFrame{
			FrameName: fv,
		}
		if v := evmap[fi]; v != nil {
			f.Event = &pb.AnimationEvent{
				Name:  v.EventName,
				Value: v.EventValue,
			}
		}
		pbclip.Frames = append(pbclip.Frames, f)
	}
	return pbclip, nil
}

func getNameByPattern(pattern, name string, posindex, nindex, frame int) string {
	v := strings.Replace(pattern, "{name}", name, -1)
	v = strings.Replace(v, "{#}", strconv.Itoa(posindex), -1)
	v = strings.Replace(v, "{#pos}", strconv.Itoa(posindex), -1)
	v = strings.Replace(v, "{#n}", strconv.Itoa(nindex), -1)
	v = strings.Replace(v, "{#f}", strconv.Itoa(frame), -1)
	v = strings.Replace(v, "{#frame}", strconv.Itoa(frame), -1)
	v = strings.Replace(v, "#", strconv.Itoa(posindex), -1)
	return v
}

func importClipMode(aclipmode string) pb.AnimationClipMode {
	switch strings.ToLower(aclipmode) {
	case "once", "forward":
		return pb.AnimationClipMode_ONCE
	// REVERSE IS NOT SUPPORTED!
	//case "reverse":
	//	return pb.?
	case "pingpong", "ping_pong":
		return pb.AnimationClipMode_PING_PONG
	case "loop":
		return pb.AnimationClipMode_LOOP
	case "clamp", "clamp_forever", "clampforever":
		return pb.AnimationClipMode_CLAMP_FOREVER
	}
	return pb.AnimationClipMode_ONCE
}

// . . .-. .   .-. .-. .-.   .-. . . .-. .-. .-.
// |-| |-  |   |-' |-  |(     |   |  |-' |-  `-.
// ' ` `-' `-' '   `-' ' '    '   `  '   `-' `-'

type imImporter struct {
	sprites []imSprite
}

func newImImporter() *imImporter {
	return &imImporter{
		sprites: make([]imSprite, 0, 16),
	}
}

func (i *imImporter) addSprite(name string, node *atlaspacker.RectPackerNode, img image.Image, pivot Vec2) {
	i.sprites = append(i.sprites, imSprite{
		Name:  name,
		Node:  node,
		Image: img,
		Pivot: pivot,
	})
}

func (i *imImporter) Sprites() []imSprite {
	return i.sprites
}

type imSprite struct {
	Name  string
	Node  *atlaspacker.RectPackerNode
	Image image.Image
	Pivot Vec2

	// after the bin packer is calculated

	FinalRect       FrameRect
	FinalImageIndex int
}

// this is not concurrent safe
type atlasIMCache struct {
	m map[string]image.Image
}

// imget
func (c *atlasIMCache) getSubImage(src image.Image, rect FrameRect) image.Image {
	if c.m == nil {
		c.m = make(map[string]image.Image)
	}
	if img, ok := c.m[rect.String()]; ok {
		return img
	}
	subim := image.NewRGBA(image.Rect(0, 0, rect.W, rect.H))
	draw.Draw(subim, subim.Bounds(), src, rect.ToRect().Min, draw.Over)
	c.m[rect.String()] = subim
	return subim
}

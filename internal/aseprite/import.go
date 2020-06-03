package aseprite

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
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
		Images:    make([][]byte, 0),
		Filters:   make([]pb.ImageFilter, 0),
		Frames:    make(map[string]*pb.Frame),
		AnimClips: make(map[string]*pb.AnimationClip),
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
	})
	if err != nil {
		return nil, err
	}
	if err := wrapAnimations(ctx, file, *input.Template); err != nil {
		return nil, err
	}
	return file, nil
}

func importAtlas(ctx context.Context, tpl AtlasImporter, ase AsepriteInput, pkr *atlaspacker.BinTreeRectPacker, imptr *imImporter) error {
	img, err := png.Decode(bytes.NewReader(ase.ImageData))
	if err != nil {
		return fmt.Errorf("error decoding png image: %w", err)
	}
	//idm := make(map[image.Image]*atlaspacker.RectPackerNode)
	//idmr := make(map[*atlaspacker.RectPackerNode]image.Image)
	//
	strat := getStrategy(tpl)
	if err := isValidForStrat(tpl, strat); err != nil {
		return err
	}
	rules := importByRules{
		Template:  tpl,
		Filename:  ase.Filename,
		FrameData: ase.FrameData,
		Img:       img,
		Pkr:       pkr,
		Imptr:     imptr,
	}
	switch strat {
	case Slices:
		return importAtlasBySlices(ctx, rules)
	case FrameTags:
		return importAtlasByFrameTags(ctx, rules)
	case Frames:
		return importAtlasByFrames(ctx, rules)
	}
	panic("invalid strategy should be detected at isValidForStrat")
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
			if frame.OutputName == "" {
				// no way to import a frame without Output name
				return true
			}
			clipim := imbank.getSubImage(r.Img, i.Frame)
			if _, ok := idm[clipim]; !ok {
				node := r.Pkr.AddRect(clipim.Bounds())
				idm[clipim] = node
				r.Imptr.addSprite(frame.OutputName, node, clipim, Vec2{})
			} else {
				r.Imptr.addSprite(frame.OutputName, idm[clipim], clipim, Vec2{})
			}
		}
		return true
	})
	return nil
}

func importAtlasBySlices(ctx context.Context, r importByRules) error {
	idm := make(map[image.Image]*atlaspacker.RectPackerNode)
	//r.Template.
	imbank := &atlasIMCache{}
	//
	// get lone slices!
	for _, slctpl := range r.Template.Slices {
		// get slice from aseprite
		rslc, ri := r.FrameData.GetSliceByName(slctpl.Name)
		if ri == -1 {
			//TODO: dont fail if not strict
			return fmt.Errorf("tplimporter: missing slice '%v'", slctpl.Name)
		}
		for i, v := range rslc.Keys {
			clipim := imbank.getSubImage(r.Img, v.Bounds)
			fname := getNameByPattern(slctpl.OutputPattern, slctpl.Name, i, ri, v.Frame)
			if _, ok := idm[clipim]; !ok {
				node := r.Pkr.AddRect(clipim.Bounds())
				idm[clipim] = node
				r.Imptr.addSprite(fname, node, clipim, v.Pivot)
			} else {
				r.Imptr.addSprite(fname, idm[clipim], clipim, v.Pivot)
			}
		}
	}
	for _, animtpl := range r.Template.AnimationClips {
		rslc, ri := r.FrameData.GetSliceByName(animtpl.Slice)
		if ri == -1 {
			//TODO: dont fail if not strict
			return fmt.Errorf("tplimporter: animation -> missing slice '%v'", animtpl.Slice)
		}
		//
		clip := r.Imptr.addAnimClip(animtpl.OutputName, importClipMode(animtpl.ClipMode))
		clip.Fps = animtpl.FPS
		//
		for _, v := range rslc.Keys {
			clipim := imbank.getSubImage(r.Img, v.Bounds)
			//fname := getNameByPattern(slctpl.OutputPattern, slctpl.Name, i, ri, v.Frame)
			if _, ok := idm[clipim]; !ok {
				node := r.Pkr.AddRect(clipim.Bounds())
				idm[clipim] = node //TODO: support pivots/offsets?
				//r.Imptr.addSprite(fname, node, clipim)
				clip.AddFrame(node, clipim, v.Pivot, animtpl.Events)
			} else {
				clip.AddFrame(idm[clipim], clipim, v.Pivot, animtpl.Events)
			}
			if animtpl.EndedEvent != nil {
				clip.EndEvent = &pb.AnimationEvent{
					Name:  animtpl.EndedEvent.EventName,
					Value: animtpl.EndedEvent.EventValue,
				}
			}
		}
		//
	}
	return nil
}

func importAtlasByFrameTags(ctx context.Context, r importByRules) error {
	return errors.New("not implemented")
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
	draw.Draw(subim, subim.Bounds(), src, rect.ToRect().Min, draw.Src)
	c.m[rect.String()] = subim
	return subim
}

type imImporter struct {
	sprites []imSprite
	clips   []*imAnimClip
}

func newImImporter() *imImporter {
	return &imImporter{
		sprites: make([]imSprite, 0, 16),
		clips:   make([]*imAnimClip, 0, 8),
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

func (i *imImporter) AnimClips() []*imAnimClip {
	return i.clips
}

func (i *imImporter) addAnimClip(name string, clipMode pb.AnimationClipMode) *imAnimClip {
	v := &imAnimClip{
		Name:     name,
		Frames:   make([]*imAnimClipFrame, 0),
		ClipMode: clipMode,
		//FIXME: more stuff
	}
	i.clips = append(i.clips, v)
	return v
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

type imAnimClip struct {
	Name     string
	Fps      int
	Frames   []*imAnimClipFrame
	EndEvent *pb.AnimationEvent
	ClipMode pb.AnimationClipMode
}

func (c *imAnimClip) AddFrame(node *atlaspacker.RectPackerNode, img image.Image, pivot Vec2, events []AnimEventIO) *imAnimClipFrame {
	f := &imAnimClipFrame{
		Raw:   node,
		Pivot: pivot,
		Img:   img,
	}
	nexti := len(c.Frames)
	for _, v := range events {
		if v.Frame == nexti {
			f.Event = &pb.AnimationEvent{
				Name:  v.EventName,
				Value: v.EventValue,
			}
		}
	}
	c.Frames = append(c.Frames, f)
	return f
}

type imAnimClipFrame struct {
	Raw             *atlaspacker.RectPackerNode
	Img             image.Image
	Bounds          FrameRect
	FinalImageIndex int
	Pivot           Vec2
	Event           *pb.AnimationEvent
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
}

func buildAtlas(ctx context.Context, input buildAtlasInput) (*pb.AtlasFile, error) {
	pki := input.PackerI
	pkr := input.Pkr
	im := input.Im
	file := &pb.AtlasFile{
		Images:    make([][]byte, 0),
		Filters:   make([]pb.ImageFilter, 0),
		Frames:    make(map[string]*pb.Frame),
		AnimClips: make(map[string]*pb.AnimationClip),
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
	xclips := im.AnimClips()
	anims := make(map[string]*pb.AnimationClip)
	for i, atlas := range atlases {
		outimg := image.NewRGBA(image.Rect(0, 0, atlas.Width, atlas.Height))
		frames := make(map[string]*pb.Frame)
		for _, spr := range xsprites {
			for _, nn := range atlas.Nodes {
				if nn == spr.Node {
					spr.FinalRect = FrameRect{
						X: nn.X,
						Y: nn.Y,
						W: nn.Width,
						H: nn.Height,
					}
					spr.FinalImageIndex = i
					//TODO: prevent from drawing twice (if nn was already done)
					draw.Draw(outimg, nn.R(), spr.Image, image.ZP, draw.Over)
					frames[spr.Name] = &pb.Frame{
						Image: uint32(i),
						X:     uint32(nn.X),
						Y:     uint32(nn.Y),
						W:     uint32(nn.Width),
						H:     uint32(nn.Height),
					}
				}
			}
		}
		for _, clip := range xclips {
			//TODO: handle animation sprite fragmentation
			outclip := anims[clip.Name]
			if outclip == nil {
				outclip = &pb.AnimationClip{
					ClipMode:   clip.ClipMode,
					Fps:        float32(clip.Fps),
					EndedEvent: clip.EndEvent,
					Name:       clip.Name,
					Frames:     make([]*pb.AnimFrame, 0),
				}
				file.AnimClips[clip.Name] = outclip
			}
			anims[outclip.Name] = outclip
			for _, clipFrame := range clip.Frames {
				for _, nn := range atlas.Nodes {
					if nn == clipFrame.Raw {
						clipFrame.Bounds = FrameRect{
							X: nn.X,
							Y: nn.Y,
							W: nn.Width,
							H: nn.Height,
						}
						clipFrame.FinalImageIndex = i
						//TODO: prevent from drawing twice (if nn was already done)
						draw.Draw(outimg, nn.R(), clipFrame.Img, image.ZP, draw.Over)
					}
				}
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
	}
	// put animation clips in the right order
	for _, clip := range xclips {
		for _, outClip := range clip.Frames {
			file.AnimClips[clip.Name].Frames = append(file.AnimClips[clip.Name].Frames, &pb.AnimFrame{
				X:     uint32(outClip.Bounds.X),
				Y:     uint32(outClip.Bounds.Y),
				W:     uint32(outClip.Bounds.W),
				H:     uint32(outClip.Bounds.H),
				Event: outClip.Event,
				Image: uint32(outClip.FinalImageIndex),
				Ox:    int32(outClip.Pivot.X * -1), //TODO: debug
				Oy:    int32(outClip.Pivot.Y * -1),
			})
		}
	}
	return file, nil
}

func wrapAnimations(ctx context.Context, file *pb.AtlasFile, g AtlasImporterGroup) error {
	file.AnimGroups = make(map[string]*pb.AnimationGroup)
	for _, a := range g.Animations {
		g := &pb.AnimationGroup{
			Name:  a.Name,
			Clips: make(map[string]string),
		}
		for _, item := range a.Clips {
			if file.AnimClips == nil || file.AnimClips[item.GlobalName] == nil {
				return errors.New("animation clip " + item.GlobalName + " not found")
			}
			g.Clips[item.GlobalName] = item.LocalName
		}
		file.AnimGroups[a.Name] = g
	}
	return nil
}

func getStrategy(tpl AtlasImporter) AtlasImportStrategy {
	if tpl.ImportStrategy == Default {
		if len(tpl.Slices) > 0 {
			return Slices
		}
		if len(tpl.FrameTags) > 0 {
			return FrameTags
		}
		return Frames
	}
	return tpl.ImportStrategy
}

func isValidForStrat(tpl AtlasImporter, s AtlasImportStrategy) error {
	if s == Default {
		panic("isValidForStrat != Default")
	}
	switch s {
	case Slices:
		if len(tpl.Slices) < 1 {
			return fmt.Errorf("%s: chosen strategy is Slices, but AtlasImporter -> Slices is empty", tpl.AsepriteSheet)
		}
		return nil
	case FrameTags:
		if len(tpl.FrameTags) < 1 {
			return fmt.Errorf("%s: chosen strategy is FrameTags, but AtlasImporter -> FrameTags is empty", tpl.AsepriteSheet)
		}
	case Frames:
		if len(tpl.Frames) < 1 {
			return fmt.Errorf("%s: chosen strategy is Frames, but AtlasImporter -> Frames is empty", tpl.AsepriteSheet)
		}
	}
	return fmt.Errorf("%s: unknown import strategy %s", tpl.AsepriteSheet, s)
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

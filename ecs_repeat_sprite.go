package troupe

import (
	"github.com/hajimehoshi/ebiten"
)

// RepeatSpritePriority right after Sprite system runs
const RepeatSpritePriority int = -11

// RepeatSprite is the component data of a RepeatSpriteComponent
type RepeatSprite struct {
	RepeatX      int
	RepeatY      int
	DrawDisabled bool
}

// RepeatSpriteComponent returns the registered RepeatSpriteComponent for the world
func RepeatSpriteComponent(w Worlder) *Component {
	c := w.Component("troupe.RepeatSpriteComponent")
	if c == nil {
		var err error
		c, err = w.NewComponent(NewComponentInput{
			Name: "troupe.RepeatSpriteComponent",
			ValidateDataFn: func(data interface{}) bool {
				_, ok := data.(*RepeatSprite)
				return ok
			},
			DestructorFn: func(_ WorldDicter, entity Entity, data interface{}) {
				//sd := data.(*RepeatSprite)
				//sd.Options = nil
			},
		})
		if err != nil {
			panic(err)
		}
	}
	return c
}

// RepeatSpriteSystem upserts the System to the world
func RepeatSpriteSystem(w *World) *System {
	if sys := w.System("x.RepeatSpriteSystem"); sys != nil {
		return sys
	}
	sys := w.NewSystem("x.RepeatSpriteSystem", RepeatSpritePriority, RepeatSpriteSystemExec,
		RepeatSpriteComponent(w), SpriteComponent(w))
	if w.Get(DefaultImageOptions) == nil {
		opt := &ebiten.DrawImageOptions{}
		w.Set(DefaultImageOptions, opt)
	}
	sys.AddTag(WorldTagDraw)
	return sys
}

// RepeatSpriteSystemExec is the function executed every frame
func RepeatSpriteSystemExec(ctx Context, screen *ebiten.Image) {
	view := ctx.System().View()
	repeatComp := RepeatSpriteComponent(ctx.World())
	spriteComp := SpriteComponent(ctx.World())
	var opt *ebiten.DrawImageOptions
	var w, h float64
	var rpx, rpy float64
	var limg *ebiten.Image
	for _, item := range view.Matches() {
		sprite := item.Components[spriteComp].(*Sprite)
		rsprite := item.Components[repeatComp].(*RepeatSprite)
		if rsprite.DrawDisabled {
			continue
		}
		limg = sprite.GetPrecomputedImage()
		//
		opt = anyDrawImageOptions(sprite.Options, ctx.DefaultDrawImageOptions())
		w, h = sprite.GetPrecomputedImageDim()
		//
		rpx, rpy = float64(rsprite.RepeatX), float64(rsprite.RepeatY)
		//
		opt.GeoM.Reset()
		txx := (-.5 + sprite.OriginX) * (w * rpx)
		tyy := (-.5 + sprite.OriginY) * (h * rpy)
		opt.GeoM.Translate(txx, tyy)
		baseGeo := opt.GeoM
		//
		for xx := 0; xx < int(rpx); xx++ {
			xf := float64(xx) * w
			for yy := 0; yy < int(rpy); yy++ {
				yf := float64(yy) * h
				opt.GeoM = baseGeo
				opt.GeoM.Translate(xf, yf)
				opt.GeoM.Scale(sprite.ScaleX, sprite.ScaleY)
				opt.GeoM.Rotate(sprite.Angle)
				opt.GeoM.Translate(sprite.X, sprite.Y)
				_ = screen.DrawImage(limg, opt) //TODO: verify
			}
		}
	}
}

// RepeatSpriteArchetype order:
// SpriteComponent(w)
// RepeatSpriteComponent(w)
func RepeatSpriteArchetype(w *World) *Archetype {
	return NewArchetype(w, SpriteComponent(w), RepeatSpriteComponent(w))
}

func init() {
	DefaultComp(func(e *Engine, w *World) {
		RepeatSpriteComponent(w)
	})
	DefaultSys(func(e *Engine, w *World) {
		RepeatSpriteSystem(w)
	})
}

func anyDrawImageOptions(a ...*ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	for _, v := range a {
		if v != nil {
			return v
		}
	}
	return &ebiten.DrawImageOptions{}
}

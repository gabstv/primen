package tau

import (
	"github.com/gabstv/ecs"
	"github.com/hajimehoshi/ebiten"
)

const (
	SNRepeatSprite string = "tau.RepeatSpriteSystem"
	CNRepeatSprite string = "tau.RepeatSpriteComponent"
)

// RepeatSprite is the component data of a RepeatSpriteComponent
type RepeatSprite struct {
	RepeatX      int
	RepeatY      int
	DrawDisabled bool
}

type RepeatSpriteComponentSystem struct {
	BaseComponentSystem
}

func (cs *RepeatSpriteComponentSystem) SystemName() string {
	return SNRepeatSprite
}

func (cs *RepeatSpriteComponentSystem) SystemPriority() int {
	return -11
}

func (cs *RepeatSpriteComponentSystem) SystemExec() SystemExecFn {
	return RepeatSpriteSystemExec
}

func (cs *RepeatSpriteComponentSystem) SystemTags() []string {
	return []string{"draw"}
}

func (cs *RepeatSpriteComponentSystem) Components(w *ecs.World) []*ecs.Component {
	return []*ecs.Component{
		repeatSpriteComponentDef(w),
		spriteComponentDef(w),
	}
}

func repeatSpriteComponentDef(w *ecs.World) *ecs.Component {
	return UpsertComponent(w, ecs.NewComponentInput{
		Name: CNRepeatSprite,
		ValidateDataFn: func(data interface{}) bool {
			_, ok := data.(*RepeatSprite)
			return ok
		},
		DestructorFn: func(_ ecs.WorldDicter, entity ecs.Entity, data interface{}) {
			//sd := data.(*RepeatSprite)
			//sd.Options = nil
		},
	})
}

// RepeatSpriteSystemExec is the function executed every frame
func RepeatSpriteSystemExec(ctx Context) {
	screen := ctx.Screen()
	view := ctx.System().View()
	repeatComp := ctx.World().Component(CNRepeatSprite)
	spriteComp := ctx.World().Component(CNSprite)
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
func RepeatSpriteArchetype(w *ecs.World) *Archetype {
	return NewArchetype(w, w.Component(CNSprite), w.Component(CNRepeatSprite))
}

// func init() {
// 	DefaultComp(func(e *Engine, w *World) {
// 		RepeatSpriteComponent(w)
// 	})
// 	DefaultSys(func(e *Engine, w *World) {
// 		RepeatSpriteSystem(w)
// 	})
// }

func anyDrawImageOptions(a ...*ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	for _, v := range a {
		if v != nil {
			return v
		}
	}
	return &ebiten.DrawImageOptions{}
}

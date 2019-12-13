package xsystem

import (
	"github.com/gabstv/troupe/pkg/troupe"
	"github.com/hajimehoshi/ebiten"
)

func SkipFrames(n int) troupe.SystemMiddleware {
	return func(next troupe.SystemFn) troupe.SystemFn {
		return func(ctx troupe.Context, screen *ebiten.Image) {
			vi := ctx.System().Get("SkipFrames")
			if vi == nil {
				ctx.System().Set("SkipFrames", 0)
				next(ctx, screen)
				return
			}
			v := vi.(int)
			ctx.System().Set("SkipFrames", v+1)
			if v < n {
				return
			}
			ctx.System().Set("SkipFrames", 0)
			next(ctx, screen)
		}
	}
}

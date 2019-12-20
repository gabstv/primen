package systems

import (
	"github.com/gabstv/troupe"
	"github.com/hajimehoshi/ebiten"
)

func MidSkipFrames(n int) troupe.SystemMiddleware {
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

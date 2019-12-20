package smid

import (
	"github.com/gabstv/troupe"
	"github.com/hajimehoshi/ebiten"
)

// SkipFrames is a System middleware that skips n frames and then executes the
// next function in the system function stack
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

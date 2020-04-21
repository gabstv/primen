package smid

import (
	"github.com/gabstv/tau"
)

// SkipFrames is a System middleware that skips n frames and then executes the
// next function in the system function stack
func SkipFrames(n int) tau.Middleware {
	return func(next tau.SystemExecFn) tau.SystemExecFn {
		return func(ctx tau.Context) {
			vi := ctx.System().Get("SkipFrames")
			if vi == nil {
				ctx.System().Set("SkipFrames", 0)
				next(ctx)
				return
			}
			v := vi.(int)
			ctx.System().Set("SkipFrames", v+1)
			if v < n {
				return
			}
			ctx.System().Set("SkipFrames", 0)
			next(ctx)
		}
	}
}

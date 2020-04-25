package tau

// MidSkipFrames is a System middleware that skips n frames and then executes the
// next function in the system function stack
func MidSkipFrames(n int) Middleware {
	return func(next SystemExecFn) SystemExecFn {
		return func(ctx Context) {
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

package imgui

import "github.com/gabstv/primen/core"

type Uinstance interface {
	Render(ctx core.DrawCtx)
	ID() UID
}

type uiFuncInstance struct {
	fn func(ctx core.DrawCtx)
	id UID
}

var _ Uinstance = (*uiFuncInstance)(nil)

func (ui *uiFuncInstance) Render(ctx core.DrawCtx) {
	ui.fn(ctx)
}

func (ui *uiFuncInstance) ID() UID {
	return ui.id
}

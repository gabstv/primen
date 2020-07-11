package imgui

import (
	"sync"

	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/dom"
	"github.com/hajimehoshi/ebiten"
)

type UID int64
type uiModule int

var (
	lock         sync.Mutex
	mainRenderer *renderer.Manager
	renderTarget *ebiten.Image
	uiinsts      []*UI
	lastid       UID
	lastEngine   core.Engine
)

func (uiModule) BeforeUpdate(ctx core.UpdateCtx) {
	//TODO: solve scaling
	w32 := float32(ctx.Engine().Width())
	h32 := float32(ctx.Engine().Height())
	mainRenderer.Update(float32(ctx.DT()), w32, h32)
}

func (uiModule) AfterUpdate(ctx core.UpdateCtx) {

}

func (uiModule) BeforeDraw(ctx core.DrawCtx) {
	if renderTarget == nil {
		return
	}
	mainRenderer.BeginFrame()
	// TODO: render all doms here
	mainRenderer.EndFrame(renderTarget)
}

func (uiModule) AfterDraw(ctx core.DrawCtx) {
	if renderTarget != nil {
		return
	}
	mainRenderer.BeginFrame()
	for _, ui := range uiinsts {
		ui.Render(ctx)
	}
	// TODO: render all doms here
	mainRenderer.EndFrame(ctx.Renderer().Screen())
}

func Setup(engine core.Engine) {
	lock.Lock()
	defer lock.Unlock()
	if mainRenderer != nil {
		panic("Setup called twice")
	}
	//fa := imgui.CurrentIO().Fonts()
	//fa.AddFontDefault()
	//mainRenderer = renderer.New(&fa)
	mainRenderer = renderer.New(nil)
	mainRenderer.ClipMask = true
	engine.AddModule(uiModule(0), 0)
	lastEngine = engine
}

func AddUI(doc dom.ElementNode) UID {
	lock.Lock()
	defer lock.Unlock()

	lastid++
	id := lastid

	ui := newUI(id, doc)

	uiinsts = append(uiinsts, ui)
	return id
}

// TODO: RemoveUI(id)

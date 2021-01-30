package imgui

import (
	"sync"

	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/dom"
	"github.com/hajimehoshi/ebiten/v2"
)

type UID int64
type uiModule int

var (
	lock         sync.Mutex
	mainRenderer *renderer.Manager
	renderTarget *ebiten.Image
	uiinsts      []Uinstance
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
	// TODO: render all doms here (?)
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
	// call Update once to ensure that the screen is updated before the first Draw() call
	mainRenderer.Update(1.0/60, float32(engine.Width()), float32(engine.Height()))
}

func AddUI(doc []dom.Node) UID {
	lock.Lock()
	defer lock.Unlock()
	if mainRenderer == nil {
		panic("imgui.Setup(engine) needs to be called before AddUI")
	}

	lastid++
	id := lastid

	ui := newUI(id, doc)

	uiinsts = append(uiinsts, ui)
	return id
}

func AddRawUI(renderfn func(ctx core.DrawCtx)) UID {
	lock.Lock()
	defer lock.Unlock()
	if mainRenderer == nil {
		panic("imgui.Setup(engine) needs to be called before AddUI")
	}

	lastid++
	id := lastid

	ui := &uiFuncInstance{
		fn: renderfn,
		id: id,
	}
	uiinsts = append(uiinsts, ui)
	return id
}

func RemoveUI(id UID) {
	lock.Lock()
	defer lock.Unlock()
	index := -1
	for i := range uiinsts {
		if uiinsts[i].ID() == id {
			index = i
		}
	}
	if index == -1 {
		return
	}
	uiinsts = uiinsts[:index+copy(uiinsts[index:], uiinsts[index+1:])]
}

func SetFilter(filter ebiten.Filter) {
	mainRenderer.Filter = filter
	mainRenderer.Cache.ResetFontAtlasCache(filter)
}

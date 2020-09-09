package core

// Module defines a Primen Engine module
type Module interface {
	BeforeUpdate(ctx UpdateCtx)
	AfterUpdate(ctx UpdateCtx)
	BeforeDraw(ctx DrawCtx)
	AfterDraw(ctx DrawCtx)
}

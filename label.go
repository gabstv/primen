package primen

// import (
// 	"image"
// 	"image/color"

// 	"github.com/gabstv/primen/core"
// 	"github.com/hajimehoshi/ebiten"
// 	"golang.org/x/image/font"
// )

// type Label struct {
// 	*WorldItem
// 	*TransformItem
// 	*DrawLayerItem
// 	label func() *core.Label
// }

// func NewLabel(parent WorldTransform, fontFace font.Face, layer Layer) *Label {
// 	lbl := &Label{}
// 	w := parent.World()
// 	e := w.NewEntity()
// 	lbl.WorldItem = newWorldItem(e, w)
// 	core.SetLabelComponentData(w, e, core.Label{
// 		ScaleX: 1,
// 		ScaleY: 1,
// 		Color:  color.White,
// 		Face:   fontFace,
// 		Filter: ebiten.FilterDefault,
// 	})
// 	lbl.label = func() *core.Label { return core.GetLabelComponentData(w, e) }
// 	lbl.label().ResetTextOffset()
// 	lbl.DrawLayerItem = newDrawLayerItem(e, w)
// 	lbl.TransformItem = newTransformItem(e, parent)
// 	return lbl
// }

// func (l *Label) SetText(t string) {
// 	if l.label().Text == t {
// 		return
// 	}
// 	l.label().Text = t
// 	l.label().SetDirty()
// }

// func (l *Label) Text() string {
// 	return l.label().Text
// }

// func (l *Label) SetArea(w, h int) {
// 	l.label().Area = image.Point{
// 		X: w,
// 		Y: h,
// 	}
// 	l.label().SetDirty()
// }

// func (l *Label) SetFilter(filter ebiten.Filter) {
// 	l.label().Filter = filter
// 	l.label().SetDirty()
// }

// // SetOrigin sets the label origin reference
// //
// // 0, 0 is top left
// //
// // 1, 1 is bottom right
// func (l *Label) SetOrigin(ox, oy float64) {
// 	l.label().OriginX = ox
// 	l.label().OriginY = oy
// }

// func (l *Label) ComputedSize() (w, h int) {
// 	p := l.label().ComputedSize()
// 	return p.X, p.Y
// }

// func (l *Label) SetColor(c color.Color) {
// 	l.label().Color = c
// 	l.label().SetDirty()
// }

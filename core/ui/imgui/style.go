package imgui

import (
	"strings"

	"github.com/dop251/goja"
	"github.com/gabstv/primen/core"
	"github.com/gabstv/primen/internal/z"
	"github.com/inkyblackness/imgui-go/v2"
)

type styleUpdateFn func(nodeid, v string, memory *uiMemory) (styles, colors int)

func styleFloatGetter(name string, style imgui.StyleVarID) styleUpdateFn {
	return func(nodeid, v string, memory *uiMemory) (styles, colors int) {
		if vf, ok := memory.ImGuiFloat(nodeid + "->" + name + "f"); ok {
			imgui.PushStyleVarFloat(style, vf)
			return 1, 0
		}
		vf, ok := z.Float32V(v)
		if !ok {
			return 0, 0
		}
		memory.Add(nodeid+"->"+name+"f", vf)
		imgui.PushStyleVarFloat(style, vf)
		return 1, 0
	}
}

func styleVec2Getter(name string, style imgui.StyleVarID) styleUpdateFn {
	return func(nodeid, v string, memory *uiMemory) (styles, colors int) {
		if vf, ok := memory.ImGuiVec2(nodeid + "->" + name + "vec2"); ok {
			imgui.PushStyleVarVec2(style, vf)
			return 1, 0
		}
		vf, ok := parseVec2(v)
		if !ok {
			return 0, 0
		}
		memory.Add(nodeid+"->"+name+"vec2", vf)
		imgui.PushStyleVarVec2(style, vf)
		return 1, 0
	}
}

func colorVec4Getter(name string, style imgui.StyleColorID) styleUpdateFn {
	return func(nodeid, v string, memory *uiMemory) (styles, colors int) {
		if vf, ok := memory.ImGuiVec4(nodeid + "->" + name + "vec4c"); ok {
			imgui.PushStyleColor(style, vf)
			return 0, 1
		}
		vf, ok := parseVec4(v)
		if !ok {
			return 0, 0
		}
		memory.Add(nodeid+"->"+name+"vec4c", vf)
		imgui.PushStyleColor(style, vf)
		return 0, 1
	}
}

var stparser = map[string]styleUpdateFn{
	"st-alpha":                 styleFloatGetter("alpha", imgui.StyleVarAlpha),
	"st-window-padding":        styleVec2Getter("window-padding", imgui.StyleVarWindowPadding),
	"st-window-rounding":       styleFloatGetter("window-rounding", imgui.StyleVarWindowRounding),
	"st-window-border-size":    styleFloatGetter("window-border-size", imgui.StyleVarWindowBorderSize),
	"st-window-min-size":       styleVec2Getter("window-min-size", imgui.StyleVarWindowMinSize),
	"st-window-title-align":    styleVec2Getter("window-title-align", imgui.StyleVarWindowTitleAlign),
	"st-child-rounding":        styleFloatGetter("child-rounding", imgui.StyleVarChildRounding),
	"st-child-border-size":     styleFloatGetter("child-border-size", imgui.StyleVarChildBorderSize),
	"st-popup-rounding":        styleFloatGetter("popup-rounding", imgui.StyleVarPopupRounding),
	"st-popup-border-size":     styleFloatGetter("popup-border-size", imgui.StyleVarPopupBorderSize),
	"st-frame-padding":         styleVec2Getter("frame-padding", imgui.StyleVarFramePadding),
	"st-frame-rounding":        styleFloatGetter("frame-rounding", imgui.StyleVarFrameRounding),
	"st-frame-border-size":     styleFloatGetter("frame-border-size", imgui.StyleVarFrameBorderSize),
	"st-item-spacing":          styleVec2Getter("item-spacing", imgui.StyleVarItemSpacing),
	"st-item-inner-spacing":    styleVec2Getter("item-inner-spacing", imgui.StyleVarItemInnerSpacing),
	"st-indent-spacing":        styleFloatGetter("indent-spacing", imgui.StyleVarIndentSpacing),
	"st-scrollbar-size":        styleFloatGetter("scrollbar-size", imgui.StyleVarScrollbarSize),
	"st-scrollbar-rounding":    styleFloatGetter("scrollbar-rounding", imgui.StyleVarScrollbarRounding),
	"st-grab-min-size":         styleFloatGetter("grab-min-size", imgui.StyleVarGrabMinSize),
	"st-grab-rounding":         styleFloatGetter("grab-rounding", imgui.StyleVarGrabRounding),
	"st-tab-rounding":          styleFloatGetter("tab-rounding", imgui.StyleVarTabRounding),
	"st-button-text-align":     styleVec2Getter("button-text-align", imgui.StyleVarButtonTextAlign),
	"st-selectable-text-align": styleVec2Getter("selectable-text-align", imgui.StyleVarSelectableTextAlign),
	// COLORS
	"st-color-text":                    colorVec4Getter("color-text", imgui.StyleColorText),
	"st-color-text-disabled":           colorVec4Getter("color-text-disabled", imgui.StyleColorTextDisabled),
	"st-color-window-bg":               colorVec4Getter("color-window-bg", imgui.StyleColorWindowBg),
	"st-color-child-bg":                colorVec4Getter("color-child-bg", imgui.StyleColorChildBg),
	"st-color-popup-bg":                colorVec4Getter("color-popup-bg", imgui.StyleColorPopupBg),
	"st-color-border":                  colorVec4Getter("color-border", imgui.StyleColorBorder),
	"st-color-border-shadow":           colorVec4Getter("color-border-shadow", imgui.StyleColorBorderShadow),
	"st-color-frame-bg":                colorVec4Getter("color-frame-bg", imgui.StyleColorFrameBg),
	"st-color-frame-bg-hovered":        colorVec4Getter("color-frame-bg-hovered", imgui.StyleColorFrameBgHovered),
	"st-color-frame-bg-active":         colorVec4Getter("color-frame-bg-active", imgui.StyleColorFrameBgActive),
	"st-color-title-bg":                colorVec4Getter("color-title-bg", imgui.StyleColorTitleBg),
	"st-color-title-bg-active":         colorVec4Getter("color-title-bg-active", imgui.StyleColorTitleBgActive),
	"st-color-title-bg-collapsed":      colorVec4Getter("color-title-bg-collapsed", imgui.StyleColorTitleBgCollapsed),
	"st-color-menu-bar-bg":             colorVec4Getter("color-menu-bar-bg", imgui.StyleColorMenuBarBg),
	"st-color-scrollbar-bg":            colorVec4Getter("color-scrollbar-bg", imgui.StyleColorScrollbarBg),
	"st-color-scrollbar-grab":          colorVec4Getter("color-scrollbar-grab", imgui.StyleColorScrollbarGrab),
	"st-color-scrollbar-grab-hovered":  colorVec4Getter("color-scrollbar-grab-hovered", imgui.StyleColorScrollbarGrabHovered),
	"st-color-scrollbar-grab-active":   colorVec4Getter("color-scrollbar-grab-active", imgui.StyleColorScrollbarGrabActive),
	"st-color-check-mark":              colorVec4Getter("color-check-mark", imgui.StyleColorCheckMark),
	"st-color-slider-grab":             colorVec4Getter("color-slider-grab", imgui.StyleColorSliderGrab),
	"st-color-slider-grab-active":      colorVec4Getter("color-slider-grab-active", imgui.StyleColorSliderGrabActive),
	"st-color-button":                  colorVec4Getter("color-button", imgui.StyleColorButton),
	"st-color-button-hovered":          colorVec4Getter("color-button-hovered", imgui.StyleColorButtonHovered),
	"st-color-button-active":           colorVec4Getter("color-button-active", imgui.StyleColorButtonActive),
	"st-color-header":                  colorVec4Getter("color-header", imgui.StyleColorHeader),
	"st-color-header-hovered":          colorVec4Getter("color-header-hovered", imgui.StyleColorHeaderHovered),
	"st-color-header-active":           colorVec4Getter("color-header-active", imgui.StyleColorHeaderActive),
	"st-color-separator":               colorVec4Getter("color-separator", imgui.StyleColorSeparator),
	"st-color-separator-hovered":       colorVec4Getter("color-separator-hovered", imgui.StyleColorSeparatorHovered),
	"st-color-separator-active":        colorVec4Getter("color-separator-active", imgui.StyleColorSeparatorActive),
	"st-color-resize-grip":             colorVec4Getter("color-resize-grip", imgui.StyleColorResizeGrip),
	"st-color-resize-grip-hovered":     colorVec4Getter("color-resize-grip-hovered", imgui.StyleColorResizeGripHovered),
	"st-color-resize-grip-active":      colorVec4Getter("color-resize-grip-active", imgui.StyleColorResizeGripActive),
	"st-color-tab":                     colorVec4Getter("color-tab", imgui.StyleColorTab),
	"st-color-tab-hovered":             colorVec4Getter("color-tab-hovered", imgui.StyleColorTabHovered),
	"st-color-tab-active":              colorVec4Getter("color-tab-active", imgui.StyleColorTabActive),
	"st-color-tab-unfocused":           colorVec4Getter("color-tab-unfocused", imgui.StyleColorTabUnfocused),
	"st-color-tab-unfocused-active":    colorVec4Getter("color-tab-unfocused-active", imgui.StyleColorTabUnfocusedActive),
	"st-color-plot-lines":              colorVec4Getter("color-plot-lines", imgui.StyleColorPlotLines),
	"st-color-plot-lines-hovered":      colorVec4Getter("color-plot-lines-hovered", imgui.StyleColorPlotLinesHovered),
	"st-color-plot-histogram":          colorVec4Getter("color-plot-histogram", imgui.StyleColorPlotHistogram),
	"st-color-plot-histogram-hovered":  colorVec4Getter("color-plot-histogram-hovered", imgui.StyleColorPlotHistogramHovered),
	"st-color-text-selected-bg":        colorVec4Getter("color-text-selected-bg", imgui.StyleColorTextSelectedBg),
	"st-color-drag-drop-target":        colorVec4Getter("color-drag-drop-target", imgui.StyleColorDragDropTarget),
	"st-color-nav-highlight":           colorVec4Getter("color-nav-highlight", imgui.StyleColorNavHighlight),
	"st-color-nav-windowing-highlight": colorVec4Getter("color-nav-windowing-highlight", imgui.StyleColorNavWindowingHighlight),
	"st-color-nav-windowing-darkening": colorVec4Getter("color-nav-windowing-darkening", imgui.StyleColorNavWindowingDarkening),
	"st-color-modal-window-darkening":  colorVec4Getter("color-modal-window-darkening", imgui.StyleColorModalWindowDarkening),
}

func pushStyles(attributes map[string]string, memory *uiMemory) (styles, colors int) {
	if attributes["id"] == "" {
		println("warning: [pushStyles] element didn't have an ID")
		attributes["id"] = z.Rs()
	}
	for name, rawval := range attributes {
		if len(name) > 3 {
			if name[:3] == "st-" {
				if f, ok := stparser[name]; ok {
					s1, c1 := f(attributes["id"], rawval, memory)
					styles += s1
					colors += c1
				}
			}
		}
	}
	return
}

func popStyles(nstyles, ncolors int) {
	imgui.PopStyleVarV(nstyles)
	imgui.PopStyleColorV(ncolors)
}

// parse position attributes
func setupWindowPos(ctx core.DrawCtx, attributes map[string]string, memory *uiMemory, jsvm *goja.Runtime) {
	switch attributes["position"] {
	case "fixed":
		yset := false
		y := float32(0)
		xset := false
		x := float32(0)
		bset := false
		b := float32(0)
		rset := false
		r := float32(0)
		if vt := attributes["top"]; vt != "" {
			if strings.HasPrefix(vt, "js:") {
				if v, err := jsvm.RunString(vt[3:]); err == nil {
					y = float32(v.ToFloat())
					yset = true
				}
			} else {
				y = z.Float32(vt, 0)
				yset = true
			}
		}
		if vt := attributes["left"]; vt != "" {
			if strings.HasPrefix(vt, "js:") {
				if v, err := jsvm.RunString(vt[3:]); err == nil {
					x = float32(v.ToFloat())
					xset = true
				}
			} else {
				x = z.Float32(vt, 0)
				xset = true
			}
		}
		if xset && yset {
			imgui.SetNextWindowPos(imgui.Vec2{
				X: x,
				Y: y,
			})
			if vt := attributes["bottom"]; vt != "" {
				if strings.HasPrefix(vt, "js:") {
					if v, err := jsvm.RunString(vt[3:]); err == nil {
						b = float32(v.ToFloat())
						bset = true
					}
				} else {
					b = z.Float32(vt, 0)
					bset = true
				}
			}
			if vt := attributes["right"]; vt != "" {
				if strings.HasPrefix(vt, "js:") {
					if v, err := jsvm.RunString(vt[3:]); err == nil {
						r = float32(v.ToFloat())
						rset = true
					}
				} else {
					r = z.Float32(vt, 0)
					rset = true
				}
			}
			if rset && bset {
				w, h := ctx.Renderer().Screen().Size()
				sz := imgui.Vec2{
					X: float32(w) - r - x,
					Y: float32(h) - b - y,
				}
				imgui.SetNextWindowSize(sz)
			}
		}
	}
}

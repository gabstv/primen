package imgui

import (
	"strconv"
	"strings"

	"github.com/inkyblackness/imgui-go/v2"
)

var windowflagmap = map[string]int{
	"none":                        imgui.WindowFlagsNone,
	"no-title-bar":                imgui.WindowFlagsNoTitleBar,
	"no-resize":                   imgui.WindowFlagsNoResize,
	"no-move":                     imgui.WindowFlagsNoMove,
	"no-scrollbar":                imgui.WindowFlagsNoScrollbar,
	"no-scroll-with-mouse":        imgui.WindowFlagsNoScrollWithMouse,
	"no-collapse":                 imgui.WindowFlagsNoCollapse,
	"always-auto-resize":          imgui.WindowFlagsAlwaysAutoResize,
	"no-background":               imgui.WindowFlagsNoBackground,
	"no-saved-settings":           imgui.WindowFlagsNoSavedSettings,
	"no-mouse-inputs":             imgui.WindowFlagsNoMouseInputs,
	"menu-bar":                    imgui.WindowFlagsMenuBar,
	"horizontal-scrollbar":        imgui.WindowFlagsHorizontalScrollbar,
	"no-focus-on-appearing":       imgui.WindowFlagsNoFocusOnAppearing,
	"no-bring-to-front-on-focus":  imgui.WindowFlagsNoBringToFrontOnFocus,
	"always-vertical-scrollbar":   imgui.WindowFlagsAlwaysVerticalScrollbar,
	"always-horizontal-scrollbar": imgui.WindowFlagsAlwaysHorizontalScrollbar,
	"always-use-window-padding":   imgui.WindowFlagsAlwaysUseWindowPadding,
	"no-nav-inputs":               imgui.WindowFlagsNoNavInputs,
	"no-nav-focus":                imgui.WindowFlagsNoNavFocus,
	"unsaved-document":            imgui.WindowFlagsUnsavedDocument,
	"no-nav":                      imgui.WindowFlagsNoNav,
	"no-decoration":               imgui.WindowFlagsNoDecoration,
	"no-inputs":                   imgui.WindowFlagsNoInputs,
}

// checks if first char is a digit (disregards +- signals)
func isRuneNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

func parseWindowFlags(attributes map[string]string) int {
	raw := attributes["flags"]
	if raw == "" {
		return 0
	}
	if isRuneNumeric(rune(raw[0])) {
		v, _ := strconv.Atoi(raw)
		return v
	}
	f := 0
	raw = strings.ReplaceAll(raw, " ", "")
	for _, fi := range strings.Split(raw, "|") {
		f = f | windowflagmap[fi]
	}
	return f
}

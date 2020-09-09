package style

import (
	"strconv"
	"strings"

	"github.com/dop251/goja"
	"github.com/inkyblackness/imgui-go/v2"
)

type Context struct {
	Size imgui.Vec2
}

func (c Context) SizeZero() bool {
	return c.Size.X == 0 && c.Size.Y == 0
}

// func setNodeLayout(node dom.ElementNode, data *store.DB, jsvm *goja.Runtime) (layout Context) {
// 	initjs := false
// 	attrs := node.Attributes()
// 	//if attrs.HasAttr("w", "width")
// 	if v, rel, ok := getattrparsejsNumber(jsvm, &initjs, attrs.FirstAttr("w", "width")); ok {
// 		if rel {
// 			cravail := imgui.ContentRegionAvail()
// 			layout.Size.X = cravail.X * (float32(v) / 100.0)
// 		} else {
// 			layout.Size.X = float32(v)
// 		}
// 	}
// 	if v, rel, ok := getattrparsejsNumber(jsvm, &initjs, attrs.FirstAttr("h", "height")); ok {
// 		if rel {
// 			cravail := imgui.ContentRegionAvail()
// 			layout.Size.X = cravail.Y * (float32(v) / 100.0)
// 		} else {
// 			layout.Size.X = float32(v)
// 		}
// 	}
// 	return
// }

func vhasJS(v string) (bool, string) {
	if strings.HasPrefix(v, "js:") {
		return true, v[3:]
	}
	return false, ""
}

func updateLocalJS(inited *bool, jsvm *goja.Runtime) {
	if *inited {
		return
	}
	jsvm.Set("__column_width", imgui.ColumnWidth())
	*inited = true
}

func getattrparsejsInt(jsvm *goja.Runtime, inited *bool, attrv string) (val int, relative, ok bool) {
	if attrv == "" {
		return 0, false, false
	}
	if hjs, strjs := vhasJS(attrv); hjs {
		updateLocalJS(inited, jsvm)
		vv, err := jsvm.RunString(strjs)
		if err != nil {
			return 0, false, false
		}
		return int(vv.ToInteger()), false, true
	}
	if strings.HasSuffix(attrv, "%") {
		relative = true
		attrv = attrv[:len(attrv)-1]
	}
	vv, err := strconv.Atoi(attrv)
	if err != nil {
		return 0, relative, false
	}
	return vv, relative, true
}

func getattrparsejsNumber(jsvm *goja.Runtime, inited *bool, attrv string) (val float64, relative, ok bool) {
	if attrv == "" {
		return 0, false, false
	}
	if hjs, strjs := vhasJS(attrv); hjs {
		updateLocalJS(inited, jsvm)
		vv, err := jsvm.RunString(strjs)
		if err != nil {
			return 0, false, false
		}
		return vv.ToFloat(), false, true
	}
	if strings.HasSuffix(attrv, "%") {
		relative = true
		attrv = attrv[:len(attrv)-1]
	}
	vv, err := strconv.ParseFloat(attrv, 64)
	if err != nil {
		return 0, relative, false
	}
	return vv, relative, true
}

func getattrparsejsString(jsvm *goja.Runtime, inited *bool, attrv string) (string, bool) {
	if hjs, strjs := vhasJS(attrv); hjs {
		updateLocalJS(inited, jsvm)
		vv, err := jsvm.RunString(strjs)
		if err != nil {
			return "", false
		}
		return vv.String(), true
	}
	return attrv, true
}

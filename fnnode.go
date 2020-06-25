package primen

import (
	"github.com/gabstv/primen/core"
)

type FnNode struct {
	*Node
	wf core.WatchFunction
}

func NewRootFnNode(w World) *FnNode {
	n := &FnNode{
		Node: NewRootNode(w),
	}
	core.SetFunctionComponentData(w, n.e, core.Function{})
	n.wf = core.WatchFunctionComponentData(w, n.e)
	return n
}

// Function retrieves the function component data
func (n *FnNode) Function() *core.Function {
	return n.wf.Data()
}

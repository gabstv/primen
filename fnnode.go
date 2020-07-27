package primen

import (
	"github.com/gabstv/primen/components"
)

type FnNode struct {
	*Node
	wf components.WatchFunction
}

func NewRootFnNode(w World) *FnNode {
	n := &FnNode{
		Node: NewRootNode(w),
	}
	components.SetFunctionComponentData(w, n.e, components.Function{})
	n.wf = components.WatchFunctionComponentData(w, n.e)
	return n
}

// Function retrieves the function component data
func (n *FnNode) Function() *components.Function {
	return n.wf.Data()
}

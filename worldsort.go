package primen

import (
	"github.com/gabstv/primen/core"
)

type worldContainer struct {
	world    *core.GameWorld
	priority int
}

// sortedWorldContainer implements sort.Interface for []worldContainer based on
// the priority field.
type sortedWorldContainer []worldContainer

func (a sortedWorldContainer) Len() int           { return len(a) }
func (a sortedWorldContainer) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortedWorldContainer) Less(i, j int) bool { return a[i].priority > a[j].priority }

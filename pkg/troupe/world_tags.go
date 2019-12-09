package troupe

import (
	"github.com/gabstv/troupe/pkg/troupe/ecs"
)

// WorldTag is a tag used to filter systems of a world
type WorldTag = string

// TagDelta is the key set of the time taken between frames (in seconds)
const TagDelta = "delta"

const (
	// WorldTagDraw -> systems that draw things
	WorldTagDraw WorldTag = "draw"
	// WorldTagUpdate -> systems that update things (!= draw)
	WorldTagUpdate WorldTag = "update"
)

type worldContainer struct {
	world    *ecs.World
	priority int
}

// sortedWorldContainer implements sort.Interface for []worldContainer based on
// the priority field.
type sortedWorldContainer []worldContainer

func (a sortedWorldContainer) Len() int           { return len(a) }
func (a sortedWorldContainer) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortedWorldContainer) Less(i, j int) bool { return a[i].priority > a[j].priority }

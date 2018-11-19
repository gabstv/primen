package groove

// WorldTag is a tag used to filter systems of a world
type WorldTag = string

const (
	// WorldTagDraw -> systems that draw things
	WorldTagDraw WorldTag = "draw"
	// WorldTagUpdate -> systems that update things (!= draw)
	WorldTagUpdate WorldTag = "update"
)

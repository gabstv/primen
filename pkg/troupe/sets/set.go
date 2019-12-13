package sets

var itemExists = struct{}{}

type SetBase interface {
	Empty() bool
	Size() int
	Clear()
}

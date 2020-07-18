package geom

// ZR is a zero vector
var ZR Rect = Rect{}

type Rect struct {
	Min Vec
	Max Vec
}

func (r Rect) Equals(other Rect) bool {
	return r.Min.Equals(other.Min) && r.Max.Equals(other.Max)
}

func (r Rect) IsZero() bool {
	return r.Equals(ZR)
}

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

func (r Rect) ContainsVec(v Vec) bool {
	return r.Min.X <= v.X && v.X <= r.Max.X && r.Min.Y <= v.Y && v.Y <= r.Max.Y
}

// Width returns the width
func (r Rect) Width() float64 {
	return r.Max.X - r.Min.X
}

// Height returns the height
func (r Rect) Height() float64 {
	return r.Max.Y - r.Min.Y
}

// Size returns the vector size (width, height)
func (r Rect) Size() Vec {
	return Vec{r.Width(), r.Height()}
}

func (r Rect) At(v Vec) Rect {
	dx := r.Width()
	dy := r.Height()
	return Rect{
		Min: Vec{v.X, v.Y},
		Max: Vec{v.X + dx, v.Y + dy},
	}
}

func (r Rect) AddVec(v Vec) Rect {
	return Rect{
		Min: r.Min.Add(v),
		Max: r.Max.Add(v),
	}
}

func (r Rect) SubVec(v Vec) Rect {
	return Rect{
		Min: r.Min.Sub(v),
		Max: r.Max.Sub(v),
	}
}

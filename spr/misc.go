package spr

type SourceData struct {
	File   string  `json:"file"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

package aseprite

type AtlasImporter struct {
	Sprites      []FrameIO `json:"sprites"`
	SpriteSheets []string  `json:"sprite_sheets"`
	Output       string    `json:"output,omitempty"`
}

type FrameIO struct {
	Filename   string `json:"filename"`
	OutputName string `json:"output_name"`
	SheetIndex int    `json:"sheet_index"`
	SheetName  string `json:"sheet_name,omitempty"`
}

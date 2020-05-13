package aseprite

type AtlasImporter struct {
	Sprites      []FrameIO `json:"sprites"`
	SpriteSheets []string  `json:"sprite_sheets"`
	Output       string    `json:"output,omitempty"`
	ImageFilter  string    `json:"image_filter,omitempty"`
}

func (i AtlasImporter) SpriteWithFilename(filename string) (frame FrameIO, exists bool) {
	for _, v := range i.Sprites {
		if v.Filename == filename {
			frame = v
			exists = true
			return
		}
	}
	return
}

type FrameIO struct {
	Filename   string `json:"filename"`
	OutputName string `json:"output_name"`
	SheetIndex int    `json:"sheet_index"`
	SheetName  string `json:"sheet_name,omitempty"`
}

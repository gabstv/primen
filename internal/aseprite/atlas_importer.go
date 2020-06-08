package aseprite

// AtlasImporterGroup is the root of an import template
type AtlasImporterGroup struct {
	Templates   []AtlasImporter `json:"templates"`
	Output      string          `json:"output,omitempty"`
	ImageFilter string          `json:"image_filter,omitempty"`
	MaxWidth    int             `json:"max_width"`
	MaxHeight   int             `json:"max_height"`
	Padding     int             `json:"padding"`
	Animations  []Animation     `json:"animations"`
	Clips       []AnimationClip `json:"clips,omitempty"`
}

// AtlasImporter is a template used to import an Aseprite JSON (of a sprite sheet)
// to Primen.
type AtlasImporter struct {
	Frames                []FrameIO `json:"frames"`
	AsepriteSheet         string    `json:"asesprite_sheet"`
	ExportUndefinedFrames bool      `json:"export_undefined_frames"`
}

// FrameWithFilename returns the frame with the specified filename.
func (i AtlasImporter) FrameWithFilename(filename string) (frame FrameIO, exists bool) {
	for _, v := range i.Frames {
		if v.Filename == filename {
			frame = v
			exists = true
			return
		}
	}
	return
}

// FrameIO is the template to import a sprite by the frame itself.
// This is used when the import strategy is set to Frames, or when it is set to Default and
// no frame tags or slices are available.
type FrameIO struct {
	Filename string `json:"filename"`
	Pivot    Vec2   `json:"pivot,omitempty"` //TODO: use
}

type Animation struct {
	Name  string          `json:"name"`
	Clips []AnimationClip `json:"clips"`
}

type AnimationClip struct {
	Name       string        `json:"name"`
	ClipMode   string        `json:"clip_mode,omitempty"`
	Frames     []string      `json:"frames,omitempty"`
	Events     []AnimEventIO `json:"events"`
	EndedEvent *AnimEventIO  `json:"ended_event,omitempty"`
	FPS        int           `json:"fps"`
}

type AnimEventIO struct {
	Frame      int    `json:"frame,omitempty"`
	EventName  string `json:"event_name"`
	EventValue string `json:"event_value"`
}

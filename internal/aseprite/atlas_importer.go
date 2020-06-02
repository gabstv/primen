package aseprite

// AtlasImportStrategy determines the strategy in which the sprite sheets should
// be imported to Primen
type AtlasImportStrategy string

const (
	// Default uses a guess approach to import sprites (and animations).
	// It first looks for Slices on the atlas importer. If slices are not present
	// in the importer, it will try FrameTags. If FrameTags are empty, it looks for
	// Sprites.
	Default AtlasImportStrategy = "default"
	// Frames uses the slice of frames to import sprites.
	Frames AtlasImportStrategy = "frames"
	// Slices uses the slices and slice keys to import sprites.
	Slices AtlasImportStrategy = "slices"
	// FrameTags uses the frame tags (defined on the Aseprite timeline) to import sprites.
	FrameTags AtlasImportStrategy = "frame_tags"
)

type AtlasImporterGroup struct {
	Templates   []AtlasImporter `json:"templates"`
	Output      string          `json:"output,omitempty"`
	ImageFilter string          `json:"image_filter,omitempty"`
	MaxWidth    int             `json:"max_width"`
	MaxHeight   int             `json:"max_height"`
	Padding     int             `json:"padding"`
}

// AtlasImporter is a template used to import an Aseprite JSON (of a sprite sheet)
// to Primen.
type AtlasImporter struct {
	ImportStrategy AtlasImportStrategy `json:"import_strategy"`
	Animations     []AnimationIO       `json:"animations"`
	Slices         []SliceIO           `json:"slices"`
	FrameTags      []FrameTagIO        `json:"frame_tags"`
	Frames         []FrameIO           `json:"frames"`
	AsepriteSheet  string              `json:"asesprite_sheet"`
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
	Filename   string `json:"filename"`
	OutputName string `json:"output_name"`
}

// SliceIO is the template to import sprites by a slice.
type SliceIO struct {
	Name          string `json:"name"`
	OutputPattern string `json:"output_pattern"`
}

type FrameTagIO struct {
	Name          string `json:"name"`
	OutputPattern string `json:"output_pattern"`
}

type AnimationIO struct {
	Slice      string        `json:"slice,omitempty"`
	FrameTag   string        `json:"frame_tag,omitempty"`
	OutputName string        `json:"output_name"`
	ClipMode   string        `json:"clip_mode,omitempty"`
	Events     []AnimEventIO `json:"events"`
	EndedEvent *AnimEventIO  `json:"ended_event,omitempty"`
	FPS        int           `json:"fps"`
}

type AnimEventIO struct {
	Frame      int    `json:"frame,omitempty"`
	EventName  string `json:"event_name"`
	EventValue string `json:"event_value"`
}

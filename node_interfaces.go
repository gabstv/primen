package primen

import (
	"github.com/gabstv/primen/components"
)

type NodeWithTransform interface {
	Transform() *components.Transform
}

type NodeWithAudioPlayer interface {
	AudioPlayer() *components.AudioPlayer
}

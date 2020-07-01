package primen

import "github.com/gabstv/primen/core"

type NodeWithTransform interface {
	Transform() *core.Transform
}

type NodeWithAudioPlayer interface {
	AudioPlayer() *core.AudioPlayer
}

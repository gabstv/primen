package primen

import (
	"github.com/gabstv/primen/components"
)

type AudioPlayerNode struct {
	*mObjectContainer
	wap components.WatchAudioPlayer
}

func NewRootAudioPlayerNode(w World, input components.NewAudioPlayerInput) *AudioPlayerNode {
	tr := &AudioPlayerNode{
		mObjectContainer: &mObjectContainer{
			mObject: &mObject{
				e: w.NewEntity(),
				w: w,
			},
		},
	}
	components.SetAudioPlayerComponentData(w, tr.e, components.NewAudioPlayer(input))
	tr.wap = components.WatchAudioPlayerComponentData(w, tr.Entity())
	return tr
}

func NewChildAudioPlayerNode(parent ObjectContainer, input components.NewAudioPlayerInput) *AudioPlayerNode {
	if parent == nil {
		panic("parent can't be nil")
	}
	tr := &AudioPlayerNode{
		mObjectContainer: &mObjectContainer{
			mObject: &mObject{
				e: parent.World().NewEntity(),
				w: parent.World(),
			},
		},
	}
	components.SetAudioPlayerComponentData(tr.w, tr.e, components.NewAudioPlayer(input))
	tr.wap = components.WatchAudioPlayerComponentData(tr.w, tr.Entity())
	tr.SetParent(parent)
	return tr
}

func (t *AudioPlayerNode) AudioPlayer() *components.AudioPlayer {
	return t.wap.Data()
}

func (t *AudioPlayerNode) Destroy() {
	t.wap.Data().Pause()
	t.wap = nil
	t.mObjectContainer.Destroy()
}

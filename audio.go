package primen

import (
	"github.com/gabstv/primen/core"
)

type AudioPlayerNode struct {
	*mObjectContainer
	wap core.WatchAudioPlayer
}

func NewRootAudioPlayerNode(w World, input core.NewAudioPlayerInput) *AudioPlayerNode {
	tr := &AudioPlayerNode{
		mObjectContainer: &mObjectContainer{
			mObject: &mObject{
				e: w.NewEntity(),
				w: w,
			},
		},
	}
	core.SetAudioPlayerComponentData(w, tr.e, core.NewAudioPlayer(input))
	tr.wap = core.WatchAudioPlayerComponentData(w, tr.Entity())
	return tr
}

func NewChildAudioPlayerNode(parent ObjectContainer, input core.NewAudioPlayerInput) *AudioPlayerNode {
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
	core.SetAudioPlayerComponentData(tr.w, tr.e, core.NewAudioPlayer(input))
	tr.wap = core.WatchAudioPlayerComponentData(tr.w, tr.Entity())
	tr.SetParent(parent)
	return tr
}

func (t *AudioPlayerNode) AudioPlayer() *core.AudioPlayer {
	return t.wap.Data()
}

func (t *AudioPlayerNode) Destroy() {
	t.wap.Data().Pause()
	t.wap = nil
	t.mObjectContainer.Destroy()
}

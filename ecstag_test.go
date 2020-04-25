package tau

import (
	"testing"

	"github.com/hajimehoshi/ebiten"
	"github.com/stretchr/testify/assert"
)

func TestTagSystem(t *testing.T) {
	e := NewEngine(&NewEngineInput{
		Scale:  1,
		Width:  400,
		Height: 300,
	})
	w := NewWorld(e)
	w.Set(DefaultImageOptions, &ebiten.DrawImageOptions{})
	SetupSystem(w, TagCS)
	c := w.Component(CNTag)

	ents := w.NewEntities(10)

	w.AddComponentToEntity(ents[0], c, &Tag{
		Tags:  []string{"hard", "orange"},
		Dirty: true,
	})
	w.AddComponentToEntity(ents[1], c, &Tag{
		Tags:  []string{"hard", "orange", "fast"},
		Dirty: true,
	})
	w.AddComponentToEntity(ents[2], c, &Tag{
		Tags:  []string{"hard", "orange", "slow"},
		Dirty: true,
	})
	w.AddComponentToEntity(ents[3], c, &Tag{
		Tags:  []string{"soft", "orange", "slow"},
		Dirty: true,
	})
	w.AddComponentToEntity(ents[4], c, &Tag{
		Tags:  []string{"soft", "orange", "fast"},
		Dirty: true,
	})
	w.AddComponentToEntity(ents[5], c, &Tag{
		Tags:  []string{"orange", "fast"},
		Dirty: true,
	})
	w.AddComponentToEntity(ents[6], c, &Tag{
		Tags:  []string{"orange"},
		Dirty: true,
	})
	w.AddComponentToEntity(ents[7], c, &Tag{
		Tags:  []string{"black", "fast", "soft"},
		Dirty: true,
	})
	w.AddComponentToEntity(ents[8], c, &Tag{
		Tags:  []string{"soft"},
		Dirty: true,
	})
	w.AddComponentToEntity(ents[9], c, &Tag{
		Tags:  []string{"black", "black", "black"},
		Dirty: true,
	})
	assert.Empty(t, FindWithTag(w, "soft")) // cache is nil at this point
	w.Run(1)
	assert.Equal(t, 7, len(FindWithTag(w, "orange")))
	tt := c.Data(ents[9]).(*Tag)
	assert.Equal(t, 1, len(tt.Tags))

	// remove a tag, but be sloppy and dont change the Dirty tag:
	zent := c.Data(ents[0]).(*Tag)
	println(zent.Tags)
	zent.Tags = []string{"tangerine"}
	println(zent.Tags)
	w.Run(1)
	assert.Equal(t, 7, len(FindWithTag(w, "orange")))
	// run a lot
	println("OKDOK1")
	for i := 0; i < 30; i++ {
		w.Run(1)
	}
	println("OKDOK2")
	assert.Equal(t, 6, len(FindWithTag(w, "orange")))

	t8 := &Tag{}
	t8.Add("cook")
	assert.Equal(t, "cook", t8.Tags[0])
	t8.Add("benjen")
	t8.Remove("cook")
	assert.Equal(t, "benjen", t8.Tags[0])
}

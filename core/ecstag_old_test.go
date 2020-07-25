package core_test

// import (
// 	"testing"

// 	"github.com/gabstv/primen"
// 	"github.com/gabstv/primen/core"
// 	"github.com/gabstv/ebiten"
// 	"github.com/stretchr/testify/assert"
// )

// func TestTagSystem(t *testing.T) {
// 	e := primen.NewEngine(&primen.NewEngineInput{
// 		Scale:  1,
// 		Width:  400,
// 		Height: 300,
// 	})
// 	w := core.NewWorld(e)
// 	w.Set(core.DefaultImageOptions, &ebiten.DrawImageOptions{})
// 	core.SetupSystem(w, core.TagCS)
// 	c := w.Component(core.CNTag)

// 	ents := w.NewEntities(10)

// 	w.AddComponentToEntity(ents[0], c, &core.Tag{
// 		Tags:  []string{"hard", "orange"},
// 		Dirty: true,
// 	})
// 	w.AddComponentToEntity(ents[1], c, &core.Tag{
// 		Tags:  []string{"hard", "orange", "fast"},
// 		Dirty: true,
// 	})
// 	w.AddComponentToEntity(ents[2], c, &core.Tag{
// 		Tags:  []string{"hard", "orange", "slow"},
// 		Dirty: true,
// 	})
// 	w.AddComponentToEntity(ents[3], c, &core.Tag{
// 		Tags:  []string{"soft", "orange", "slow"},
// 		Dirty: true,
// 	})
// 	w.AddComponentToEntity(ents[4], c, &core.Tag{
// 		Tags:  []string{"soft", "orange", "fast"},
// 		Dirty: true,
// 	})
// 	w.AddComponentToEntity(ents[5], c, &core.Tag{
// 		Tags:  []string{"orange", "fast"},
// 		Dirty: true,
// 	})
// 	w.AddComponentToEntity(ents[6], c, &core.Tag{
// 		Tags:  []string{"orange"},
// 		Dirty: true,
// 	})
// 	w.AddComponentToEntity(ents[7], c, &core.Tag{
// 		Tags:  []string{"black", "fast", "soft"},
// 		Dirty: true,
// 	})
// 	w.AddComponentToEntity(ents[8], c, &core.Tag{
// 		Tags:  []string{"soft"},
// 		Dirty: true,
// 	})
// 	w.AddComponentToEntity(ents[9], c, &core.Tag{
// 		Tags:  []string{"black", "black", "black"},
// 		Dirty: true,
// 	})
// 	assert.Empty(t, core.FindWithTag(w, "soft")) // cache is nil at this point
// 	w.Run(1)
// 	assert.Equal(t, 7, len(core.FindWithTag(w, "orange")))
// 	tt := c.Data(ents[9]).(*core.Tag)
// 	assert.Equal(t, 1, len(tt.Tags))

// 	// remove a tag, but be sloppy and dont change the Dirty tag:
// 	zent := c.Data(ents[0]).(*core.Tag)
// 	println(zent.Tags)
// 	zent.Tags = []string{"tangerine"}
// 	println(zent.Tags)
// 	w.Run(1)
// 	assert.Equal(t, 7, len(core.FindWithTag(w, "orange")))
// 	// run a lot
// 	println("OKDOK1")
// 	for i := 0; i < 30; i++ {
// 		w.Run(1)
// 	}
// 	println("OKDOK2")
// 	assert.Equal(t, 6, len(core.FindWithTag(w, "orange")))

// 	t8 := &core.Tag{}
// 	t8.Add("cook")
// 	assert.Equal(t, "cook", t8.Tags[0])
// 	t8.Add("benjen")
// 	t8.Remove("cook")
// 	assert.Equal(t, "benjen", t8.Tags[0])
// }

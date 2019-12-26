package utils

import (
	"github.com/gabstv/troupe"
	"github.com/gabstv/troupe/utils/sets"
	"github.com/gabstv/troupe/utils/smid"
	"github.com/hajimehoshi/ebiten"
)

var tagPresent = struct{}{}

// Tag is the data of a tag component.
type Tag struct {
	// public and private struct fields
	Tags  []string
	Dirty bool

	lastTags []string
	lastMap  map[string]struct{}
}

func (t *Tag) Add(tag string) {
	t.Tags = append(t.Tags, tag)
	t.Dirty = true
}

func (t *Tag) Remove(tag string) {
	ri := -1
	for i, v := range t.Tags {
		if v == tag {
			ri = i
			break
		}
	}
	if ri == -1 {
		return
	}
	t.Tags = append(t.Tags[:ri], t.Tags[ri+1:]...)
	t.Dirty = true
}

// TagComponent will get the registered tag component of the world.
// If a component is not present, it will create a new component
// using world.NewComponent
func TagComponent(w troupe.WorldDicter) *troupe.Component {
	c := w.Component("troupe/utils.TagComponent")
	if c == nil {
		var err error
		c, err = w.NewComponent(troupe.NewComponentInput{
			Name: "troupe/utils.TagComponent",
			ValidateDataFn: func(data interface{}) bool {
				if data == nil {
					return false
				}
				_, ok := data.(*Tag)
				return ok
			},
			DestructorFn: func(_ troupe.WorldDicter, entity troupe.Entity, data interface{}) {
				if t, _ := data.(*Tag); t != nil {
					t.Dirty = false
					t.lastTags = nil
					t.lastMap = nil
				}
			},
		})
		if err != nil {
			panic(err)
		}
	}
	return c
}

// TagSystem creates the tag system
func TagSystem(w *troupe.World) *troupe.System {
	if sys := w.System("troupe/utils.TagSystem"); sys != nil {
		return sys
	}
	sys := w.NewSystem("troupe/utils.TagSystem", 0, troupe.SysWrapFn(TagSystemExec, tagSystemMidDirty(), smid.SkipFrames(30)), TagComponent(w))
	sys.AddTag(troupe.WorldTagUpdate)
	return sys
}

func FindWithTag(w *troupe.World, tag string, tags ...string) []troupe.Entity {
	if tag == "" {
		return []troupe.Entity{}
	}
	sys := TagSystem(w)
	ci := sys.Get("cache")
	if ci == nil {
		return []troupe.Entity{}
	}
	c := ci.(*tagSystemBakeCache)
	if c.Sets[tag] == nil {
		return []troupe.Entity{}
	}
	ents := sets.NewEntitySet(c.Sets[tag].Values()...)
	if len(tags) < 1 {
		return ents.Values()
	}
	for _, tt := range tags {
		en2 := c.Sets[tt]
		if en2 == nil {
			return []troupe.Entity{}
		}
		el := ents.Values()
		for _, k := range el {
			if !en2.Contains(k) {
				ents.Remove(k)
			}
		}
		if ents.Empty() {
			return []troupe.Entity{}
		}
	}
	return ents.Values()
}

func tagSystemMidDirty() troupe.SystemMiddleware {
	return func(next troupe.SystemFn) troupe.SystemFn {
		return func(ctx troupe.Context, screen *ebiten.Image) {
			defer next(ctx, screen)
			c := TagComponent(ctx.World())
			matches := ctx.System().View().Matches()
			dtags := make([]*Tag, 0, 8)
			dents := make([]troupe.Entity, 0, 8)
			for _, m := range matches {
				t := m.Components[c].(*Tag)
				if t.Dirty {
					dtags = append(dtags, t)
					dents = append(dents, m.Entity)
				}
			}
			tagSystemBake(ctx, screen, dents, dtags)
		}
	}
}

func isStringSliceDirty(current, previous []string) bool {
	if previous == nil {
		return true
	}
	if len(current) != len(previous) {
		return true
	}
	for i := range current {
		if current[i] != previous[i] {
			return true
		}
	}
	return false
}

// TagSystemExec is the main function of the TagSystem
func TagSystemExec(ctx troupe.Context, screen *ebiten.Image) {
	v := ctx.System().View()
	matches := v.Matches()
	tagcomp := TagComponent(ctx.World())
	dtags := make([]*Tag, 0)
	dents := make([]troupe.Entity, 0)
	for _, m := range matches {
		t := m.Components[tagcomp].(*Tag)
		if isStringSliceDirty(t.Tags, t.lastTags) {
			dtags = append(dtags, t)
			dents = append(dents, m.Entity)
		}
	}
	tagSystemBake(ctx, screen, dents, dtags)
}

type tagBuf struct {
	items []string
	n     int
}

func newTagBuf(size int) *tagBuf {
	if size <= 0 {
		size = 256
	}
	return &tagBuf{
		items: make([]string, size),
		n:     0,
	}
}

func (b *tagBuf) Reset() {
	b.n = 0
}

func (b *tagBuf) Add(v ...string) {
	for _, vi := range v {
		b.add(vi)
	}
}

func (b *tagBuf) add(v string) {
	if len(b.items) <= b.n {
		b.items = append(b.items, v)
		b.n++
		return
	}
	b.items[b.n] = v
	b.n++
}

func (b *tagBuf) List() []string {
	return b.items[:b.n]
}

type tagSystemBakeCache struct {
	AddBuf *tagBuf
	DelBuf *tagBuf
	Sets   map[string]sets.EntitySet
}

func tagSystemBake(ctx troupe.Context, screen *ebiten.Image, dentities []troupe.Entity, dtags []*Tag) {
	if len(dentities) < 1 {
		return
	}
	var cache *tagSystemBakeCache
	if x := ctx.System().Get("cache"); x != nil {
		cache = x.(*tagSystemBakeCache)
	} else {
		cache = &tagSystemBakeCache{
			AddBuf: newTagBuf(32),
			DelBuf: newTagBuf(32),
			Sets:   make(map[string]sets.EntitySet),
		}
		ctx.System().Set("cache", cache)
	}

	for i := range dentities {
		cache.AddBuf.Reset()
		cache.DelBuf.Reset()
		t := dtags[i]
		t.Dirty = false
		if t.lastMap == nil {
			t.lastMap = make(map[string]struct{})
			t.lastTags = make([]string, 0, len(t.Tags))
		}
		curMap := make(map[string]struct{})
		for _, tt := range t.Tags {
			if _, ok := t.lastMap[tt]; !ok {
				cache.AddBuf.Add(tt)
			}
			//TODO: autocorrect duplicated tags
			curMap[tt] = tagPresent
		}
		for _, tt := range t.lastTags {
			if _, ok := curMap[tt]; !ok {
				cache.DelBuf.Add(tt)
			}
		}
		if len(curMap) != len(t.Tags) {
			// autocorrect duplicated tags
			t.Tags = make([]string, 0, len(curMap))
			for k := range curMap {
				t.Tags = append(t.Tags, k)
			}
		}
		dels := cache.DelBuf.List()
		for _, k := range dels {
			if s := cache.Sets[k]; s != nil {
				s.Remove(dentities[i])
			}
		}
		adds := cache.AddBuf.List()
		for _, k := range adds {
			v := cache.Sets[k]
			if v == nil {
				cache.Sets[k] = sets.NewEntitySet(dentities[i])
			} else {
				v.Add(dentities[i])
			}
		}
		t.lastMap = curMap
		t.lastTags = make([]string, len(t.Tags))
		copy(t.lastTags, t.Tags)
	}
}

func init() {
	troupe.DefaultComp(func(e *troupe.Engine, w *troupe.World) {
		TagComponent(w)
	})
	troupe.DefaultSys(func(e *troupe.Engine, w *troupe.World) {
		TagSystem(w)
	})
}

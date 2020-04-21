package utils

import (
	"github.com/gabstv/ecs"
	"github.com/gabstv/tau"
	"github.com/gabstv/tau/utils/sets"
	"github.com/gabstv/tau/utils/smid"
)

const (
	SNTag = "tau/utils.TagSystem"
	CNTag = "tau/utils.TagComponent"
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

var (
	TagCS *TagComponentSystem = new(TagComponentSystem)
)

type TagComponentSystem struct {
	tau.BaseComponentSystem
}

func (cs *TagComponentSystem) SystemName() string {
	return SNTag
}

func (cs *TagComponentSystem) SystemInit() tau.SystemInitFn {
	return func(w *ecs.World, sys *ecs.System) {

	}
}

func (cs *TagComponentSystem) SystemExec() tau.SystemExecFn {
	return tau.SystemWrap(TagSystemExec, tagSystemMidDirty(), smid.SkipFrames(30))
}

func (cs *TagComponentSystem) Components(w ecs.Worlder) []*ecs.Component {
	return []*ecs.Component{
		tau.UpsertComponent(w, ecs.NewComponentInput{
			Name: CNTag,
			ValidateDataFn: func(data interface{}) bool {
				if data == nil {
					return false
				}
				_, ok := data.(*Tag)
				return ok
			},
			DestructorFn: func(w ecs.WorldDicter, entity ecs.Entity, data interface{}) {
				if t, _ := data.(*Tag); t != nil {
					t.Dirty = false
					t.lastTags = nil
					t.lastMap = nil
				}
			},
		}),
	}
}

func FindWithTag(w *ecs.World, tag string, tags ...string) []ecs.Entity {
	if tag == "" {
		return []ecs.Entity{}
	}
	sys := w.System(SNTag)
	ci := sys.Get("cache")
	if ci == nil {
		return []ecs.Entity{}
	}
	c := ci.(*tagSystemBakeCache)
	if c.Sets[tag] == nil {
		return []ecs.Entity{}
	}
	ents := sets.NewEntitySet(c.Sets[tag].Values()...)
	if len(tags) < 1 {
		return ents.Values()
	}
	for _, tt := range tags {
		en2 := c.Sets[tt]
		if en2 == nil {
			return []ecs.Entity{}
		}
		el := ents.Values()
		for _, k := range el {
			if !en2.Contains(k) {
				ents.Remove(k)
			}
		}
		if ents.Empty() {
			return []ecs.Entity{}
		}
	}
	return ents.Values()
}

func tagSystemMidDirty() tau.Middleware {
	return func(next tau.SystemExecFn) tau.SystemExecFn {
		return func(ctx tau.Context) {
			defer next(ctx)
			c := ctx.World().Component(CNTag)
			matches := ctx.System().View().Matches()
			dtags := make([]*Tag, 0, 8)
			dents := make([]ecs.Entity, 0, 8)
			for _, m := range matches {
				t := m.Components[c].(*Tag)
				if t.Dirty {
					dtags = append(dtags, t)
					dents = append(dents, m.Entity)
				}
			}
			tagSystemBake(ctx, dents, dtags)
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
func TagSystemExec(ctx tau.Context) {
	v := ctx.System().View()
	matches := v.Matches()
	tagcomp := ctx.World().Component(CNTag)
	dtags := make([]*Tag, 0)
	dents := make([]ecs.Entity, 0)
	for _, m := range matches {
		t := m.Components[tagcomp].(*Tag)
		if isStringSliceDirty(t.Tags, t.lastTags) {
			dtags = append(dtags, t)
			dents = append(dents, m.Entity)
		}
	}
	tagSystemBake(ctx, dents, dtags)
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

func tagSystemBake(ctx tau.Context, dentities []ecs.Entity, dtags []*Tag) {
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
	tau.RegisterComponentSystem(&TagComponentSystem{})
}

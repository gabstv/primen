package core

import (
	"sort"
	"sync"

	"github.com/gabstv/ecs/v2"
)

type SLVal interface {
	Less(v interface{}) bool
	Destroy()
}

type EntitySortedList struct {
	l      sync.RWMutex
	set    map[ecs.Entity]*entitySortedListEntry
	list   []*entitySortedListEntry
	mincap int
}

type entitySortedListEntry struct {
	Key   ecs.Entity
	Value SLVal
	Index int
}

func NewEntitySortedList(mincap int) *EntitySortedList {
	return &EntitySortedList{
		set:    make(map[ecs.Entity]*entitySortedListEntry),
		list:   make([]*entitySortedListEntry, 0, mincap),
		mincap: mincap,
	}
}

func (sl *EntitySortedList) AddOrUpdate(key ecs.Entity, value SLVal) bool {
	sl.l.RLock()
	e, ok := sl.set[key]
	sl.l.RUnlock()
	if ok {
		if e.Value == value {
			return false
		}
		// update and sort
		e.Value.Destroy()
		e.Value = value
		sl.l.Lock()
		sl.sort()
		sl.l.Unlock()
		return false
	}
	sl.l.Lock()
	item := &entitySortedListEntry{
		Key:   key,
		Value: value,
		Index: len(sl.list),
	}
	sl.set[key] = item
	sl.list = append(sl.list, item)
	sl.sort()
	sl.l.Unlock()
	return true
}

func (sl *EntitySortedList) Get(key ecs.Entity) (SLVal, bool) {
	sl.l.RLock()
	defer sl.l.RUnlock()
	if v, ok := sl.set[key]; ok {
		return v.Value, true
	}
	return nil, false
}

func (sl *EntitySortedList) Delete(key ecs.Entity) bool {
	sl.l.Lock()
	defer sl.l.Unlock()
	vv, ok := sl.set[key]
	if !ok {
		return false
	}
	for i := vv.Index + 1; i < len(sl.list); i++ {
		sl.list[i].Index = i - 1
	}
	sl.list = append(sl.list[:vv.Index], sl.list[vv.Index+1:]...)
	sl.set[key].Value.Destroy()
	delete(sl.set, key)
	return true
}

func (sl *EntitySortedList) LastValue() SLVal {
	sl.l.RLock()
	defer sl.l.RUnlock()
	if len(sl.list) < 1 {
		return nil
	}
	return sl.list[len(sl.list)-1].Value
}

func (sl *EntitySortedList) FirstValue() SLVal {
	sl.l.RLock()
	defer sl.l.RUnlock()
	if len(sl.list) < 1 {
		return nil
	}
	return sl.list[0].Value
}

// Reset clears the list.
func (sl *EntitySortedList) Reset() {
	sl.l.Lock()
	defer sl.l.Unlock()
	sl.set = make(map[ecs.Entity]*entitySortedListEntry)
	sl.list = make([]*entitySortedListEntry, 0, sl.mincap)
}

// Each iterates on all entries, sorted.
//
// Do not delete or add entries while looping (it is read locked while looping).
func (sl *EntitySortedList) Each(fn func(key ecs.Entity, value SLVal) bool) {
	sl.l.RLock()
	defer sl.l.RUnlock()
	for _, v := range sl.list {
		if !fn(v.Key, v.Value) {
			return
		}
	}
}

func (sl *EntitySortedList) sort() {
	sort.Sort(entitySortedListSt(sl.list))
}

type entitySortedListSt []*entitySortedListEntry

func (a entitySortedListSt) Len() int { return len(a) }
func (a entitySortedListSt) Swap(i, j int) {
	a[i].Index, a[j].Index = j, i
	a[i], a[j] = a[j], a[i]
}
func (a entitySortedListSt) Less(i, j int) bool {
	a[i].Index = i
	a[j].Index = j
	return a[i].Value.Less(a[j].Value)
}

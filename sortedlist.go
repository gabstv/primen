package tau

import (
	"sort"
	"sync"

	"github.com/gabstv/ecs"
)

type SortedListLess func(ival, jval interface{}) bool

// SortedListCompare should return TRUE if haystackitemval is less or equal to needleval
type SortedListCompare func(needleval, haystackitemval interface{}) bool

type IntSortedList struct {
	l         sync.RWMutex
	set       map[int64]*intSortedListEntry
	list      []*intSortedListEntry
	lessfn    SortedListLess
	comparefn SortedListCompare
	mincap    int
}

type intSortedListEntry struct {
	Key   int64
	Value interface{}
}

func NewIntSortedList(less SortedListLess, comp SortedListCompare, mincap int) *IntSortedList {
	return &IntSortedList{
		set:       make(map[int64]*intSortedListEntry),
		list:      make([]*intSortedListEntry, 0, mincap),
		lessfn:    less,
		comparefn: comp,
		mincap:    mincap,
	}
}

func (sl *IntSortedList) AddOrUpdate(key int64, value interface{}) bool {
	sl.l.RLock()
	e, ok := sl.set[key]
	sl.l.RUnlock()
	if ok {
		if e.Value == value {
			return false
		}
		// update and sort
		e.Value = value
		sl.sort()
		return false
	}
	item := &intSortedListEntry{
		Key:   key,
		Value: value,
	}
	sl.l.Lock()
	sl.set[key] = item
	sl.list = append(sl.list, item)
	sl.sort()
	sl.l.Unlock()
	return true
}

// Find returns the key and value of the item.
func (sl *IntSortedList) Find(val interface{}) (int64, bool) {
	sl.l.RLock()
	defer sl.l.RUnlock()
	n := sort.Search(len(sl.list), func(i int) bool {
		return sl.comparefn(val, sl.list[i].Value)
	})
	if n == len(sl.list) {
		return 0, false
	}
	if sl.list[n].Value != val {
		return 0, false
	}
	return sl.list[n].Key, true
}

func (sl *IntSortedList) Get(key int64) (interface{}, bool) {
	sl.l.RLock()
	defer sl.l.RUnlock()
	if v, ok := sl.set[key]; ok {
		return v.Value, true
	}
	return nil, false
}

func (sl *IntSortedList) Delete(key int64) bool {
	sl.l.Lock()
	defer sl.l.Unlock()
	vv, ok := sl.set[key]
	if !ok {
		return false
	}
	n := sort.Search(len(sl.list), func(i int) bool {
		return sl.comparefn(vv.Value, sl.list[i].Value)
	})
	if n == len(sl.list) {
		return false
	}
	if sl.list[n].Key != key {
		return false
	}
	sl.list = append(sl.list[:n], sl.list[n+1:]...)
	delete(sl.set, key)
	return true
}

func (sl *IntSortedList) LastValue() interface{} {
	sl.l.RLock()
	defer sl.l.RUnlock()
	if len(sl.list) < 1 {
		return nil
	}
	return sl.list[len(sl.list)-1].Value
}

func (sl *IntSortedList) FirstValue() interface{} {
	sl.l.RLock()
	defer sl.l.RUnlock()
	if len(sl.list) < 1 {
		return nil
	}
	return sl.list[0].Value
}

// Reset clears the list.
func (sl *IntSortedList) Reset() {
	sl.l.Lock()
	defer sl.l.Unlock()
	sl.set = make(map[int64]*intSortedListEntry)
	sl.list = make([]*intSortedListEntry, 0, sl.mincap)
}

// Each iterates on all entries, sorted.
//
// Do not delete or add entries while looping (it is read locked while looping).
func (sl *IntSortedList) Each(fn func(key int64, value interface{}) bool) {
	sl.l.RLock()
	defer sl.l.RUnlock()
	for _, v := range sl.list {
		if !fn(v.Key, v.Value) {
			return
		}
	}
}

func (sl *IntSortedList) sort() {
	sort.Slice(sl.list, sl.sortfn())
}

func (sl *IntSortedList) sortfn() func(i, j int) bool {
	l := sl.list
	return func(i, j int) bool {
		ie := l[i]
		je := l[j]
		return sl.lessfn(ie.Value, je.Value)
	}
}

//
//
//
//
//
//
//
//
//

type EntitySortedList struct {
	l         sync.RWMutex
	set       map[ecs.Entity]*entitySortedListEntry
	list      []*entitySortedListEntry
	lessfn    SortedListLess
	comparefn SortedListCompare
	mincap    int
}

type entitySortedListEntry struct {
	Key   ecs.Entity
	Value interface{}
}

func NewEntitySortedList(less SortedListLess, comp SortedListCompare, mincap int) *EntitySortedList {
	return &EntitySortedList{
		set:       make(map[ecs.Entity]*entitySortedListEntry),
		list:      make([]*entitySortedListEntry, 0, mincap),
		lessfn:    less,
		comparefn: comp,
		mincap:    mincap,
	}
}

func (sl *EntitySortedList) AddOrUpdate(key ecs.Entity, value interface{}) bool {
	sl.l.RLock()
	e, ok := sl.set[key]
	sl.l.RUnlock()
	if ok {
		if e.Value == value {
			return false
		}
		// update and sort
		e.Value = value
		sl.sort()
		return false
	}
	item := &entitySortedListEntry{
		Key:   key,
		Value: value,
	}
	sl.l.Lock()
	sl.set[key] = item
	sl.list = append(sl.list, item)
	sl.sort()
	sl.l.Unlock()
	return true
}

// Find returns the key and value of the item.
func (sl *EntitySortedList) Find(val interface{}) (ecs.Entity, bool) {
	sl.l.RLock()
	defer sl.l.RUnlock()
	n := sort.Search(len(sl.list), func(i int) bool {
		return sl.comparefn(val, sl.list[i].Value)
	})
	if n == len(sl.list) {
		return 0, false
	}
	if sl.list[n].Value != val {
		return 0, false
	}
	return sl.list[n].Key, true
}

func (sl *EntitySortedList) Get(key ecs.Entity) (interface{}, bool) {
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
	n := sort.Search(len(sl.list), func(i int) bool {
		return sl.comparefn(vv.Value, sl.list[i].Value)
	})
	if n == len(sl.list) {
		return false
	}
	if sl.list[n].Key != key {
		return false
	}
	sl.list = append(sl.list[:n], sl.list[n+1:]...)
	delete(sl.set, key)
	return true
}

func (sl *EntitySortedList) LastValue() interface{} {
	sl.l.RLock()
	defer sl.l.RUnlock()
	if len(sl.list) < 1 {
		return nil
	}
	return sl.list[len(sl.list)-1].Value
}

func (sl *EntitySortedList) FirstValue() interface{} {
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
func (sl *EntitySortedList) Each(fn func(key ecs.Entity, value interface{}) bool) {
	sl.l.RLock()
	defer sl.l.RUnlock()
	for _, v := range sl.list {
		if !fn(v.Key, v.Value) {
			return
		}
	}
}

func (sl *EntitySortedList) sort() {
	sort.Slice(sl.list, sl.sortfn())
}

func (sl *EntitySortedList) sortfn() func(i, j int) bool {
	l := sl.list
	return func(i, j int) bool {
		ie := l[i]
		je := l[j]
		return sl.lessfn(ie.Value, je.Value)
	}
}

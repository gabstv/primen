// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sets

import (
	"bytes"
	"fmt"

	"github.com/gabstv/tau"
)

type EntitySet interface {
	Add(elements ...tau.Entity)
	Remove(elements ...tau.Entity)
	Contains(elements ...tau.Entity) bool
	Values() []tau.Entity
	SetBase
}

type entitySet struct {
	items map[tau.Entity]struct{}
}

// NewEntitySet instantiates a new empty set and adds the entities (if present)
func NewEntitySet(values ...tau.Entity) EntitySet {
	set := &entitySet{items: make(map[tau.Entity]struct{})}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

// Add adds the items (one or more) to the set.
func (set *entitySet) Add(items ...tau.Entity) {
	for _, item := range items {
		set.items[item] = itemExists
	}
}

// Remove removes the items (one or more) from the set.
func (set *entitySet) Remove(items ...tau.Entity) {
	for _, item := range items {
		delete(set.items, item)
	}
}

// Contains check if items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (set *entitySet) Contains(items ...tau.Entity) bool {
	for _, item := range items {
		if _, contains := set.items[item]; !contains {
			return false
		}
	}
	return true
}

// Empty returns true if set does not contain any elements.
func (set *entitySet) Empty() bool {
	return set.Size() == 0
}

// Size returns number of elements within the set.
func (set *entitySet) Size() int {
	return len(set.items)
}

// Clear clears all values in the set.
func (set *entitySet) Clear() {
	set.items = make(map[tau.Entity]struct{})
}

// Values returns all items in the set.
func (set *entitySet) Values() []tau.Entity {
	values := make([]tau.Entity, set.Size())
	count := 0
	for item := range set.items {
		values[count] = item
		count++
	}
	return values
}

// String returns a string representation of EntitySet
func (set *entitySet) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("EntitySet {")

	for k := range set.items {
		buf.WriteString(fmt.Sprintf("E%v; ", k))
	}
	buf.Truncate(buf.Len() - 2)
	buf.WriteRune('}')
	return buf.String()
}

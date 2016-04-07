// Copyright 2016 Pikkpoiss
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gamejam

import (
	"fmt"
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=eventobserverlist.go gen "Something=EventObserver"
//go:generate genny -in=$GOFILE -out=scenelist.go gen "Something=Scene"

type Something generic.Type
type SomethingListID int

var ErrSomethingNotInList = fmt.Errorf("Could not find item in list")

// A SomethingNode represents a doubly-linked Something.
type SomethingNode struct {
	Something
	id   SomethingListID
	next *SomethingNode
	prev *SomethingNode
}

func (n *SomethingNode) Next() *SomethingNode {
	return n.next
}

func (n *SomethingNode) NodeID() SomethingListID {
	return n.id
}

// Unlinks the node from the current list.
// Returns the next node for convenience while iterating.
func (n *SomethingNode) Unlink() (next *SomethingNode) {
	next = n.next
	if n.next != nil {
		n.next.prev = n.prev
	}
	if n.prev != nil {
		n.prev.next = n.next
	}
	n.next = nil
	n.prev = nil
	return
}

// A SomethingList manages a list of SomethingNode items.
type SomethingList struct {
	nextID SomethingListID
	head   *SomethingNode
}

func NewSomethingList(items ...Something) (l *SomethingList) {
	l = &SomethingList{
		nextID: 1,
	}
	for i := len(items) - 1; i >= 0; i-- {
		l.Prepend(items[i])
	}
	return
}

func (l *SomethingList) Head() *SomethingNode {
	return l.head
}

func (l *SomethingList) Prepend(item Something) (id SomethingListID) {
	id = l.nextID
	l.nextID++
	var node = &SomethingNode{
		Something: item,
		id:        id,
		next:      l.head,
		prev:      nil,
	}
	if l.head != nil {
		l.head.prev = node
	}
	l.head = node
	return
}

// Attempts to remove `item` from the list.
// Returns ErrSomethingNotInList if the ID did not exist in list.
func (l *SomethingList) Remove(id SomethingListID) (removed Something, err error) {
	var node = l.Head()
	for node != nil {
		if node.NodeID() == id {
			node.Unlink()
			removed = node.Something
			return // Should only ever be one of an ID in a list.
		}
		node = node.Next()
	}
	err = ErrSomethingNotInList
	return
}

// Unlinks all Something items from this SomethingList.
func (l *SomethingList) Empty() {
	var node = l.Head()
	for node != nil {
		node = node.Unlink()
	}
}

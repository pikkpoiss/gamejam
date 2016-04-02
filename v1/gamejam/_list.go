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

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=eventobserver$GOFILE gen "Something=EventObserver"

type Something generic.Type

// A SomethingNode represents a doubly-linked Something.
type SomethingNode struct {
	Something
	next *SomethingNode
	prev *SomethingNode
}

func (n *SomethingNode) Next() *SomethingNode {
	return n.next
}

func (n *SomethingNode) Unlink() {
	if n.next != nil {
		n.next.prev = n.prev
	}
	if n.prev != nil {
		n.prev.next = n.next
	}
	n.next = nil
	n.prev = nil
}

// A SomethingList manages a list of SomethingNode items.
type SomethingList struct {
	head *SomethingNode
}

func NewSomethingList(items ...Something) (l *SomethingList) {
	l = &SomethingList{}
	for i := len(items) - 1; i >= 0; i-- {
		l.Prepend(items[i])
	}
	return
}

func (l *SomethingList) Head() SomethingNode {
	return l.head
}

func (l *SomethingList) Prepend(item Something) {
	var item = &SomethingNode{
		Something: item,
		next:      l.head,
		prev:      nil,
	}
	if l.head != nil {
		l.head.prev = item
	}
	l.head = item
}

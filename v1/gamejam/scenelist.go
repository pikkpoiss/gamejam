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
)

var (
	ErrorNotInList = fmt.Errorf("Item not in a list")
)

type SceneListItem interface {
	Scene
	Next() SceneListItem
	Unlink() error
}

type SceneList interface {
	Head() SceneListItem
	Remove(id SceneID) (removed SceneListItem, err error)
	Prepend(s Scene)
}

type BaseSceneList struct {
	head *BaseSceneListItem
}

func NewBaseSceneList(scenes ...Scene) (l *BaseSceneList) {
	l = &BaseSceneList{}
	for i := len(scenes) - 1; i >= 0; i-- {
		l.Prepend(scenes[i])
	}
	return
}

func (l *BaseSceneList) Head() SceneListItem {
	return l.head
}

func (l *BaseSceneList) Prepend(s Scene) {
	var item = &BaseSceneListItem{
		Scene: s,
		list:  l,
		next:  l.head,
		prev:  nil,
	}
	if l.head != nil {
		l.head.prev = item
	}
	l.head = item
}

func (l *BaseSceneList) Remove(id SceneID) (removed SceneListItem, err error) {
	var item = l.Head()
	for item != nil {
		if item.GetID() == id {
			item.Unlink()
			removed = item
			return
		}
	}
	err = fmt.Errorf("Scene ID %v was not found in list", id)
	return
}

func (l *BaseSceneList) unlink(item *BaseSceneListItem) {
	if l.head == item {
		l.head = item.next
	}
}

type BaseSceneListItem struct {
	Scene
	next *BaseSceneListItem
	prev *BaseSceneListItem
	list *BaseSceneList
}

func (i *BaseSceneListItem) Next() SceneListItem {
	if i.next == nil {
		// Allow checking SceneListItem == nil, otherwise will be
		// a non-nil interface pointing to a nil value.
		return nil
	}
	return i.next
}

func (i *BaseSceneListItem) Unlink() error {
	if i.list == nil {
		return ErrorNotInList
	}
	i.list.unlink(i)
	i.list = nil
	i.next.prev = i.prev
	if i.prev != nil {
		i.prev.next = i.next
	}
	i.next = nil
	i.prev = nil
	return nil
}

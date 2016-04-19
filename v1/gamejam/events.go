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

type Event interface{}

type EventObserverID int

type EventObserver func(event Event)

type Events interface {
	AddEventObserver(obs EventObserver) (id EventObserverID)
	RemoveEventObserver(obs EventObserverID) (err error)
	Notify(event Event)
	DeleteObservers()
}

type BaseEvents struct {
	list *EventObserverList
}

func NewBaseEvents() *BaseEvents {
	return &BaseEvents{
		list: NewEventObserverList(),
	}
}

func (n *BaseEvents) AddEventObserver(obs EventObserver) (id EventObserverID) {
	id = EventObserverID(n.list.Prepend(obs).EventObserverListID())
	return
}

func (n *BaseEvents) RemoveEventObserver(id EventObserverID) (err error) {
	_, err = n.list.Remove(EventObserverListID(id))
	return
}

func (n *BaseEvents) DeleteObservers() {
	n.list.Empty()
}

func (n *BaseEvents) Notify(event Event) {
	var node = n.list.Head()
	for node != nil {
		node.EventObserver(event)
		node = node.Next()
	}
}

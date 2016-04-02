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

type EventNotifier interface {
	AddObserver(obs EventObserver) (err error)
	RemoveObserver(obs EventObserver) (err error)
	Delete()
}

type BaseEventNotifier struct {
	list *EventObserverList
}

func NewBaseEventNotifier() *BaseEventNotifier {
	return &BaseEventNotifier{
		list: NewEventObserverList(),
	}
}

func (n *BaseEventNotifier) AddObserver(obs EventObserver) (err error) {
	n.list.Prepend(obs)
	return
}

func (n *BaseEventNotifier) RemoveObserver(obs EventObserver) (err error) {
	err = n.list.Remove(obs)
	return
}

func (n *BaseEventNotifier) Delete() {
	n.list.Delete()
}

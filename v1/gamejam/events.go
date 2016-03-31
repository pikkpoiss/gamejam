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

type Event interface {
}

type EventListener interface {
	OnEvent(evt Event) (err error)
}

type EventTrigger interface {
	Trigger(evt Event) (err error)
	AddEventListener(listener EventListener) (err error)
	RemoveEventListener(listener EventListener) (err error)
	Delete()
}

type BaseEventTrigger struct {
}

func (t *BaseEventTrigger) Trigger(evt Event) (err error) {
	return
}

func (t *BaseEventTrigger) AddEventListener(listener EventListener) (err error) {
	return
}

func (t *BaseEventTrigger) RemoveEventListener(listener EventListener) (err error) {
	return
}

func (t *BaseEventTrigger) Delete() {
}

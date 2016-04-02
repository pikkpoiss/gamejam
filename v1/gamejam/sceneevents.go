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

type SceneCreateEvent struct {
}

type SceneDeleteEvent struct {
}

type SceneEventListener interface {
	OnSceneCreate(evt SceneCreateEvent) (err error)
	OnSceneDelete(evt SceneDeleteEvent) (err error)
}

type SceneEventTrigger interface {
	
}

type baseSceneEventListener struct {
	SceneEventListener
}

func (l *baseSceneEventListener) OnEvent(evt Event) (err error) {
	switch event := evt.(type) {
	case SceneCreateEvent:
		err = l.OnSceneCreate(event)
	case SceneDeleteEvent:
		err = l.OnSceneDelete(event)
	}
	return
}

func NewSceneEventListener(impl SceneEventListener) EventListener {
	return &baseSceneEventListener{
		SceneEventListener: impl,
	}
}

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

// Fired after the scene has finished loading all resources.
type SceneLoadedEvent struct {
	Scene
}

// Fired before the scene starts unloading all resources.
type SceneUnloadEvent struct {
	Scene
}

type SceneEventObserver interface {
	OnSceneLoaded(event SceneLoadedEvent)
	OnSceneUnload(event SceneUnloadEvent)
}

func BindSceneEvents(events Events, obs SceneEventObserver) (id EventObserverID) {
	id = events.AddEventObserver(func(evt Event) {
		switch event := evt.(type) {
		case SceneLoadedEvent:
			obs.OnSceneLoaded(event)
		case SceneUnloadEvent:
			obs.OnSceneUnload(event)
		}
		return
	})
	return
}

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

type SceneManager interface {
	Load(r Resources) (err error)
	AddScene(s Scene) (err error)
	RemoveScene(id SceneID) (err error)
	Head() SceneListItem
	Update()
	Render()
}

type BaseSceneManager struct {
	scenelist SceneList
	resources Resources
	nextID    SceneID
}

func NewBaseSceneManager(scenes ...Scene) (m *BaseSceneManager) {
	m = &BaseSceneManager{
		scenelist: NewBaseSceneList(scenes...),
		resources: nil,
		nextID:    1,
	}
	return
}

func (m *BaseSceneManager) Load(r Resources) (err error) {
	if r == nil {
		err = fmt.Errorf("Load called with nil Resources")
		return
	}
	m.resources = r
	var item = m.Head()
	for item != nil {
		if err = m.loadScene(item); err != nil {
			return
		}
		item = item.Next()
	}
	return
}

func (m *BaseSceneManager) Head() SceneListItem {
	return m.scenelist.Head()
}

func (m *BaseSceneManager) loadScene(s Scene) (err error) {
	s.SetID(m.nextID)
	m.nextID++
	err = s.OnLoad(m.resources)
	return
}

func (m *BaseSceneManager) AddScene(s Scene) (err error) {
	m.scenelist.Prepend(s)
	if m.resources == nil {
		// Load hasn't been called yet.
		return
	}
	err = m.loadScene(s)
	return
}

func (m *BaseSceneManager) RemoveScene(id SceneID) (err error) {
	var removed SceneListItem
	if removed, err = m.scenelist.Remove(id); err != nil {
		return
	}
	err = removed.OnUnload(m.resources)
	return
}

func (m *BaseSceneManager) Update() {
	var item = m.Head()
	for item != nil {
		item.Update(m)
		item = item.Next()
	}
}

func (m *BaseSceneManager) Render() {
	var item = m.Head()
	for item != nil {
		item.Render()
		item = item.Next()
	}
}

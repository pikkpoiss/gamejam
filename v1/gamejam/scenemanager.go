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
	AddScene(s Scene) (err error)
	RemoveScene(s Scene) (err error)
	Head() *SceneNode
	Update() (err error)
	Render() (err error)
	Delete() (err error)
}

type BaseSceneManager struct {
	removelist []SceneID
	scenelist  *SceneList
	resources  Resources
}

func NewBaseSceneManager(res Resources, scenes ...Scene) (m *BaseSceneManager, err error) {
	m = &BaseSceneManager{
		removelist: nil,
		scenelist:  NewSceneList(),
		resources:  res,
	}
	for i := len(scenes) - 1; i >= 0; i-- {
		if err = m.AddScene(scenes[i]); err != nil {
			return
		}
	}
	return
}

func (m *BaseSceneManager) Head() *SceneNode {
	return m.scenelist.Head()
}

func (m *BaseSceneManager) AddScene(s Scene) (err error) {
	//BindSceneEventObserver(s, m)
	if err = s.Load(m.resources); err != nil {
		return
	}
	var id = m.scenelist.Prepend(s)
	s.SetSceneID(SceneID(id))
	return
}

func (m *BaseSceneManager) RemoveScene(s Scene) (err error) {
	m.removelist = append(m.removelist, s.SceneID())
	return
}

func (m *BaseSceneManager) Update() (err error) {
	var (
		item  = m.Head()
		scene Scene
		id    SceneID
	)
	for item != nil {
		item.Update(m)
		item = item.Next()
	}
	if m.removelist != nil {
		for _, id = range m.removelist {
			// Inefficient.
			if scene, err = m.scenelist.Remove(SceneListID(id)); err != nil {
				return
			}
			if err = scene.Unload(m.resources); err != nil {
				return
			}
		}
		m.removelist = nil
	}
	return
}

func (m *BaseSceneManager) Render() (err error) {
	var item = m.Head()
	for item != nil {
		item.Render()
		item = item.Next()
	}
	return
}

func (m *BaseSceneManager) Delete() (err error) {
	fmt.Printf("BaseSceneManager Delete\n")
	var node = m.scenelist.Head()
	for node != nil {
		if err = node.Unload(m.resources); err != nil {
			return
		}
		node = node.Next()
	}
	m.scenelist.Empty()
	return
}

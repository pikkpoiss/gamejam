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

type SceneID int

type Scene interface {
	AddComponent(c Component)
	Load(r Resources) (err error)
	Unload(r Resources) (err error)
	Render()
	Update(mgr SceneManager)
	SetSceneID(id SceneID)
	SceneID() SceneID
}

type BaseScene struct {
	components map[ComponentID]Component
	id SceneID
}

func NewBaseScene() *BaseScene {
	return &BaseScene{
		components: map[ComponentID]Component{},
	}
}

func (s *BaseScene) AddComponent(c Component) {
	c.SetScene(s)
	s.components[c.GetID()] = c
}

func (s *BaseScene) Load(r Resources) (err error) {
	return
}

func (s *BaseScene) Render() {
}

func (s *BaseScene) SetSceneID(id SceneID) {
	s.id = id
}

func (s *BaseScene) SceneID() SceneID {
	return s.id
}

func (s *BaseScene) Unload(r Resources) (err error) {
	var (
		id ComponentID
		c  Component
	)
	for id, c = range s.components {
		s.components[id] = nil
		c.Delete()
	}
	//s.DeleteObservers()
	return
}

func (s *BaseScene) Update(mgr SceneManager) {
}

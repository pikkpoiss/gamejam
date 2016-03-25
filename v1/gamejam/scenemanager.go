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

type SceneManager interface {
	Init(r Resources) (err error)
	GetScene() Scene
	SetScene(s Scene) (err error)
}

type BaseSceneManager struct {
	scene     Scene
	resources Resources
}

func NewBaseSceneManager(s Scene) (m *BaseSceneManager) {
	m = &BaseSceneManager{
		scene:     s,
		resources: nil,
	}
	return
}

func (m *BaseSceneManager) Init(r Resources) (err error) {
	m.resources = r
	return m.SetScene(m.scene)
}

func (m *BaseSceneManager) GetScene() Scene {
	return m.scene
}

func (m *BaseSceneManager) SetScene(s Scene) (err error) {
	if err = s.Load(m.resources); err != nil {
		return
	}
	if m.scene != nil && m.scene != s {
		m.scene.Unload(m.resources)
	}
	m.scene = s
	return
}

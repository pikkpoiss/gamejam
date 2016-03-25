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
	SetScene(s Scene) (done chan error)
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
	err = <-m.SetScene(m.scene)
	return
}

func (m *BaseSceneManager) GetScene() Scene {
	return m.scene
}

func (m *BaseSceneManager) loadScene(scene Scene, done chan error) {
	var (
		doneScene = make(chan error, 1)
		err       error
	)
	scene.Load(m.resources, doneScene)
	err = <-doneScene
	if err == nil {
		if m.scene != nil && m.scene != scene {
			m.scene.Unload(m.resources)
		}
		m.scene = scene
	}
	done <- err
}

func (m *BaseSceneManager) SetScene(s Scene) chan error {
	var done = make(chan error)
	go m.loadScene(s, done)
	return done
}

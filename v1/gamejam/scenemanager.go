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
	Init(r Resources)
	GetScene() Scene
	SetScene(s Scene)
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
	m.SetScene(s)
	return
}

func (m *BaseSceneManager) Init(r Resources) {
	m.resources = r
	m.SetScene(m.scene)
}

func (m *BaseSceneManager) GetScene() Scene {
	return m.scene
}

func (m *BaseSceneManager) loadScene(scene Scene) {
	var (
		done = make(chan error)
		err  error
	)
	scene.Load(m.resources, done)
	err = <-done
	if err != nil {
		panic(err) // TODO: Better support for propagating errors?
	}
	m.scene = scene
}

func (m *BaseSceneManager) SetScene(s Scene) {
	go m.loadScene(s)
}

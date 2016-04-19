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

package main

import (
	"flag"
	"fmt"
	"github.com/pikkpoiss/gamejam/v1/gamejam"
	"runtime"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

func main() {
	var (
		app  gamejam.App
		main *gamejam.Main
		err  error
	)
	flag.Parse()
	if app, err = NewApp(); err != nil {
		panic(err)
	}
	main = gamejam.NewMain(app)
	if err = main.Run(); err != nil {
		panic(err)
	}
}

type App struct {
}

func NewApp() (app *App, err error) {
	app = &App{}
	return
}

func (a *App) GetWindowData() (data *gamejam.WindowData, err error) {
	data = &gamejam.WindowData{
		WindowWidth:  640,
		WindowHeight: 480,
		WindowTitle:  "03-scene",
	}
	return
}

func (a *App) GetAppData() (data *gamejam.AppData, err error) {
	data = &gamejam.AppData{
		Resources: gamejam.NewBaseResources(),
	}
	data.SceneManager, err = gamejam.NewBaseSceneManager(data.Resources, NewScene())
	return
}

type Scene struct {
	*gamejam.BaseScene
	count int
}

func NewScene() *Scene {
	return &Scene{
		BaseScene: gamejam.NewBaseScene(),
	}
}

func (s *Scene) Load(r gamejam.Resources) (err error) {
	if err = s.BaseScene.Load(r); err != nil {
		return
	}
	var sheet gamejam.ResourceType
	if sheet, err = r.Get(gamejam.NewTexturePackerSheetLoader(
		"./examples/resources/spritesheet.json",
		gamejam.SmoothingNearest,
	)); err != nil {
		return
	}
	fmt.Printf("LOADED %v\n", sheet)
	return
}

func (s *Scene) Render() {
	fmt.Printf("\rRENDER %v", s.count)
	s.count++
}

func (s *Scene) Unload(res gamejam.Resources) (err error) {
	fmt.Println("DELETE")
	err = s.BaseScene.Unload(res)
	return
}

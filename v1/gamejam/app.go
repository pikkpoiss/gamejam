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
	"github.com/golang/glog"
	"github.com/pikkpoiss/gamejam/v1/base/core"
)

type App interface {
	GetWindowData() (data *WindowData, err error)
	GetAppData() (data *AppData, err error)
}

type WindowData struct {
	WindowWidth  int
	WindowHeight int
	WindowTitle  string
}

type AppData struct {
	SceneManager SceneManager
	Resources    Resources
}

type Main struct {
	App
}

func NewMain(app App) *Main {
	return &Main{
		App: app,
	}
}

func (a *Main) Run() (err error) {
	var (
		context *core.Context
		winData *WindowData
		appData *AppData
	)
	if context, err = core.NewContext(); err != nil {
		return
	}
	defer context.Delete()
	defer glog.Flush()
	if winData, err = a.GetWindowData(); err != nil {
		return
	}
	if err = context.CreateWindow(
		winData.WindowWidth,
		winData.WindowHeight,
		winData.WindowTitle,
	); err != nil {
		return
	}
	if appData, err = a.GetAppData(); err != nil {
		return
	}
	if appData.SceneManager == nil {
		err = fmt.Errorf("SceneManager must be set")
		return
	}
	if appData.Resources == nil {
		err = fmt.Errorf("Resources must be set")
		return
	}
	defer appData.SceneManager.Delete()
	for !context.ShouldClose() {
		context.Events.Poll()
		context.Clear()
		appData.SceneManager.Update()
		appData.SceneManager.Render()
	}
	return
}

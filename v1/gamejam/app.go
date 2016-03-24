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
	"github.com/golang/glog"
	"github.com/pikkpoiss/gamejam/v1/base/core"
)

type App interface {
	Update()
	Render()
	Run() (err error)
}

type BaseApp struct {
	WindowWidth  int
	WindowHeight int
	WindowTitle  string
}

func (a *BaseApp) Update() {
}

func (a *BaseApp) Render() {
}

func (a *BaseApp) Run() (err error) {
	var (
		context *core.Context
	)
	if context, err = core.NewContext(); err != nil {
		return
	}
	defer context.Delete()
	if err = context.CreateWindow(
		a.WindowWidth,
		a.WindowHeight,
		a.WindowTitle,
	); err != nil {
		return
	}
	for !context.ShouldClose() {
		context.Events.Poll()
		context.Clear()
	}
	glog.Flush()
	return
}
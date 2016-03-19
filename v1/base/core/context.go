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

package core

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Context struct {
	window        *glfw.Window
	OpenGLVersion string
	ShaderVersion string
	cursor        bool
	fullscreen    bool
	w             int
	h             int
	name          string
	initialized   bool
	Events        *Events
}

func NewContext() (context *Context, err error) {
	if err = glfw.Init(); err != nil {
		return
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.StencilBits, 8)
	glfw.WindowHint(glfw.DepthBits, 24)
	context = &Context{
		cursor:     true,
		fullscreen: false,
	}
	return
}

func (c *Context) Camera(worldCenter, worldSize mgl32.Vec3) (*Camera, error) {
	return NewCamera(
		worldCenter,
		worldSize,
		mgl32.Vec2{float32(c.w), float32(c.h)},
	)
}

func (c *Context) SetResizable(val bool) {
	if val {
		glfw.WindowHint(glfw.Resizable, 1)
	} else {
		glfw.WindowHint(glfw.Resizable, 0)
	}
}

func (c *Context) SetCursor(val bool) {
	c.cursor = val
	if c.window != nil {
		if c.cursor {
			c.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		} else {
			c.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		}
	}
}

func (c *Context) SetFullscreen(val bool) error {
	if c.fullscreen == val {
		return nil
	}
	c.fullscreen = val
	return c.createWindow()
}

func (c *Context) SetSwapInterval(val int) {
	glfw.SwapInterval(val)
}

func (c *Context) createWindow() (err error) {
	var (
		monitor *glfw.Monitor = nil
	)
	if c.window != nil {
		win := c.window
		c.window = nil
		win.Destroy()
	}
	if c.fullscreen == true {
		monitor = glfw.GetPrimaryMonitor()
	}
	if c.window, err = glfw.CreateWindow(c.w, c.h, c.name, monitor, nil); err != nil {
		return
	}
	c.Events = newEvents(c.window)
	c.window.MakeContextCurrent()
	return
}

func (c *Context) CreateWindow(w, h int, name string) (err error) {
	c.w = w
	c.h = h
	c.name = name
	c.createWindow()
	c.SetCursor(c.cursor)
	gl.Init()
	if e := gl.GetError(); e != 0 {
		if e != gl.INVALID_ENUM {
			err = fmt.Errorf("OpenGL glInit error: %X\n", e)
			return
		}
	}
	c.OpenGLVersion = glfw.GetVersionString()
	c.ShaderVersion = gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Disable(gl.CULL_FACE)
	glfw.SwapInterval(1)
	return
}

func (c *Context) ShouldClose() bool {
	return c.window.ShouldClose()
}

func (c *Context) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (c *Context) SwapBuffers() {
	c.window.SwapBuffers()
}

func (c *Context) Delete() {
	if c.window != nil {
		c.window.Destroy()
	}
	glfw.Terminate()
}

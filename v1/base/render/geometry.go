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

package render

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pikkpoiss/gamejam/v1/base/core"
	"unsafe"
)

type Point struct {
	Position mgl32.Vec3
	Texture  mgl32.Vec2
	Frame    float32
}

var Square = []Point{
	Point{
		Position: mgl32.Vec3{-0.5, -0.5, 0},
		Texture:  mgl32.Vec2{0, 0},
		Frame:    0,
	},
	Point{
		Position: mgl32.Vec3{0.5, 0.5, 0},
		Texture:  mgl32.Vec2{1, 1},
		Frame:    0,
	},
	Point{
		Position: mgl32.Vec3{-0.5, 0.5, 0},
		Texture:  mgl32.Vec2{0, 1},
		Frame:    0,
	},
	Point{
		Position: mgl32.Vec3{-0.5, -0.5, 0},
		Texture:  mgl32.Vec2{0, 0},
		Frame:    0,
	},
	Point{
		Position: mgl32.Vec3{0.5, -0.5, 0},
		Texture:  mgl32.Vec2{1, 0},
		Frame:    0,
	},
	Point{
		Position: mgl32.Vec3{0.5, 0.5, 0},
		Texture:  mgl32.Vec2{1, 1},
		Frame:    0,
	},
}

type Geometry struct {
	Points []Point
	Dirty  bool
	vbo    *core.ArrayBuffer
	stride uintptr
}

func NewGeometry(capacity int) (out *Geometry) {
	var (
		point  Point
		stride uintptr = unsafe.Sizeof(point)
	)
	out = &Geometry{
		Points: make([]Point, 0, capacity),
		Dirty:  true,
		stride: stride,
		vbo:    core.NewArrayBuffer(),
	}
	return
}

func NewGeometryFromPoints(points []Point) (out *Geometry) {
	out = NewGeometry(len(points))
	out.Points = append(out.Points, points...)
	return
}

func (g *Geometry) Bind() {
	g.vbo.Bind()
}

func (g *Geometry) Delete() {
	if g.vbo != nil {
		g.vbo.Delete()
		g.vbo = nil
	}
}

func (g *Geometry) Upload() {
	if g.Dirty {
		g.vbo.Upload(g.Points, len(g.Points)*int(g.stride))
		g.Dirty = false
	}
}

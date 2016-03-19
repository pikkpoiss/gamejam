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

package loaders

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pikkpoiss/gamejam/v1/base/render"
	"github.com/pikkpoiss/gamejam/v1/base/sprites"
	"strings"
)

type TextMapping struct {
	defaultSprite int
	mapping       map[rune]int
	sheet         *sprites.Sheet
}

func NewTextMapping(sheet *sprites.Sheet, defaultSprite string) (out *TextMapping, err error) {
	var (
		sprite *sprites.Sprite
	)
	out = &TextMapping{
		sheet:   sheet,
		mapping: map[rune]int{},
	}
	if sprite, err = sheet.Sprite(defaultSprite); err != nil {
		return
	}
	out.defaultSprite = sprite.Index()
	return

}

func (m *TextMapping) Set(r rune, key string) (err error) {
	var (
		sprite *sprites.Sprite
	)
	if sprite, err = m.sheet.Sprite(key); err != nil {
		return
	}
	m.mapping[r] = sprite.Index()
	return
}

func (m *TextMapping) Get(r rune) (index int) {
	var exists bool
	if index, exists = m.mapping[r]; !exists {
		index = m.defaultSprite
	}
	return
}

type TextLoader struct {
}

func NewTextLoader() *TextLoader {
	return &TextLoader{}
}

func (l *TextLoader) add(x, y float32, index int, scale float32, geo *render.Geometry) {
	var (
		findex = float32(index)
		unit   = scale
	)
	x = x * scale
	y = y * scale
	geo.Points = append(geo.Points, []render.Point{
		render.Point{
			Position: mgl32.Vec3{x, y, 0},
			Texture:  mgl32.Vec2{0, 0},
			Frame:    findex,
		},
		render.Point{
			Position: mgl32.Vec3{x + unit, y + unit, 0},
			Texture:  mgl32.Vec2{1, 1},
			Frame:    findex,
		},
		render.Point{
			Position: mgl32.Vec3{x, y + unit, 0},
			Texture:  mgl32.Vec2{0, 1},
			Frame:    findex,
		},
		render.Point{
			Position: mgl32.Vec3{x, y, 0},
			Texture:  mgl32.Vec2{0, 0},
			Frame:    findex,
		},
		render.Point{
			Position: mgl32.Vec3{x + unit, y, 0},
			Texture:  mgl32.Vec2{1, 0},
			Frame:    findex,
		},
		render.Point{
			Position: mgl32.Vec3{x + unit, y + unit, 0},
			Texture:  mgl32.Vec2{1, 1},
			Frame:    findex,
		},
	}...)
}

func (l *TextLoader) Load(mapping *TextMapping, scale float32, grid string) (geometry *render.Geometry, err error) {
	var (
		lines []string
		line  string
		char  rune
		y     int
		x     int
	)
	lines = strings.Split(strings.TrimSpace(grid), "\n")
	if len(lines) == 0 {
		err = fmt.Errorf("No lines in input data")
		return
	}
	geometry = render.NewGeometry(len(lines) * len(lines[0]) * 6)
	for y, line = range lines {
		for x, char = range line {
			l.add(float32(x), float32(y), mapping.Get(char), scale, geometry)
		}
	}
	return
}

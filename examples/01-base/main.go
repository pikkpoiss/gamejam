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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/glog"
	"github.com/pikkpoiss/gamejam/v1/base/core"
	"github.com/pikkpoiss/gamejam/v1/base/loaders"
	"github.com/pikkpoiss/gamejam/v1/base/render"
	"github.com/pikkpoiss/gamejam/v1/base/sprites"
	"github.com/pikkpoiss/gamejam/v1/base/text"
	"github.com/pikkpoiss/gamejam/v1/base/util"
	"image/color"
	"runtime"
)

const BATCH = `
AAA
BBB
`

type Inst struct {
	Key string
	X   float32
	Y   float32
	R   float32
}

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

func main() {
	flag.Parse()

	const (
		WinTitle              = "01-base"
		WinWidth              = 640
		WinHeight             = 480
		PixelsPerUnit float32 = 100
	)

	var (
		context         *core.Context
		sheet           *sprites.Sheet
		camera          *core.Camera
		framerate       *util.Framerate
		font            *text.FontFace
		fg              = color.RGBA{255, 255, 255, 255}
		bg              = color.RGBA{0, 0, 0, 255}
		err             error
		inst            *render.Instance
		rot             int = 0
		textMapping     *loaders.TextMapping
		batchData       *render.Geometry
		renderer        *render.Renderer
		spriteInstances *sprites.SpriteInstanceList
		textInstances   *text.TextInstanceList
		batchInstances  *render.InstanceList
		square          *render.Geometry
	)
	if context, err = core.NewContext(); err != nil {
		panic(err)
	}
	if err = context.CreateWindow(WinWidth, WinHeight, WinTitle); err != nil {
		panic(err)
	}
	if renderer, err = render.NewRenderer(100); err != nil {
		panic(err)
	}

	if sheet, err = loaders.NewTexturePackerLoader().Load(
		"examples/resources/spritesheet.json",
		core.SmoothingNearest,
	); err != nil {
		panic(err)
	}
	if textMapping, err = loaders.NewTextMapping(
		sheet,
		"numbered_squares_03",
	); err != nil {
		panic(err)
	}
	textMapping.Set('A', "numbered_squares_01")
	textMapping.Set('B', "numbered_squares_tall_16")
	if batchData, err = loaders.NewTextLoader().Load(
		textMapping,
		1,
		BATCH,
	); err != nil {
		panic(err)
	}

	spriteInstances = sprites.NewSpriteInstanceList(sheet, PixelsPerUnit)
	batchInstances = render.NewInstanceList()

	batchInstances.NewInstance()

	square = render.NewGeometryFromPoints(render.Square)

	textInstances = text.NewTextInstanceList(text.Config{
		TextureWidth:  512,
		TextureHeight: 512,
		PixelsPerUnit: PixelsPerUnit,
	})
	if framerate, err = util.NewFramerateRenderer(); err != nil {
		panic(err)
	}
	if camera, err = context.Camera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{6.4, 4.8, 2}); err != nil {
		panic(err)
	}
	if font, err = text.NewFontFace("examples/resources/Roboto-Light.ttf", 24, fg, bg); err != nil {
		panic(err)
	}
	for _, s := range []Inst{
		Inst{Key: "This is text!", X: 0, Y: -1.0, R: 0},
		Inst{Key: "More text!", X: 1.0, Y: 1.0, R: 15},
	} {
		inst = textInstances.NewInstance()
		if err = textInstances.SetText(inst, s.Key, font); err != nil {
			panic(err)
		}
		inst.SetPosition(mgl32.Vec3{s.X, s.Y, 0})
		inst.SetRotation(s.R)
	}
	for _, s := range []Inst{
		Inst{Key: "numbered_squares_01", X: 0, Y: 0, R: 0},
		Inst{Key: "numbered_squares_02", X: -1.5, Y: -1.5, R: -15},
		Inst{Key: "numbered_squares_03", X: -2.0, Y: -2.0, R: -30},
	} {
		inst = spriteInstances.NewInstance()
		if err = spriteInstances.SetFrame(inst, s.Key); err != nil {
			panic(err)
		}
		inst.SetPosition(mgl32.Vec3{s.X, s.Y, 0})
		inst.SetRotation(s.R)
	}

	for !context.ShouldClose() {
		context.Events.Poll()
		context.Clear()

		renderer.Bind()
		sheet.Bind()

		renderer.Render(camera, sheet, batchData, batchInstances)
		renderer.Render(camera, sheet, square, spriteInstances)

		textInstances.Bind()
		renderer.Render(camera, textInstances.Sheet(), square, textInstances)
		textInstances.Unbind()

		renderer.Unbind()

		framerate.Bind()
		framerate.Render(camera)
		framerate.Unbind()

		context.SwapBuffers()

		if err = textInstances.SetText(
			textInstances.Head(),
			fmt.Sprintf("Rotation %v", rot%100),
			font,
		); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			break
		}
		inst.SetRotation(float32(rot))
		rot += 1
	}
	if err = core.WritePNG("test-packed.png", textInstances.Sheet().Image()); err != nil {
		panic(err)
	}
	textInstances.Delete()
	glog.Flush()
}

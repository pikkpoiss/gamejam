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
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	WorldCenter mgl32.Vec3
	WorldSize   mgl32.Vec3
	ScreenSize  mgl32.Vec2
	PxPerUnit   mgl32.Vec2
	Projection  mgl32.Mat4
	View        mgl32.Mat4
	Inverse     mgl32.Mat4
}

func NewCamera(worldCenter, worldSize mgl32.Vec3, screenSize mgl32.Vec2) (c *Camera, err error) {
	c = &Camera{
		View: mgl32.Ident4(),
	}
	c.SetScreenSize(screenSize)
	err = c.SetWorldBounds(worldCenter, worldSize)
	return
}

func (c *Camera) SetScreenSize(screenSize mgl32.Vec2) {
	c.ScreenSize = screenSize
	c.calcPxPerUnit()
}

func (c *Camera) SetWorldBounds(worldCenter, worldSize mgl32.Vec3) (err error) {
	c.WorldCenter = worldCenter
	c.WorldSize = worldSize
	c.calcPxPerUnit()
	var (
		defaultMat4 mgl32.Mat4
		half        = c.WorldSize.Mul(0.5)
		min         = c.WorldCenter.Sub(half)
		max         = c.WorldCenter.Add(half)
	)
	c.Projection = mgl32.Ortho(
		min.X(),
		max.X(),
		min.Y(),
		max.Y(),
		max.Z(),
		min.Z())
	c.Inverse = c.Projection.Inv()
	if c.Inverse == defaultMat4 {
		err = fmt.Errorf("Projection matrix not invertible")
	}
	return
}

func (c *Camera) calcPxPerUnit() {
	c.PxPerUnit = mgl32.Vec2{
		c.ScreenSize.X() / c.WorldSize.X(),
		c.ScreenSize.Y() / c.WorldSize.Y(),
	}
}

func (c *Camera) unproject(pt mgl32.Vec2) mgl32.Vec2 {
	var (
		screen = pt.Vec4(1, 1)
		out    mgl32.Vec4
	)
	out = c.Inverse.Mul4x1(screen)
	out = out.Mul(1.0 / out[3])
	return out.Vec2()
}

func (c *Camera) ScreenToWorldCoords(screenCoords mgl32.Vec2) mgl32.Vec2 {
	// http://stackoverflow.com/questions/7692988/
	var (
		half = c.ScreenSize.Mul(0.5)
		pt   = mgl32.Vec2{
			(screenCoords.X() - half.X()) / half.X(),
			(half.Y() - screenCoords.Y()) / half.Y(),
		}
	)
	return c.unproject(pt)
}

func (c *Camera) project(pt mgl32.Vec2) mgl32.Vec2 {
	var (
		screen = pt.Vec4(1, 1)
		out    mgl32.Vec4
	)
	out = c.Projection.Mul4x1(screen)
	return out.Vec2()
}

func (c *Camera) WorldToScreenCoords(pt mgl32.Vec2) mgl32.Vec2 {
	var (
		pct  = c.project(pt)
		half = c.ScreenSize.Mul(0.5)
		out  = mgl32.Vec2{
			pct.X()*half.X() + half.X(),
			pct.Y()*half.Y() + half.Y(),
		}
	)
	return out
}

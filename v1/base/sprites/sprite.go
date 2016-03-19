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

package sprites

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pikkpoiss/gamejam/v1/base/render"
	"image"
)

type Sprite struct {
	index  int
	bounds mgl32.Vec2
	offset mgl32.Vec2
}

func (s *Sprite) Index() int {
	return s.index
}

func (s *Sprite) ImageBounds() image.Rectangle {
	return image.Rectangle{
		image.Point{int(s.offset.X()), int(s.offset.Y())},
		image.Point{int(s.offset.X() + s.bounds.X()), int(s.offset.Y() + s.bounds.Y())},
	}
}

func (s *Sprite) textureBounds(textureBounds mgl32.Vec2) render.UniformSprite {
	return render.NewUniformSprite(
		s.bounds.X()/textureBounds.X(),
		s.bounds.Y()/textureBounds.Y(),
		s.offset.X()/textureBounds.X(),
		1.0-(s.offset.Y()+s.bounds.Y()-1.0)/textureBounds.Y(),
	)
}

func (s *Sprite) WorldDimensions(pxPerUnit float32) mgl32.Vec2 {
	return mgl32.Vec2{
		s.bounds.X() / pxPerUnit,
		s.bounds.Y() / pxPerUnit,
	}
}

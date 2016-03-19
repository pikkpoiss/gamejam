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
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pikkpoiss/gamejam/v1/base/core"
	"github.com/pikkpoiss/gamejam/v1/base/render"
	"unsafe"
)

type Sheet struct {
	keys            map[string]*Sprite
	texture         *core.Texture
	ubo             *core.UniformBuffer
	Count           int
	version         int
	uploadedVersion int
}

func NewSheet() *Sheet {
	return &Sheet{
		keys:            map[string]*Sprite{},
		version:         0,
		uploadedVersion: -1,
		ubo:             core.NewUniformBuffer(),
	}
}

func (s *Sheet) SetTexture(texture *core.Texture) {
	s.deleteTexture()
	s.texture = texture
}

func (s *Sheet) Bind() {
	if s.texture != nil {
		s.texture.Bind()
	}
	s.upload()
}

func (s *Sheet) Unbind() {
	if s.texture != nil {
		s.texture.Unbind()
	}
}

func (s *Sheet) deleteTexture() {
	if s.texture != nil {
		s.texture.Delete()
		s.texture = nil
	}
}

func (s *Sheet) Delete() {
	s.deleteTexture()
	if s.ubo != nil {
		s.ubo.Delete()
		s.ubo = nil
	}
}

func (s *Sheet) AddSprite(key string, bounds, offset mgl32.Vec2) (out *Sprite) {
	var index int
	index = s.Count
	out = &Sprite{
		index:  index,
		bounds: bounds,
		offset: offset,
	}
	s.keys[key] = out
	s.Count++
	s.version++
	return
}

func (s *Sheet) Exists(key string) (exists bool) {
	_, exists = s.keys[key]
	return
}

func (s *Sheet) Sprite(key string) (out *Sprite, err error) {
	var exists bool
	if out, exists = s.keys[key]; !exists {
		err = fmt.Errorf("Invalid tile key %v", key)
		return
	}
	return
}

func (s *Sheet) upload() (err error) {
	if s.version == s.uploadedVersion {
		return
	}
	var (
		sprite *Sprite
		entry  render.UniformSprite
		data   = make([]render.UniformSprite, s.Count)
		size   = s.Count * int(unsafe.Sizeof(entry))
	)
	if s.texture == nil {
		err = fmt.Errorf("No texture associated with sheet")
		return
	}
	for _, sprite = range s.keys {
		data[sprite.index] = sprite.textureBounds(s.texture.Size)
	}
	s.ubo.Upload(data, size)
	s.uploadedVersion = s.version
	return
}

func (s *Sheet) Texture() *core.Texture {
	return s.texture
}

func (s *Sheet) BufferID() uint32 {
	return s.ubo.BufferID()
}

func (s *Sheet) Size() int {
	return s.ubo.Size()
}

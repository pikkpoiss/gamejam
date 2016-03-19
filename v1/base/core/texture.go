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
	"bytes"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image"
)

type TextureSmoothing int

const (
	SmoothingNearest TextureSmoothing = gl.NEAREST
	SmoothingLinear  TextureSmoothing = gl.LINEAR
)

type Texture struct {
	id           uint32
	Size         mgl32.Vec2
	OriginalSize mgl32.Vec2
}

func LoadTexture(path string, smoothing TextureSmoothing) (texture *Texture, err error) {
	var img image.Image
	if img, err = LoadPNG(path); err != nil {
		return
	}
	return GetTexture(img, smoothing)
}

func GetTexture(img image.Image, smoothing TextureSmoothing) (texture *Texture, err error) {
	var (
		originalBounds = img.Bounds()
		textureId      uint32
		resizedImg     = getPow2Image(img)
		resizedBounds  = resizedImg.Bounds()
	)
	if textureId, err = getGLTexture(img, smoothing); err != nil {
		return
	}
	texture = &Texture{
		id: textureId,
		Size: mgl32.Vec2{
			float32(resizedBounds.Dx()),
			float32(resizedBounds.Dy()),
		},
		OriginalSize: mgl32.Vec2{
			float32(originalBounds.Dx()),
			float32(originalBounds.Dy()),
		},
	}
	return
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t *Texture) Delete() {
	if t.id != 0 {
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.DeleteTextures(1, &t.id)
		t.id = 0
	}
}

func getGLTexture(img image.Image, smoothing TextureSmoothing) (t uint32, err error) {
	var (
		data   *bytes.Buffer
		bounds image.Rectangle
		width  int
		height int
	)
	if data, err = imageBytes(img); err != nil {
		return
	}
	bounds = img.Bounds()
	width = bounds.Max.X - bounds.Min.X
	height = bounds.Max.Y - bounds.Min.Y
	gl.GenTextures(1, &t)
	gl.BindTexture(gl.TEXTURE_2D, t)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(smoothing))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(smoothing))
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(width),
		int32(height),
		0,
		gl.RGBA,
		gl.UNSIGNED_INT_8_8_8_8,
		gl.Ptr(data.Bytes()),
	)
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return
}

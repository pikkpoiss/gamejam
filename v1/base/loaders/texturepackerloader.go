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
	"encoding/json"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pikkpoiss/gamejam/v1/base/core"
	"github.com/pikkpoiss/gamejam/v1/base/sprites"
	"io/ioutil"
	"path"
)

type texturePackerFloatCoords struct {
	X float32 `json:x,omitempty`
	Y float32 `json:y,omitempty`
}

type texturePackerIntCoords struct {
	X int `json:x,omitempty`
	Y int `json:y,omitempty`
	W int `json:w,omitempty`
	H int `json:h,omitempty`
}

type texturePackerFrame struct {
	Filename         string                   `json:filename`
	Frame            texturePackerIntCoords   `json:frame`
	Rotated          bool                     `json:rotated`
	Trimmed          bool                     `json:trimmed`
	SpriteSourceSize texturePackerIntCoords   `json:spriteSourceSize`
	SourceSize       texturePackerIntCoords   `json:sourceSize`
	Pivot            texturePackerFloatCoords `json:pivot`
}

type texturePackerMeta struct {
	Image  string                 `json:image`
	Format string                 `json:format`
	Size   texturePackerIntCoords `json:size`
	Scale  string                 `json:scale`
}

type texturePackerJSONArray struct {
	Frames []texturePackerFrame `json:frames`
	Meta   texturePackerMeta    `json:meta`
}

type TexturePackerLoader struct {
}

func NewTexturePackerLoader() *TexturePackerLoader {
	return &TexturePackerLoader{}
}

func (l *TexturePackerLoader) Load(jsonPath string, smoothing core.TextureSmoothing) (sheet *sprites.Sheet, err error) {
	var (
		dir         string
		data        []byte
		texturePath string
		parsed      texturePackerJSONArray
		texture     *core.Texture
	)
	dir = path.Dir(jsonPath)
	if data, err = ioutil.ReadFile(jsonPath); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(data), &parsed); err != nil {
		return
	}
	sheet = sprites.NewSheet()
	for _, frame := range parsed.Frames {
		sheet.AddSprite(
			frame.Filename,
			mgl32.Vec2{
				float32(frame.Frame.W),
				float32(frame.Frame.H),
			},
			mgl32.Vec2{
				float32(frame.Frame.X),
				float32(frame.Frame.Y),
			},
		)
	}
	texturePath = path.Join(dir, parsed.Meta.Image)
	if texture, err = core.LoadTexture(texturePath, smoothing); err != nil {
		return
	}
	sheet.SetTexture(texture)
	return
}

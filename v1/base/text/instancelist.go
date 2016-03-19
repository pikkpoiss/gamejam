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

package text

import (
	"github.com/golang/glog"
	"github.com/pikkpoiss/gamejam/v1/base/core"
	"github.com/pikkpoiss/gamejam/v1/base/render"
	"github.com/pikkpoiss/gamejam/v1/base/sprites"
	"image/draw"
)

type Config struct {
	TextureWidth  int
	TextureHeight int
	PixelsPerUnit float32
}

type TextInstanceList struct {
	*render.InstanceList
	cfg   Config
	sheet *sprites.PackedSheet
}

func NewTextInstanceList(cfg Config) *TextInstanceList {
	return &TextInstanceList{
		InstanceList: render.NewInstanceList(),
		cfg:          cfg,
		sheet: sprites.NewPackedSheet(
			cfg.TextureWidth,
			cfg.TextureHeight,
		),
	}
}

func (l *TextInstanceList) SetText(instance *render.Instance, text string, font *FontFace) (err error) {
	var (
		img    draw.Image
		sprite *sprites.Sprite
	)
	if instance == nil {
		return // No error.
	}
	if img, err = font.GetImage(text); err != nil {
		return
	}
	if err = l.sheet.Pack(text, img); err != nil {
		// Attempt to compact the texture.
		if err = l.repackImage(); err != nil {
			return
		}
		if err = l.sheet.Pack(text, img); err != nil {
			return
		}
	}
	if sprite, err = l.sheet.Sprite(text); err != nil {
		return
	}
	instance.Frame = sprite.Index()
	instance.SetScale(sprite.WorldDimensions(l.cfg.PixelsPerUnit).Vec3(1.0))
	instance.MarkChanged()
	instance.Key = text
	if err = l.generateTexture(); err != nil {
		return
	}
	return
}

func (l *TextInstanceList) generateTexture() (err error) {
	var (
		texture *core.Texture
	)
	if texture, err = core.GetTexture(
		l.sheet.Image(),
		core.SmoothingLinear,
	); err != nil {
		return
	}
	l.sheet.SetTexture(texture)
	return
}

func (l *TextInstanceList) repackImage() (err error) {
	var (
		newImage *sprites.PackedSheet
		instance *render.Instance
		sprite   *sprites.Sprite
	)
	if glog.V(1) {
		glog.Info("Repacking image")
	}
	newImage = sprites.NewPackedSheet(
		l.sheet.Width,
		l.sheet.Height,
	)
	instance = l.Head()
	for instance != nil {
		if err = newImage.Copy(instance.Key, l.sheet); err != nil {
			return
		}
		if sprite, err = newImage.Sheet.Sprite(instance.Key); err != nil {
			return
		}
		instance.Frame = sprite.Index()
		instance.MarkChanged()
		instance = instance.Next()
	}
	l.sheet = newImage
	if err = l.generateTexture(); err != nil {
		return
	}
	if glog.V(1) {
		glog.Info("Done repacking")
	}
	return
}

func (l *TextInstanceList) Bind() {
	l.sheet.Bind()
}

func (l *TextInstanceList) Unbind() {
	l.sheet.Unbind()
}

func (l *TextInstanceList) Delete() {
	l.sheet.Delete()
	l.sheet = nil
}

func (l *TextInstanceList) Sheet() *sprites.PackedSheet {
	return l.sheet
}

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
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pikkpoiss/gamejam/v1/base/core"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
)

type FontFace struct {
	font    *truetype.Font
	charw   float32
	charh   float32
	offy    float32
	fg      color.Color
	bg      color.Color
	context *freetype.Context
}

func NewFontFace(path string, pixels float32, fg, bg color.Color) (fontface *FontFace, err error) {
	var (
		font      *truetype.Font
		fontbytes []byte
		bounds    fixed.Rectangle26_6
		context   *freetype.Context
		points    float32
		dpi       float32 = 96
	)
	if fontbytes, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if font, err = freetype.ParseFont(fontbytes); err != nil {
		return
	}
	points = pixels * 72 / dpi
	bounds = font.Bounds(fixed.I(int(pixels)))
	context = freetype.NewContext()
	context.SetFont(font)
	context.SetFontSize(float64(points))
	context.SetDPI(float64(dpi))
	fontface = &FontFace{
		font:    font,
		charw:   float32(bounds.Max.X-bounds.Min.X) / 64,
		charh:   float32(bounds.Max.Y-bounds.Min.Y) / 64,
		offy:    float32(bounds.Min.Y) / 64,
		fg:      fg,
		bg:      bg,
		context: context,
	}
	return
}

func (ff *FontFace) GetImage(text string) (img draw.Image, err error) {
	var (
		src image.Image
		bg  image.Image
		dst draw.Image
		pt  fixed.Point26_6
		w   int
		h   int
	)
	src = image.NewUniform(ff.fg)
	bg = image.NewUniform(ff.bg)
	w = int(float32(len(text)) * ff.charw)
	h = int(ff.charh)
	dst = image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(dst, dst.Bounds(), bg, image.ZP, draw.Src)
	ff.context.SetSrc(src)
	ff.context.SetDst(dst)
	ff.context.SetClip(dst.Bounds())
	pt = freetype.Pt(0, int(ff.charh+ff.offy))
	if pt, err = ff.context.DrawString(text, pt); err != nil {
		return
	}
	img = image.NewRGBA(image.Rect(0, 0, int(pt.X/64), int(pt.Y/64)))
	draw.Draw(img, img.Bounds(), dst, image.Pt(0, -int(ff.offy)), draw.Src)
	return
}

func (ff *FontFace) GetText(text string) (t *core.Texture, err error) {
	var (
		img image.Image
	)
	if img, err = ff.GetImage(text); err != nil {
		return
	}
	t, err = core.GetTexture(img, gl.NEAREST)
	return
}

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
	"github.com/golang/glog"
	"image"
	"image/draw"
)

type PackedSheet struct {
	*Sheet
	Width   int
	Height  int
	img     draw.Image
	shelves []*packingShelf
}

func NewPackedSheet(w, h int) (i *PackedSheet) {
	return &PackedSheet{
		Width:   w,
		Height:  h,
		img:     image.NewRGBA(image.Rect(0, 0, w, h)),
		shelves: []*packingShelf{newShelf()},
		Sheet:   NewSheet(),
	}
}

func (s *PackedSheet) Image() image.Image {
	return s.img
}

func (s *PackedSheet) Pack(key string, img image.Image) (err error) {
	return s.packSprite(key, img, img.Bounds())
}

func (s *PackedSheet) Copy(key string, src *PackedSheet) (err error) {
	var (
		bounds image.Rectangle
		sprite *Sprite
	)
	if sprite, err = src.Sprite(key); err != nil {
		return
	}
	bounds = sprite.ImageBounds()
	return s.packSprite(key, src.img, bounds)
}

func (s *PackedSheet) packSprite(key string, src image.Image, srcBounds image.Rectangle) (err error) {
	var (
		j         int
		shelf     *packingShelf
		score     int
		bestScore int             = -1
		bestShelf int             = -1
		texBounds image.Rectangle = s.img.Bounds()
		w         int             = srcBounds.Max.X - srcBounds.Min.X
		h         int             = srcBounds.Max.Y - srcBounds.Min.Y
		maxW      int             = texBounds.Max.X
	)
	if s.Exists(key) {
		// Don't need to pack since it's already in here
		return
	}
	for j, shelf = range s.shelves {
		if shelf.CanAdd(w, h, maxW) {
			score = shelf.BestAreaFit(w, h, maxW)
			if score > bestScore {
				bestScore = score
				bestShelf = j
			}
		}
	}
	if bestShelf == -1 {
		shelf = s.shelves[len(s.shelves)-1]
		if shelf.y+shelf.height+h > texBounds.Max.Y {
			// New packingShelf would exceed current image size
			err = fmt.Errorf("Cannot fit text into texture")
			return
		}
		s.shelves = append(s.shelves, s.shelves[len(s.shelves)-1].Close())
		bestShelf = len(s.shelves) - 1
	}
	shelf = s.shelves[bestShelf]
	var (
		x, y     = shelf.Add(w, h)
		destPt   = image.Pt(x, y)
		destRect = image.Rectangle{destPt, destPt.Add(image.Pt(w, h))}
	)
	s.AddSprite(key, mgl32.Vec2{float32(w), float32(h)}, mgl32.Vec2{float32(x), float32(y)})
	if glog.V(2) {
		glog.Infof("packRegion(%v): dest %v src %v", key, destRect, srcBounds.Min)
	}
	draw.Draw(s.img, destRect, src, srcBounds.Min, draw.Src)
	return
}

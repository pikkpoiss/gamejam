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

type packingShelf struct {
	x      int
	y      int
	height int
	isOpen bool
}

func newShelf() *packingShelf {
	return &packingShelf{
		x:      0,
		y:      0,
		height: 0,
		isOpen: true,
	}
}

func (s *packingShelf) FitsX(w, maxW int) bool {
	return s.x+w <= maxW
}

func (s *packingShelf) FitsY(h int) bool {
	return s.height >= h
}

func (s *packingShelf) RemainingX(maxW int) int {
	return maxW - s.x
}

func (s *packingShelf) CanAdd(w, h, maxW int) bool {
	if !s.FitsX(w, maxW) {
		return false
	}
	if !s.isOpen && !s.FitsY(h) {
		return false
	}
	return true
}

func (s *packingShelf) Add(w, h int) (origX, origY int) {
	origX = s.x
	origY = s.y
	if s.height < h {
		s.height = h
	}
	s.x += w
	return
}

func (s *packingShelf) Close() (out *packingShelf) {
	out = newShelf()
	out.y = s.y + s.height
	s.isOpen = false
	return
}

func (s *packingShelf) BestAreaFit(w, h, maxW int) int {
	var (
		packingShelfArea = s.RemainingX(maxW) * s.height
		wordArea  = w * h
	)
	return packingShelfArea - wordArea
}

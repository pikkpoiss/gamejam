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
	"github.com/pikkpoiss/gamejam/v1/base/render"
)

type SpriteInstances interface {
	render.Instances
	SetFrame(instance *render.Instance, frame string) (err error)
}

type SpriteInstanceList struct {
	*render.InstanceList
	pixelsPerUnit float32 // TODO: Should be in sheet
	sheet         *Sheet
}

func NewSpriteInstanceList(sheet *Sheet, pixelsPerUnit float32) *SpriteInstanceList {
	return &SpriteInstanceList{
		InstanceList:  render.NewInstanceList(),
		pixelsPerUnit: pixelsPerUnit,
		sheet:         sheet,
	}
}

func (l *SpriteInstanceList) SetFrame(instance *render.Instance, frame string) (err error) {
	var s *Sprite
	if instance == nil {
		return // No error
	}
	if s, err = l.sheet.Sprite(frame); err != nil {
		return
	}
	instance.Frame = s.Index()
	instance.SetScale(s.WorldDimensions(l.pixelsPerUnit).Vec3(1.0))
	instance.MarkChanged()
	instance.Key = frame
	return
}

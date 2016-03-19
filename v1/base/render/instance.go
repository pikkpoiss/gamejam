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

package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Instance struct {
	model    mgl32.Mat4
	position mgl32.Vec3
	scale    mgl32.Vec3
	rotation float32
	Frame    int
	Key      string // TODO: move to an interface{} data pointer.
	color    mgl32.Vec4
	dirty    bool
	next     *Instance
	prev     *Instance
	list     Instances
}

func newInstance() *Instance {
	return &Instance{
		scale:    mgl32.Vec3{1.0, 1.0, 1.0},
		position: mgl32.Vec3{0.0, 0.0, 0.0},
		color:    mgl32.Vec4{0.0, 0.0, 0.0, 0.0},
		rotation: 0,
		dirty:    true,
	}
}

func (i *Instance) SetScale(s mgl32.Vec3) {
	if i.scale.X() != s.X() || i.scale.Y() != s.Y() || i.scale.Z() != s.Z() {
		i.scale = s
		i.dirty = true
	}
}

func (i *Instance) SetPosition(p mgl32.Vec3) {
	if i.position.X() != p.X() || i.position.Y() != p.Y() || i.position.Z() != p.Z() {
		i.position = p
		i.dirty = true
	}
}

func (i *Instance) SetRotation(r float32) {
	if i.rotation != r {
		i.rotation = r
		i.dirty = true
	}
}

func (i *Instance) GetModel() mgl32.Mat4 {
	if i.dirty {
		var model mgl32.Mat4
		model = mgl32.Translate3D(
			i.position.X(),
			i.position.Y(),
			i.position.Z(),
		)
		model = model.Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(i.rotation)))
		model = model.Mul4(mgl32.Scale3D(i.scale.X(), i.scale.Y(), i.scale.Z()))
		i.model = model
		i.dirty = false
	}
	return i.model
}

func (i *Instance) Color() mgl32.Vec4 {
	return i.color
}

func (i *Instance) SetColor(r, g, b, a float32) {
	i.color[0] = r
	i.color[1] = g
	i.color[2] = b
	i.color[3] = a
	i.dirty = true
}

func (i *Instance) Next() *Instance {
	return i.next
}

func (i *Instance) Remove() {
	if i.next != nil {
		i.next.prev = i.prev
	}
	if i.prev != nil {
		i.prev.next = i.next
	}
	i.next = nil
	i.prev = nil
	i.list = nil
}

func (i *Instance) InsertAfter(inst *Instance) {
	if inst == nil {
		return
	}
	if i.next != nil {
		i.next.prev = inst
	}
	inst.next = i.next
	inst.prev = i
	inst.list = i.list
	i.next = inst
}

func (i *Instance) MarkChanged() {
	i.dirty = true
}

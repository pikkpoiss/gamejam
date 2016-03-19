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

type Instances interface {
	Head() *Instance
	Prepend(inst *Instance)
	NewInstance() (inst *Instance)
}

type InstanceList struct {
	count int
	root  Instance
}

func NewInstanceList() (l *InstanceList) {
	l = &InstanceList{}
	l.root.list = l
	return
}

func (l *InstanceList) Head() *Instance {
	return l.root.next
}

func (l *InstanceList) Prepend(inst *Instance) {
	l.root.InsertAfter(inst)
	l.count++
}

func (l *InstanceList) NewInstance() (inst *Instance) {
	inst = newInstance()
	inst.SetPosition(mgl32.Vec3{0, 0, 0})
	inst.SetScale(mgl32.Vec3{1.0, 1.0, 1.0})
	inst.SetRotation(0)
	inst.Frame = 0
	l.Prepend(inst)
	return
}

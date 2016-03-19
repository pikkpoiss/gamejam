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
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Buffer interface {
}

type GLBuffer struct {
	id          uint32
	bufferBytes int
	target      uint32
}

func NewGLBuffer(target uint32) (b *GLBuffer) {
	b = &GLBuffer{
		target: target,
	}
	gl.GenBuffers(1, &b.id)
	b.Bind()
	return
}

func (b *GLBuffer) BufferID() uint32 {
	return b.id
}

func (b *GLBuffer) Bind() {
	gl.BindBuffer(b.target, b.id)
}

func (b *GLBuffer) Delete() {
	gl.DeleteBuffers(1, &b.id)
}

func (b *GLBuffer) Upload(data interface{}, size int) {
	b.Bind()
	if size > b.bufferBytes {
		b.bufferBytes = size
		gl.BufferData(b.target, size, gl.Ptr(data), gl.STREAM_DRAW)
	} else {
		gl.BufferSubData(b.target, 0, size, gl.Ptr(data))
	}
}

func (b *GLBuffer) Size() int {
	return b.bufferBytes
}

type UniformBuffer struct {
	*GLBuffer
}

func NewUniformBuffer() (b *UniformBuffer) {
	b = &UniformBuffer{
		GLBuffer: NewGLBuffer(gl.UNIFORM_BUFFER),
	}
	return
}

type ArrayBuffer struct {
	*GLBuffer
}

func NewArrayBuffer() (b *ArrayBuffer) {
	b = &ArrayBuffer{
		GLBuffer: NewGLBuffer(gl.ARRAY_BUFFER),
	}
	return
}

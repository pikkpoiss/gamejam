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

package util

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pikkpoiss/gamejam/v1/base/core"
	"time"
	"unsafe"
)

const FRAMERATE_FRAGMENT = `#version 150
precision mediump float;
uniform vec4 v_Color;
out vec4 v_FragData;
void main() {
  v_FragData = v_Color;
}` + "\x00"

const FRAMERATE_VERTEX = `#version 150
in vec2 v_Position;
uniform mat4 m_ModelView;
uniform mat4 m_Projection;
void main() {
  gl_Position = m_Projection * m_ModelView * vec4(v_Position, 0.0, 1.0);
}` + "\x00"

type framerateDataPoint struct {
	pos mgl32.Vec2
}

type Framerate struct {
	shader        *core.Program
	vbo           uint32
	vboBytes      int
	stride        int32
	locColor      int32
	locModelView  int32
	locProjection int32
	data          *framerateData
}

func NewFramerateRenderer() (r *Framerate, err error) {
	r = &Framerate{
		shader: core.NewProgram(),
		data:   newFramerateData(120),
	}
	if err = r.shader.Load(FRAMERATE_VERTEX, FRAMERATE_FRAGMENT); err != nil {
		return
	}
	r.shader.Bind()
	gl.GenBuffers(1, &r.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	var (
		point       framerateDataPoint
		locPosition = uint32(gl.GetAttribLocation(r.shader.ID(), gl.Str("v_Position\x00")))
		offPosition = gl.PtrOffset(int(unsafe.Offsetof(point.pos)))
	)
	r.stride = int32(unsafe.Sizeof(point))
	r.locColor = gl.GetUniformLocation(r.shader.ID(), gl.Str("v_Color\x00"))
	r.locModelView = gl.GetUniformLocation(r.shader.ID(), gl.Str("m_ModelView\x00"))
	r.locProjection = gl.GetUniformLocation(r.shader.ID(), gl.Str("m_Projection\x00"))
	gl.EnableVertexAttribArray(locPosition)
	gl.VertexAttribPointer(locPosition, 2, gl.FLOAT, false, r.stride, offPosition)
	return
}

func (r *Framerate) Bind() {
	r.shader.Bind()
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
}

func (r *Framerate) Unbind() {
	r.shader.Unbind()
}

func (r *Framerate) Delete() {
	r.shader.Delete()
}

func (r *Framerate) Render(camera *core.Camera) (err error) {
	r.data.Sample()
	var (
		modelView     = mgl32.Ident4()
		dataBytes int = int(r.data.Count) * int(r.stride)
	)
	gl.Uniform4f(r.locColor, 255.0/255.0, 0, 0, 255.0/255.0)
	gl.UniformMatrix4fv(r.locModelView, 1, false, &modelView[0])
	gl.UniformMatrix4fv(r.locProjection, 1, false, &camera.Projection[0])

	if dataBytes > r.vboBytes {
		r.vboBytes = dataBytes
		gl.BufferData(gl.ARRAY_BUFFER, dataBytes, gl.Ptr(r.data.Points), gl.STREAM_DRAW)
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, dataBytes, gl.Ptr(r.data.Points))
	}

	gl.DrawArrays(gl.POINTS, 0, r.data.Count)
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

type framerateData struct {
	Count  int32
	times  []time.Duration
	index  int32
	last   time.Time
	max    time.Duration
	Points []framerateDataPoint
}

func newFramerateData(size int32) *framerateData {
	return &framerateData{
		Count:  size,
		times:  make([]time.Duration, size),
		index:  0,
		Points: make([]framerateDataPoint, size),
	}
}

func (d *framerateData) Sample() {
	var (
		now  = time.Now()
		diff = now.Sub(d.last)
		i    int32
	)
	d.times[d.index] = diff
	d.index = (d.index + 1) % d.Count
	if diff > d.max {
		d.max = diff
	}
	for i = 0; i < d.Count; i++ {
		d.Points[i].pos = mgl32.Vec2{
			float32(i) / float32(d.Count),
			float32(d.times[(i+d.index)%d.Count]) / float32(d.max),
		}
	}
}

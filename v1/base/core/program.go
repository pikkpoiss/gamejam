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
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"strings"
	"unsafe"
)

type Uniform struct {
	location int32
}

func (u *Uniform) Mat4(m mgl32.Mat4) {
	gl.UniformMatrix4fv(u.location, 1, false, &m[0])
}

type Program struct {
	vao     uint32
	program uint32
}

func NewProgram() *Program {
	return &Program{}
}

func (p *Program) Delete() {
	gl.DeleteVertexArrays(1, &p.vao)
	p.vao = 0
	gl.DeleteProgram(p.program)
	p.program = 0
}

func (p *Program) Bind() {
	gl.BindVertexArray(p.vao)
	gl.UseProgram(p.program)
}

func (p *Program) Unbind() {
	gl.BindVertexArray(0)
}

func (p *Program) ID() uint32 {
	return p.program
}

func (p *Program) Load(vertex, fragment string) (err error) {
	if err = p.createVAO(); err != nil {
		return
	}
	if err = p.buildProgram(vertex, fragment); err != nil {
		return
	}
	return
}

func (p *Program) Uniform(name string) *Uniform {
	var nameStr = gl.Str(fmt.Sprintf("%v\x00", name))
	return &Uniform{
		location: gl.GetUniformLocation(p.ID(), nameStr),
	}
}

func (p *Program) UniformBlock(name string, binding uint32) *UniformBlock {
	var (
		nameStr = gl.Str(fmt.Sprintf("%v\x00", name))
		index   = uint32(gl.GetUniformBlockIndex(p.program, nameStr))
	)
	gl.UniformBlockBinding(p.program, index, binding)
	return &UniformBlock{
		binding: binding,
	}
}

func (p *Program) Attrib(name string, stride uintptr) *VertexAttribute {
	var nameStr = gl.Str(fmt.Sprintf("%v\x00", name))
	return &VertexAttribute{
		location: uint32(gl.GetAttribLocation(p.program, nameStr)),
		stride:   stride,
	}
}

func (p *Program) createVAO() error {
	gl.GenVertexArrays(1, &p.vao)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR gl.GenVertexArray %X", e)
	}
	gl.BindVertexArray(p.vao)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR array.Bind %X", e)
	}
	return nil
}

func (p *Program) compileShader(stype uint32, source string) (shader uint32, err error) {
	csource := gl.Str(source)
	shader = gl.CreateShader(stype)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)
	var status int32
	if gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status); status == gl.FALSE {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(shader, length, nil, gl.Str(log))
		err = fmt.Errorf("ERROR shader compile:\n%s", log)
	}
	return
}

func (p *Program) linkProgram(vertex uint32, fragment uint32) (program uint32, err error) {
	program = gl.CreateProgram()
	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.BindFragDataLocation(program, 0, gl.Str("v_FragData\x00"))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR program.BindFragDataLocation %X", e)
		return
	}
	gl.LinkProgram(program)
	var status int32
	if gl.GetProgramiv(program, gl.LINK_STATUS, &status); status == gl.FALSE {
		var length int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(program, length, nil, gl.Str(log))
		err = fmt.Errorf("ERROR program link:\n%s", log)
	}
	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)
	return
}

func (p *Program) buildProgram(vsrc string, fsrc string) (err error) {
	var (
		vertex   uint32
		fragment uint32
		glsrc    string
	)
	glsrc = fmt.Sprintf("%v\x00", vsrc)
	if vertex, err = p.compileShader(gl.VERTEX_SHADER, glsrc); err != nil {
		return
	}
	glsrc = fmt.Sprintf("%v\x00", fsrc)
	if fragment, err = p.compileShader(gl.FRAGMENT_SHADER, glsrc); err != nil {
		return
	}
	p.program, err = p.linkProgram(vertex, fragment)
	return
}

type UniformBlock struct {
	binding uint32
}

func (b *UniformBlock) Bind(bufferID uint32, size int) {
	gl.BindBufferRange(gl.UNIFORM_BUFFER, b.binding, bufferID, 0, size)
}

type VertexAttribute struct {
	location uint32
	stride   uintptr
}

func (a *VertexAttribute) vertexAttrib(l uint32, size int32, xtype uint32, offset uintptr, divisor uint32) {
	var offsetPtr = gl.PtrOffset(int(offset))
	gl.EnableVertexAttribArray(a.location + l)
	gl.VertexAttribPointer(a.location+l, size, xtype, false, int32(a.stride), offsetPtr)
	gl.VertexAttribDivisor(a.location+l, divisor)
}

func (a *VertexAttribute) Float(offset uintptr, divisor uint32) {
	a.vertexAttrib(0, 1, gl.FLOAT, offset, divisor)
}

func (a *VertexAttribute) Vec2(offset uintptr, divisor uint32) {
	a.vertexAttrib(0, 2, gl.FLOAT, offset, divisor)
}

func (a *VertexAttribute) Vec3(offset uintptr, divisor uint32) {
	a.vertexAttrib(0, 3, gl.FLOAT, offset, divisor)
}

func (a *VertexAttribute) Vec4(offset uintptr, divisor uint32) {
	a.vertexAttrib(0, 4, gl.FLOAT, offset, divisor)
}

func (a *VertexAttribute) Mat4(offset uintptr, divisor uint32) {
	var (
		float    float32
		sizeVec4 = unsafe.Sizeof(float) * 4
	)
	a.vertexAttrib(0, 4, gl.FLOAT, offset, divisor)
	a.vertexAttrib(1, 4, gl.FLOAT, offset+sizeVec4, divisor)
	a.vertexAttrib(2, 4, gl.FLOAT, offset+2*sizeVec4, divisor)
	a.vertexAttrib(3, 4, gl.FLOAT, offset+3*sizeVec4, divisor)
}

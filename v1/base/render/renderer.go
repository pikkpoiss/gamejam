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
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pikkpoiss/gamejam/v1/base/core"
	"unsafe"
)

const FRAGMENT = `#version 150

precision mediump float;

in vec2 v_TexturePos;
in vec2 v_TextureMin;
in vec2 v_TextureDim;
in vec4 v_BaseColor;
uniform sampler2D u_Texture;
out vec4 v_FragData;

void main() {
  vec2 v_TexturePosition = v_TextureMin + mod(v_TexturePos, v_TextureDim);
  v_FragData = clamp(texture(u_Texture, v_TexturePosition) + v_BaseColor, 0.0, 1.0);
}`

const VERTEX = `#version 150

#define MAX_TILES 1024

struct Tile {
  vec4 texture;
};

layout (std140) uniform TextureData {
  Tile Tiles[MAX_TILES];
};

in vec3 v_Position;
in vec2 v_Texture;
in float f_VertexFrame;
in float f_InstanceFrame;
in vec4 v_Color;
in mat4 m_Model;
uniform mat4 m_View;
uniform mat4 m_Projection;
out vec2 v_TexturePos;
out vec2 v_TextureMin;
out vec2 v_TextureDim;
out vec4 v_BaseColor;

void main() {
  Tile t_Tile = Tiles[int(f_VertexFrame + f_InstanceFrame)];
  v_TextureMin = t_Tile.texture.zw;
  v_TextureDim = t_Tile.texture.xy;
  v_TexturePos = v_Texture * v_TextureDim;
  v_BaseColor = v_Color;
  gl_Position = m_Projection * m_View * m_Model * vec4(v_Position, 1.0);
}`

type renderInstance struct {
	model mgl32.Mat4
	frame float32
	color mgl32.Vec4
}

type Renderer struct {
	shader      *core.Program
	vbo         *core.ArrayBuffer
	textureData *core.UniformBlock
	uView       *core.Uniform
	uProj       *core.Uniform
	bufferSize  int
	buffer      []renderInstance
	stride      uintptr
}

func NewRenderer(bufferSize int) (r *Renderer, err error) {
	var (
		instance       renderInstance
		instanceStride = unsafe.Sizeof(instance)
	)
	r = &Renderer{
		shader:     core.NewProgram(),
		bufferSize: bufferSize,
		buffer:     make([]renderInstance, bufferSize),
		stride:     instanceStride,
	}
	if err = r.shader.Load(VERTEX, FRAGMENT); err != nil {
		return
	}
	r.shader.Bind()

	r.vbo = core.NewArrayBuffer()

	r.shader.Attrib("f_InstanceFrame", instanceStride).Float(unsafe.Offsetof(instance.frame), 1)
	r.shader.Attrib("m_Model", instanceStride).Mat4(unsafe.Offsetof(instance.model), 1)
	r.shader.Attrib("v_Color", instanceStride).Vec4(unsafe.Offsetof(instance.color), 1)

	r.textureData = r.shader.UniformBlock("TextureData", 1)

	r.uView = r.shader.Uniform("m_View")
	r.uProj = r.shader.Uniform("m_Projection")

	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (r *Renderer) Bind() {
	r.shader.Bind()
}

func (r *Renderer) registerGeometry(geometry *Geometry) {
	var (
		pt       Point
		ptStride = unsafe.Sizeof(pt)
	)
	geometry.Bind()
	geometry.Upload()
	r.shader.Attrib("v_Position", ptStride).Vec3(unsafe.Offsetof(pt.Position), 0)
	r.shader.Attrib("v_Texture", ptStride).Vec2(unsafe.Offsetof(pt.Texture), 0)
	r.shader.Attrib("f_VertexFrame", ptStride).Float(unsafe.Offsetof(pt.Frame), 0)
}

func (r *Renderer) registerTextureData(buffer UniformBufferSheet) {
	r.textureData.Bind(buffer.BufferID(), buffer.Size())
}

func (r *Renderer) Unbind() {
	r.shader.Unbind()
}

func (r *Renderer) Delete() {
	if r.shader != nil {
		r.shader.Delete()
		r.shader = nil
	}
	if r.vbo != nil {
		r.vbo.Delete()
		r.vbo = nil
	}
}

func (r *Renderer) draw(geometry *Geometry, count int) (err error) {
	if count <= 0 {
		return
	}
	r.vbo.Upload(r.buffer, count*int(r.stride))
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(len(geometry.Points)), int32(count))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (r *Renderer) Render(
	camera *core.Camera,
	sheet UniformBufferSheet,
	geometry *Geometry,
	instances Instances,
) (err error) {
	var (
		instance *Instance
		i        *renderInstance
		index    int
	)
	r.uView.Mat4(camera.View)
	r.uProj.Mat4(camera.Projection)
	r.registerGeometry(geometry)
	r.registerTextureData(sheet)
	index = 0
	instance = instances.Head()
	for instance != nil {
		i = &r.buffer[index]
		i.frame = float32(instance.Frame)
		i.model = instance.GetModel()
		i.color = instance.Color()
		index++
		instance = instance.Next()
		if index >= r.bufferSize {
			if err = r.draw(geometry, index); err != nil {
				return
			}
			index = 0
		}
	}
	err = r.draw(geometry, index)
	return
}

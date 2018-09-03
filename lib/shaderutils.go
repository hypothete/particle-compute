package cmp

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
)

// ShaderProgram is a struct describing a set of shaders compiled into a program
type ShaderProgram struct {
	ID      uint32
	Buffers []uint32
}

// CreateShaderProgram creates the gl reference for the ShaderProgram
func CreateShaderProgram() ShaderProgram {
	var sp = ShaderProgram{ID: gl.CreateProgram()}
	return sp
}

// Attach is a wrapper for gl.AttachShader
func (sp ShaderProgram) Attach(shaderID uint32) {
	gl.AttachShader(sp.ID, shaderID)
}

// Link is a wrapper for gl.LinkProgram
func (sp ShaderProgram) Link() {
	gl.LinkProgram(sp.ID)
}

// Load takes a file path and tries to load it as a shader
func Load(path string, shaderType uint32) uint32 {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	source := string(bytes) + "\x00"

	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		panic(fmt.Errorf("failed to compile %v: %v", source, log))
	}

	return shader
}

// MakeVao initializes and returns a vertex array from the points provided.
func MakeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	return vao
}

// MakeOutputTexture returns a pointer to a texture of width and height
func MakeOutputTexture(renderedTexture uint32, width, height int32) uint32 {
	gl.GenTextures(1, &renderedTexture)
	gl.BindTexture(gl.TEXTURE_2D, renderedTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, width, height, 0, gl.RGBA, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return renderedTexture
}

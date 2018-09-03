package cmp

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Camera is used to view the scene
type Camera struct {
	Projection, View, InvViewProd mgl32.Mat4
	Position, Target, Up          mgl32.Vec3
	Fovy, Aspect, Near, Far       float32
	ViewUniform, ProjUniform      int32
}

// UpdateMatrices reads values form the camera and preps the matrices
func (c *Camera) UpdateMatrices() {
	c.View = mgl32.LookAtV(c.Position, c.Target, c.Up)
	c.Projection = mgl32.Perspective(c.Fovy, c.Aspect, c.Near, c.Far)
}

func (c *Camera) AssignUniformLocations() {
	c.ProjUniform = int32(1)
	c.ViewUniform = int32(2)
}

func (c *Camera) SetUniforms() {
	gl.UniformMatrix4fv(c.ViewUniform, 1, false, &c.View[0])
	gl.UniformMatrix4fv(c.ProjUniform, 1, false, &c.Projection[0])
}

// NewCamera is a camera constructor
func NewCamera(position, target, up mgl32.Vec3, fovy, aspect, near, far float32) Camera {
	c := new(Camera)
	c.Position = position
	c.Target = target
	c.Up = up
	c.Fovy = fovy
	c.Aspect = aspect
	c.Near = near
	c.Far = far
	c.UpdateMatrices()
	return *c
}

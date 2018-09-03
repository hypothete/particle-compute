package main

import (
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hypothete/particle-compute/lib"
)

const (
	windowWidth  = 1024
	windowHeight = 1024
	widthUnits   = 1024
	heightUnits  = 1
	numParticles = 1024
)

var points, velocities []mgl32.Vec4
var vao uint32

// gStr is a shorthand for goofy string concat
func gStr(str string) *uint8 {
	formatted := gl.Str(str + "\x00")
	return formatted
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Compute", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}

func takeScreenshot() {
	pixels := image.NewRGBA(image.Rect(0, 0, windowWidth, windowHeight))
	gl.ReadPixels(0, 0, windowWidth, windowHeight, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels.Pix))
	newImage := imaging.FlipV(pixels)
	output, err := os.Create("screenshot.png")
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(output, newImage)
	output.Close()
}

func draw(
	window *glfw.Window,
	particleProg *cmp.ShaderProgram,
	quadProg *cmp.ShaderProgram,
	cam *cmp.Camera) {

	gl.UseProgram(particleProg.ID)
	gl.DispatchCompute(widthUnits, heightUnits, 1)
	gl.MemoryBarrier(gl.VERTEX_ATTRIB_ARRAY_BARRIER_BIT)

	gl.UseProgram(quadProg.ID)
	//gl.ClearColor(0, 0, 0, 1)
	//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.POINTS, 0, numParticles)

	gl.UseProgram(0)

	window.SwapBuffers()

	if window.GetKey(glfw.KeyF3) == glfw.Press {
		takeScreenshot()
	}
	glfw.PollEvents()
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	initOpenGL()

	particleShader := cmp.Load("shaders/particles.glsl", gl.COMPUTE_SHADER)
	vertexShader := cmp.Load("shaders/vert.glsl", gl.VERTEX_SHADER)
	fragmentShader := cmp.Load("shaders/frag.glsl", gl.FRAGMENT_SHADER)

	particleProg := cmp.CreateShaderProgram()
	particleProg.Attach(particleShader)
	particleProg.Link()
	gl.UseProgram(particleProg.ID)

	posSSBO := uint32(1)
	velSSBO := uint32(2)
	particleProg.Buffers = append(particleProg.Buffers, posSSBO, velSSBO)

	// Generate the position buffer
	mm := float32(32)
	for i := 0; i < numParticles; i++ {
		x := (rand.Float32()*2 - 1) * mm
		y := (rand.Float32()*2 - 1) * mm
		z := (rand.Float32()*2 - 1) * mm
		points = append(points, mgl32.Vec4{x, y, z, 1})
	}

	gl.GenBuffers(1, &posSSBO)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, posSSBO)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, numParticles*16, gl.Ptr(points), gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, posSSBO)

	//generate the velocity buffer
	nn := float32(6.0)
	for i := 0; i < numParticles; i++ {
		x := (rand.Float32()*2 - 1) * nn
		y := (rand.Float32()*2 - 1) * nn
		z := (rand.Float32()*2 - 1) * nn
		velocities = append(velocities, mgl32.Vec4{x, y, z, 0})
	}

	gl.GenBuffers(1, &velSSBO)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, velSSBO)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, numParticles*16, gl.Ptr(velocities), gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, velSSBO)

	// set up the quad for viewing
	quadProg := cmp.CreateShaderProgram()
	quadProg.Attach(vertexShader)
	quadProg.Attach(fragmentShader)
	quadProg.Link()

	gl.UseProgram(quadProg.ID)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, posSSBO)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	cam := cmp.NewCamera(
		mgl32.Vec3{0, 0, 100},
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 1, 0},
		mgl32.DegToRad(60),
		windowWidth/windowHeight,
		0.1,
		1000.0)
	cam.AssignUniformLocations()
	cam.SetUniforms()

	gl.UseProgram(0)

	for !window.ShouldClose() {
		draw(window, &particleProg, &quadProg, &cam)
	}
}

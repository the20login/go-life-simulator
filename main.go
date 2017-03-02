package main

import (
	"image/color"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/llgcode/draw2d/draw2dkit"
	"time"
	"github.com/the20login/go-life-simulator/world"
	"github.com/the20login/go-life-simulator/world/quadtree"
	"github.com/the20login/go-life-simulator/dweller"
	"math"
	"math/rand"
)

var (
	// global rotation
	rotate        int
	width, height int
	redraw        = true
	font          draw2d.FontData
	worldInstance *world.World
)

func reshape(window *glfw.Window, w, h int) {
	gl.ClearColor(1, 1, 1, 1)
	/* Establish viewing area to cover entire window. */
	gl.Viewport(0, 0, int32(w), int32(h))
	/* PROJECTION Matrix mode. */
	gl.MatrixMode(gl.PROJECTION)
	/* Reset project matrix. */
	gl.LoadIdentity()
	/* Map abstract coords directly to window coords. */
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	/* Invert Y axis so increasing Y goes down. */
	gl.Scalef(1, -1, 1)
	/* Shift origin up to upper-left corner. */
	gl.Translatef(0, float32(-h), 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	width, height = w, h
	redraw = true
}

// Ask to refresh
func invalidate() {
	redraw = true
}

func display() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.LineWidth(1)
	gc := draw2dgl.NewGraphicContext(width, height)
	gc.SetFontData(draw2d.FontData{
		Name:   "luxi",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleBold | draw2d.FontStyleItalic})

	gc.BeginPath()
	scale := math.Min(float64(width) / float64(worldInstance.Rectangle.Width()), float64(height) / float64(worldInstance.Rectangle.Height()))
	gc.SetFillColor(color.RGBA{0, 255, 0, 0xff})
	for _, point := range worldInstance.FoodPositions() {
		draw2dkit.Rectangle(gc, point.X * scale, point.Y * scale, (point.X + 1) * scale, (point.Y + 1) * scale)
		gc.Fill()
	}
	gc.SetFillColor(color.RGBA{255, 0, 0, 0xff})
	for _, point := range worldInstance.AntsPositions() {
		draw2dkit.Rectangle(gc, point.X * scale, point.Y * scale, (point.X + 1) * scale, (point.Y + 1) * scale)
		gc.Fill()
	}

	gl.Flush() /* Single buffered, so needs a flush. */
}

func init() {
	runtime.LockOSThread()
}

func main() {
	worldInstance = world.NewWorld(quadtree.NewRectangle(quadtree.Point{0, 0}, quadtree.Point{100, 100}))
	go func(){
		time.Sleep(time.Second)
		for i:=0;i< 20;i++ {
			worldInstance.AddFood(quadtree.Point{rand.Float64() * 100, rand.Float64() * 100}, dweller.NewFood(worldInstance, quadtree.Point{50, 50}))
		}
		time.Sleep(5*time.Second)
		worldInstance.AddAnt(quadtree.Point{55, 55}, dweller.NewAnt(worldInstance, quadtree.Point{50, 50}))
	}()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	width, height = 800, 800
	window, err := glfw.CreateWindow(width, height, "Show RoundedRect", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetSizeCallback(reshape)
	window.SetKeyCallback(onKey)
	window.SetCharCallback(onChar)

	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	reshape(window, width, height)
	for !window.ShouldClose() {

			display()
			window.SwapBuffers()

		time.Sleep(100 * time.Millisecond)
		glfw.PollEvents()
	}
}

func onChar(w *glfw.Window, char rune) {
	log.Println(char)
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch {
	case key == glfw.KeyEscape && action == glfw.Press,
		key == glfw.KeyQ && action == glfw.Press:
		w.SetShouldClose(true)
	}
}

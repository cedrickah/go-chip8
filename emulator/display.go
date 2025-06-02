package emulator

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Renderer struct {
	Cols    uint8
	Rows    uint8
	Display []int
	window *sdl.Window
	surface *sdl.Surface
}

func NewRenderer(w *sdl.Window, s *sdl.Surface) *Renderer{
	colsNumber := uint8(64)
	rowsNumber := uint8(32)
	displayLength := uint16(colsNumber)*uint16(rowsNumber)
	cells := make([]int, displayLength)
	return &Renderer{Cols: colsNumber, Rows: rowsNumber, Display: cells, window: w, surface: s}
}

func (r *Renderer) SetPixel(x uint8, y uint8) bool {
	pixelLoc := uint16(x) + (uint16(y) * uint16(r.Cols))
	r.Display[pixelLoc] ^= 1
	return r.Display[pixelLoc] != 0
}

func (r *Renderer) Clear() {
	displayLength := uint16(r.Cols)*uint16(r.Rows)
	r.Display = make([]int, displayLength)
}

func (r *Renderer) Render() {
	r.surface.FillRect(nil, 0)
	
	displayLength := uint16(r.Cols)*uint16(r.Rows)
	for i := uint16(0); i < displayLength; i++ {
		x := (uint16(i) % uint16(r.Cols)) * 10
		y := math.Floor(float64(uint16(i)) / float64(uint16(r.Cols))) * 10
		if r.Display[i] != 0 {
			colour := sdl.Color{R: 255, G: 255, B: 255, A: 255}
			rect := sdl.Rect{X: int32(x), Y: int32(y), W: 10, H: 10}
			pixel := sdl.MapRGBA(r.surface.Format, colour.R, colour.G, colour.B, colour.A)
			r.surface.FillRect(&rect, pixel)
		}
	}
	r.window.UpdateSurface()
}

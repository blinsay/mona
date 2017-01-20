package main

import (
	"image"
	"image/color"
	"math/rand"
)

// Generates glyphs and applies them to
type GlyphGenerator interface {
	Init()
	Apply(img *image.Gray)
}

type RandomRectangles struct {
	Bounds   image.Rectangle
	MaxCount int

	nextUint func() uint8
	nextInt  func(int) int
}

func GenerateRects(count int, bounds image.Rectangle) *RandomRectangles {
	return &RandomRectangles{
		Bounds:   bounds,
		MaxCount: count,
		nextUint: func() uint8 {
			return uint8(rand.Uint32() & 0xff)
		},
		nextInt: func(n int) int {
			return rand.Intn(n)
		},
	}
}

func (r *RandomRectangles) Apply(img *image.Gray) {
	n := r.nextInt(r.MaxCount)
	for i := 0; i < n; i++ {
		DrawRect(img, r.randomRectangle(), r.randomColor())
	}
}

// FIXME: test!!!!
func (r *RandomRectangles) randomPoint(min, max image.Point) image.Point {
	x := min.X + r.nextInt(max.X-min.X)
	y := min.Y + r.nextInt(max.Y-min.Y)
	return image.Point{x, y}
}

func (r *RandomRectangles) randomRectangle() image.Rectangle {
	min := r.randomPoint(r.Bounds.Min, r.Bounds.Max)
	max := r.randomPoint(min, r.Bounds.Max)
	return image.Rectangle{Min: min, Max: max}
}

func (r *RandomRectangles) randomColor() color.Color {
	return color.Gray{r.nextUint()}
}

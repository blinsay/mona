package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
)

//  NOTE: This would be dope to have, but the compiler can't figure out that
//  func(image.Rectangle)*image.Grey should be the same as func(image.Rectangle) draw.Image
//
// func ConvertTo(src image.Image, imgType func(image.Rectangle) draw.Image) draw.Image {
// 	bounds := src.Bounds()
// 	dst := imgType(bounds)
// 	draw.Draw(dst, bounds, src, bounds.Min, draw.Over)
// 	return dst
// }
//
// func NewGray(r image.Rectangle) draw.Image {
// 	return image.NewGray(r)
// }

// Create a new image with the given bounds. The whole image will be uniformly
// colored with the color c.
func NewImage(bounds image.Rectangle, c color.Color) *image.Gray {
	img := image.NewGray(bounds)
	draw.Draw(img, bounds, &image.Uniform{C: c}, bounds.Min, draw.Src)
	return img
}

// Draw a rectangle onto the given image. The rectangle will be filled with the
// given color.
func DrawRect(img draw.Image, r image.Rectangle, c color.Color) {
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			img.Set(x, y, c)
		}
	}
}

// Reads the file at the given path as a JPEG. Returns any errors from opening
// the file or decoding the jpeg.
func ReadJpeg(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return jpeg.Decode(f)
}

// Convert any image.Image to an image.Gray.
func ConvertToGray(src image.Image) *image.Gray {
	bounds := src.Bounds()
	dst := image.NewGray(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Over)
	return dst
}

// Destructively copy src to dst. This skips draw.Draw in favor of directly
// copying struct fields.
//
// This might be a bad idea.
func CopyGray(dst, src *image.Gray) {
	dst.Stride = src.Stride
	dst.Rect = src.Rect
	dst.Pix = make([]uint8, len(src.Pix))
	copy(dst.Pix, src.Pix)
}

package main

import (
	"image"
	"math"
)

func EuclideanDistanceTo(source *image.Gray) func(*image.Gray) float64 {
	return func(other *image.Gray) float64 {
		if source.Bounds() != other.Bounds() {
			return math.Inf(1)
		}

		var sumSquares float64
		for i := range source.Pix {
			delta := float64(source.Pix[i]) - float64(other.Pix[i])
			sumSquares += delta * delta
		}

		return math.Sqrt(sumSquares)
	}
}

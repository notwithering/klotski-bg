package render

import (
	"main/layout"
	"math"
)

const (
	degree float64 = math.Pi / 180

	distance float64 = 1000
	fov      float64 = 30 * degree

	width  = 2560
	height = 1440

	rotation float64 = -2 * degree
)

type Point2D struct {
	X, Y       int
	Depth      float64
	LocalDepth float64
}

func Rasterize(positions []layout.Vec3) []Point2D {
	out := make([]Point2D, len(positions))
	if len(positions) == 0 {
		return out
	}

	minZ, maxZ := math.Inf(1), math.Inf(-1)
	for _, p := range positions {
		p = rotateY(p, rotation)

		z := p.Z + distance
		if z < minZ {
			minZ = z
		}
		if z > maxZ {
			maxZ = z
		}
	}
	if maxZ-minZ < 1e-6 {
		maxZ = minZ + 1
	}

	scale := float64(height) / (2 * math.Tan(fov/2))
	for i, p := range positions {
		p = rotateY(p, rotation)

		depth := p.Z + distance
		if depth < 1 {
			depth = 1
		}

		x := int((p.X/depth)*scale + float64(width)/2)
		y := int((-p.Y/depth)*scale + float64(height)/2)

		out[i] = Point2D{
			X:          x,
			Y:          y,
			Depth:      depth,
			LocalDepth: (depth - minZ) / (maxZ - minZ),
		}
	}
	return out
}

func rotateY(p layout.Vec3, theta float64) layout.Vec3 {
	s, c := math.Sin(theta), math.Cos(theta)
	return layout.Vec3{
		X: p.X*c + p.Z*s,
		Y: p.Y,
		Z: -p.X*s + p.Z*c,
	}
}

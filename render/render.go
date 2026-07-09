package render

import (
	"image"
	"image/color"
	"math"
)

func Render(points []Point2D, edges [][2]int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			img.SetRGBA(x, y, color.RGBA{R: 40, G: 40, B: 40, A: 255})
		}
	}

	for _, e := range edges {
		p0 := points[e[0]]
		p1 := points[e[1]]

		drawLine(img, p0.X, p0.Y, p1.X, p1.Y, color.RGBA{R: 235, G: 219, B: 178, A: 50})
	}
	for _, p := range points {
		radius := int(math.Round(3 * (1 - p.LocalDepth)))
		normalizedDepth := math.Max(0, math.Min(1, p.LocalDepth))

		pointColor := gradient(normalizedDepth,
			color.RGBA{R: 204, G: 36, B: 29, A: 255},
			color.RGBA{R: 152, G: 151, B: 26, A: 255},
			color.RGBA{R: 69, G: 133, B: 136, A: 255},
		).(color.RGBA)

		drawCircle(img, p.X, p.Y, radius, pointColor)
	}

	return img
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, color color.RGBA) {
	dx := int(math.Abs(float64(x1 - x0)))
	dy := -int(math.Abs(float64(y1 - y0)))
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy

	for {
		blend(img, x0, y0, color)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func drawCircle(img *image.RGBA, x, y, r int, c color.RGBA) {
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy <= r*r {
				blend(img, x+dx, y+dy, c)
			}
		}
	}
}

func blend(img *image.RGBA, x, y int, c color.RGBA) {
	b := img.Bounds()
	if x < b.Min.X || y < b.Min.Y || x >= b.Max.X || y >= b.Max.Y {
		return
	}
	if c.A == 255 {
		img.SetRGBA(x, y, c)
		return
	}
	bgc := img.RGBAAt(x, y)
	a := float64(c.A) / 255
	mix := func(cv, bgv uint8) uint8 {
		return uint8(float64(cv)*a + float64(bgv)*(1-a))
	}
	img.SetRGBA(x, y, color.RGBA{R: mix(c.R, bgc.R), G: mix(c.G, bgc.G), B: mix(c.B, bgc.B), A: 255})
}

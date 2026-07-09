package layout

import (
	"math"
	"math/rand"
)

const (
	iterations int     = 300
	area       float64 = 50
	gravity    float64 = 0.2
	seed       int64   = 1
)

type Vec3 struct {
	X, Y, Z float64
}

// 3d fruchterman-reingold sim
func ForceDirected3D(n int, edges [][2]int) []Vec3 {
	pos := make([]Vec3, n)
	if n == 0 {
		return pos
	}

	rng := rand.New(rand.NewSource(seed))
	for i := range pos {
		pos[i] = Vec3{
			X: (rng.Float64()*2 - 1) * area / 2,
			Y: (rng.Float64()*2 - 1) * area / 2,
			Z: (rng.Float64()*2 - 1) * area / 2,
		}
	}

	volume := area * area * area
	k := math.Cbrt(volume / float64(n))

	disp := make([]Vec3, n)
	temp := area / 10

	for range iterations {
		for i := range disp {
			disp[i] = Vec3{}
		}

		for i := range n {
			for j := i + 1; j < n; j++ {
				dx := pos[i].X - pos[j].X
				dy := pos[i].Y - pos[j].Y
				dz := pos[i].Z - pos[j].Z
				dist := math.Sqrt(dx*dx+dy*dy+dz*dz) + 1e-6
				f := (k * k) / dist
				ux, uy, uz := dx/dist, dy/dist, dz/dist
				disp[i].X += ux * f
				disp[i].Y += uy * f
				disp[i].Z += uz * f
				disp[j].X -= ux * f
				disp[j].Y -= uy * f
				disp[j].Z -= uz * f
			}
		}

		for _, e := range edges {
			i, j := e[0], e[1]
			dx := pos[i].X - pos[j].X
			dy := pos[i].Y - pos[j].Y
			dz := pos[i].Z - pos[j].Z
			dist := math.Sqrt(dx*dx+dy*dy+dz*dz) + 1e-6
			f := (dist * dist) / k
			ux, uy, uz := dx/dist, dy/dist, dz/dist
			disp[i].X -= ux * f
			disp[i].Y -= uy * f
			disp[i].Z -= uz * f
			disp[j].X += ux * f
			disp[j].Y += uy * f
			disp[j].Z += uz * f
		}

		for i := range pos {
			disp[i].X -= pos[i].X * gravity
			disp[i].Y -= pos[i].Y * gravity
			disp[i].Z -= pos[i].Z * gravity

			dlen := math.Sqrt(disp[i].X*disp[i].X+disp[i].Y*disp[i].Y+disp[i].Z*disp[i].Z) + 1e-6
			capped := math.Min(dlen, temp)
			pos[i].X += disp[i].X / dlen * capped
			pos[i].Y += disp[i].Y / dlen * capped
			pos[i].Z += disp[i].Z / dlen * capped
		}

		temp *= 0.99
	}

	return pos
}

package genimg

import (
	"math/rand"
)

// Gene of every individual represented as a polygon coordinates and color
type Gene struct {
	x, y, radius, angle float64
	red, green, blue    int
}

// NewRandomGene creates a gene with all values set to random
func NewRandomGene() Gene {
	return Gene{
		x:      rand.Float64(),
		y:      rand.Float64(),
		radius: rand.Float64(),
		angle:  rand.Float64(),
		red:    rand.Intn(255),
		green:  rand.Intn(255),
		blue:   rand.Intn(255),
	}
}

// Mutate creates a new gene with one of the fields changed
func (g *Gene) Mutate(power float64) {
	switch rand.Intn(7) {
	case 0:
		g.x = mutateFloat(g.x, power)
	case 1:
		g.y = mutateFloat(g.y, power)
	case 2:
		g.radius = mutateFloat(g.radius, power)
	case 3:
		g.angle = mutateFloat(g.angle, power)
	case 4:
		g.red = mutateColor(g.red, power)
	case 5:
		g.green = mutateColor(g.green, power)
	case 6:
		g.blue = mutateColor(g.blue, power)
	}
}

func ensureRange(x, lo, hi float64) float64 {
	if x > hi {
		return hi
	}
	if x < lo {
		return lo
	}
	return x
}

func mutateFloat(f, power float64) float64 {
	if f == 0 {
		return rand.Float64() * power
	}

	delta := rand.Float64()*power*2 - power
	return ensureRange(f+delta, 0.0, 1.0)
}

func mutateColor(c int, power float64) int {
	nc := c
	for nc == c {
		nc = int(255.0 * mutateFloat(float64(c)/255.0, power))
	}
	return nc
}

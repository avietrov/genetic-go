package genimg

import (
	"math/rand"
)

// Gene of every individual represented as a polygon coordinates and color
type Gene struct {
	x, y, radius, angle float64
	red, green, blue    float64
}

// NewRandomGene creates a gene with all values set to random
func NewRandomGene() Gene {
	return Gene{
		x:      rand.Float64(),
		y:      rand.Float64(),
		radius: rand.Float64(),
		angle:  rand.Float64(),
		red:    rand.Float64(),
		green:  rand.Float64(),
		blue:   rand.Float64(),
	}
}

func (g *Gene) mutate() Gene {
	newGene := Gene{
		x:      g.x,
		y:      g.y,
		radius: g.radius,
		angle:  g.angle,
		red:    g.red,
		green:  g.green,
		blue:   g.blue,
	}

	switch rand.Intn(7) {
	case 0:
		newGene.x = mutateFloat(newGene.x)
	case 1:
		newGene.y = mutateFloat(newGene.y)
	case 2:
		newGene.radius = mutateFloat(newGene.radius)
	case 3:
		newGene.angle = mutateFloat(newGene.angle)
	case 4:
		newGene.red = mutateFloat(newGene.red)
	case 5:
		newGene.green = mutateFloat(newGene.green)
	case 6:
		newGene.blue = mutateFloat(newGene.blue)
	}

	return newGene
}

func ensureRange(x float64, lo float64, hi float64) float64 {
	if x > hi {
		return hi
	}
	if x < lo {
		return lo
	}
	return x
}

func mutateFloat(f float64) float64 {
	span := 1.0 * mutationPower
	if f == 0 {
		return rand.Float64() * span
	}

	delta := rand.Float64()*span*2 - span
	return ensureRange(f+delta, 0.0, 1.0)
}

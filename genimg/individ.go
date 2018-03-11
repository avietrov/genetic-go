package genimg

import (
	"math/rand"

	"github.com/fogleman/gg"
)

// Individ represents an image in a form of an array of genes
type Individ struct {
	gene []Gene
	fit  int
}

func (i *Individ) mutate() {
	indexes := rand.Perm(len(i.gene))
	for j := 0; j < genesToMutate; j++ {
		idx := indexes[j]
		i.gene[idx] = i.gene[idx].mutate()
	}
}

func (i *Individ) render(w int, h int) *gg.Context {
	ctx := gg.NewContext(int(w), int(h))
	ctx.SetRGBA(1, 1, 1, 1)
	ctx.DrawRectangle(0, 0, 1.0*float64(w), 1.0*float64(h))
	ctx.Fill()

	for k := range i.gene {
		gene := i.gene[k]
		x := gene.x * float64(w)
		y := gene.y * float64(h)
		r := gene.radius * (float64(w) * maxPolygonSize)
		a := gene.angle
		ctx.DrawRegularPolygon(3, x, y, r, a)
		ctx.SetRGBA255(gene.red, gene.green, gene.blue, alpha)
		ctx.Fill()
	}

	return ctx
}

// ByFitness is used to sort population by their score
type ByFitness []Individ

func (a ByFitness) Len() int           { return len(a) }
func (a ByFitness) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFitness) Less(i, j int) bool { return a[i].fit < a[j].fit }

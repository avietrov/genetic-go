package genimg

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sort"

	"github.com/fogleman/gg"
)

const xOverCount = 1
const populSize = 30
const genomSize = 10
const mutationPower = 0.5
const maxPolygonSize = 0.3
const alpha = 127
const genesToMutate = 1
const precision = 0.1

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
		newGene.red = mutateColor(newGene.red)
	case 5:
		newGene.green = mutateColor(newGene.green)
	case 6:
		newGene.blue = mutateColor(newGene.blue)
	}

	return newGene
}

func mutateColor(c int) int {
	nc := c
	for nc == c {
		nc = int(float64(c) * (rand.Float64()*mutationPower*2 + (1 - mutationPower)))
	}

	if nc > 255 {
		return 255
	}

	if nc < 0 { // should not be possible?
		return 0
	}

	return nc

}

func mutateFloat(f float64) float64 {
	nf := f
	for nf == f {
		nf = f * (rand.Float64()*mutationPower*2 + (1 - mutationPower))
	}

	if nf > 1.0 {
		return 1.0
	}

	if nf < 0.0 { // should not be possible?
		return 0.0
	}

	return nf
}

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

func (i *Individ) render(w float64, h float64) *gg.Context {
	ctx := gg.NewContext(int(w), int(h))
	ctx.SetRGBA(1, 1, 1, 1)
	ctx.DrawRectangle(0, 0, 1.0*float64(w), 1.0*float64(h))
	ctx.Fill()

	for k := range i.gene {
		gene := i.gene[k]
		ctx.DrawRegularPolygon(3, gene.x*w, gene.y*h, gene.radius*w*maxPolygonSize, gene.angle)
		ctx.SetRGBA255(gene.red, gene.green, gene.blue, alpha)
		ctx.Fill()
	}

	return ctx
}

func compareColors(c1 color.Color, c2 color.Color) int {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	rd := r1 - r2
	gd := g1 - g2
	bd := b1 - b2
	return int(rd*rd + gd*gd + bd*bd)
}

func elite(popul []Individ) Individ {
	idx := int(rand.Float32() * rand.Float32() * float32(len(popul)-1))
	return popul[idx]
}

func crossover(p1 []Gene, p2 []Gene) []Gene {
	child := make([]Gene, len(p1))
	lo := rand.Intn(len(p1) - 1)
	hi := lo + rand.Intn(len(p1)-lo)
	for i := 0; i < len(child); i++ {
		if i >= lo && i < hi {
			child[i] = p2[i]
		} else {
			child[i] = p1[i]
		}
	}
	return child
}

func fitness(ctx *gg.Context, target *gg.Context) int {
	step := int(float32(target.Width()) * precision)
	if step == 0 {
		step = 1
	}
	diff := int64(0)
	count := int64(0)
	for x := 0; x < ctx.Width(); x += step {
		for y := 0; y < ctx.Height(); y += step {
			c := ctx.Image().At(x, y)
			t := target.Image().At(x, y)
			diff += int64(compareColors(c, t))
			count++
		}
	}

	return int(diff / count)
}

// ByFitness is used to sort population by their score
type ByFitness []Individ

func (a ByFitness) Len() int           { return len(a) }
func (a ByFitness) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFitness) Less(i, j int) bool { return a[i].fit < a[j].fit }

// Main does the magic
func Main() {
	rand.Seed(42)
	tc, err := gg.LoadPNG("target.png")
	if err != nil {
		panic(err)
	}

	w := float64(tc.Bounds().Size().X)
	h := float64(tc.Bounds().Size().Y)

	target := gg.NewContext(tc.Bounds().Size().X, tc.Bounds().Size().Y)
	target.DrawImage(tc, 0, 0)
	target.Scale(precision, precision)

	// initialise population
	fmt.Print("Initializing...")
	popul := make([]Individ, populSize)
	for i := 0; i < len(popul); i++ {
		gene := make([]Gene, genomSize)
		for j := 0; j < len(gene); j++ {
			gene[j] = NewRandomGene()
		}
		popul[i] = Individ{gene, 0}
		ctx := popul[i].render(w, h)
		popul[i].fit = fitness(ctx, target)
		fmt.Print(".")
	}
	fmt.Println()

	i := 0
	bestFit := math.MaxInt64
	for popul[0].fit > 0 {
		sort.Sort(ByFitness(popul))
		if len(popul) > populSize {
			popul = popul[:populSize]
		}

		if bestFit > popul[0].fit {
			popul[0].render(w, h).SavePNG(fmt.Sprintf("out/%v.png", i))
			bestFit = popul[0].fit
		}

		if i%500 == 0 {
			fmt.Printf("%v\t%v\n", i, popul[0].fit)
			if i >= 3000 {
				return
			}
		}

		maxFit := popul[len(popul)-1].fit
		for i := 0; i < xOverCount; i++ {
			parent1 := elite(popul)
			parent2 := elite(popul)
			childGene := crossover(parent1.gene, parent2.gene)
			child := Individ{childGene, 0}
			child.mutate()
			ctx := child.render(w, h)
			child.fit = fitness(ctx, target)

			if child.fit < maxFit {
				popul = append(popul, child)
			}
		}
		i++
	}
}

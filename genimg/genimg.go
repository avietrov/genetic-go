package genimg

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"

	"github.com/fogleman/gg"
)

const xOverCount = 4
const populSize = 30
const genomSize = 300
const mutationPower = 0.5
const maxPolygonSize = 0.2
const alpha = 255
const genesToMutate = 3
const precision = 0.01

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

func fitness(ctx *gg.Context, target *gg.Context) int {
	diff := int64(0)
	count := int64(0)
	for x := 0.0; x < 1.0; x += precision {
		xS := int(float64(target.Width()) * x)
		for y := 0.0; y < 1.0; y += precision {
			yS := int(float64(target.Height()) * y)
			c := ctx.Image().At(xS, yS)
			t := target.Image().At(xS, yS)
			diff += int64(compareColors(c, t))
			count++
		}
	}

	return int(diff / count)
}

func readTarget(path string) (*gg.Context, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err2 := image.Decode(file)
	if err != nil {
		return nil, err2
	}

	ctx := gg.NewContextForImage(img)
	return ctx, nil
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

func crossoverAndMutate(popul []Individ, w int, h int, target *gg.Context) Individ {
	parent1 := elite(popul)
	parent2 := elite(popul)
	childGene := crossover(parent1.gene, parent2.gene)
	child := Individ{childGene, 0}
	child.mutate()
	ctx := child.render(w, h)
	child.fit = fitness(ctx, target)
	return child
}

func crossoverAndMutateAsync(popul []Individ, w int, h int, target *gg.Context, ch chan Individ, wg *sync.WaitGroup) {
	child := crossoverAndMutate(popul, w, h, target)
	ch <- child
	wg.Done()
}

// Main does the magic
func Main() {
	rand.Seed(42)
	target, err := readTarget("target.png")
	if err != nil {
		log.Fatal(err)
	}
	w := target.Width()
	h := target.Height()

	popul := make([]Individ, populSize)
	for i := 0; i < len(popul); i++ {
		gene := make([]Gene, genomSize)
		for j := 0; j < len(gene); j++ {
			gene[j] = NewRandomGene()
		}
		popul[i] = Individ{gene, 0}
		ctx := popul[i].render(w, h)
		popul[i].fit = fitness(ctx, target)
	}

	i := 0
	bestFit := math.MaxInt64
	for popul[0].fit > 0 {
		sort.Sort(ByFitness(popul))
		if len(popul) > populSize {
			popul = popul[:populSize]
		}

		if bestFit > popul[0].fit {
			popul[0].render(w, h).SavePNG(fmt.Sprintf("out/%07d.png", i))
			bestFit = popul[0].fit
		}

		if i%1000 == 0 {
			fmt.Printf("%v\n", popul[0].fit)
			if i >= 1000000 {
				return
			}
		}
		ch := make(chan Individ, xOverCount)
		wg := sync.WaitGroup{}
		wg.Add(xOverCount)
		for i := 0; i < xOverCount; i++ {
			go crossoverAndMutateAsync(popul, w, h, target, ch, &wg)
		}
		wg.Wait()
		close(ch)

		for child := range ch {
			popul = append(popul, child)
		}

		i++
	}
}

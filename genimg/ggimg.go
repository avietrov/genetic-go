package genimg

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/fogleman/gg"
)

const maxPolygonSize = 0.2
const alpha = 127
const precision = 0.01

func compareColors(c1 color.Color, c2 color.Color) int {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	rd := r1 - r2
	gd := g1 - g2
	bd := b1 - b2
	return int(rd*rd + gd*gd + bd*bd)
}

func FitnessTo(target *gg.Context) func(i *Individ) int64 {
	return func(i *Individ) int64 {
		ctx := Render(i, target.Width(), target.Height())
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

		return diff / count
	}
}

func ReadTarget(path string) (*gg.Context, error) {
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

func Render(i *Individ, w int, h int) *gg.Context {
	ctx := gg.NewContext(int(w), int(h))
	ctx.SetRGBA(1, 1, 1, 1)
	ctx.DrawRectangle(0, 0, 1.0*float64(w), 1.0*float64(h))
	ctx.Fill()

	for k := range i.DNA {
		gene := i.DNA[k]
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

func Observer(w, h int) func(int, []Individ) {
	return func(generation int, popul []Individ) {
		if generation%1000 == 0 {
			go Render(&popul[0], w, h).SavePNG(fmt.Sprintf("out/%07d.png", generation))
			fmt.Printf("%v\n", popul[0].Fit)
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/avietrov/genetic-go/genimg"
)

var xoverC int
var populS int
var genomS int
var mutatP float64
var mutatC int
var maxGen int
var target string

func init() {
	flag.IntVar(&xoverC, "x-over-count", 2, "number of crossovers done per generation.")
	flag.IntVar(&populS, "population-size", 30, "number of individuals per generation (how many images per generation).")
	flag.IntVar(&genomS, "genom-size", 300, "how large is DNA of each individ (how many triangles per image).")
	flag.Float64Var(&mutatP, "mutation-power", 0.2, "value from 0.0+ to to 1.0, that defines stregth of mutations.")
	flag.IntVar(&mutatC, "mutation-count", 5, "value from 1 to {genom-size} that defines how many genes are mutated per individ mutation.")
	flag.IntVar(&maxGen, "max-generations", 100000, "maximum number of generations.")
	flag.StringVar(&target, "target", "target.png", "target image")
}

func main() {
	flag.Parse()

	conf := genimg.Config{
		XOverCount:     xoverC,
		PopulSize:      populS,
		GenomSize:      genomS,
		MutationPower:  mutatP,
		GenesToMutate:  mutatC,
		MaxGenerations: maxGen,
	}

	fmt.Printf("Staring with config: %#v\n", conf)
	targetImg, err := genimg.ReadTarget(target)
	if err != nil {
		log.Fatal(err)
	}
	w := targetImg.Width()
	h := targetImg.Height()

	genimg.FindFittest(conf, genimg.FitnessTo(targetImg), genimg.Observer(w, h))
}

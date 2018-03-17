package main

import (
	"flag"
	"log"

	"github.com/avietrov/genetic-go/genimg"
)

func main() {
	targetPath := flag.String("target", "target.png", "target image")
	conf := readConfig()
	flag.Parse()
	target, err := genimg.ReadTarget(*targetPath)
	if err != nil {
		log.Fatal(err)
	}
	w := target.Width()
	h := target.Height()

	genimg.FindFittest(conf, genimg.FitnessTo(target), genimg.Observer(w, h))
}

func readConfig() genimg.Config {
	var xoverC int
	var populS int
	var genomS int
	var mutatP float64
	var mutatC int
	var maxGen int

	flag.IntVar(&xoverC, "x-over-count", 2, "number of crossovers done per generation.")
	flag.IntVar(&populS, "population-size", 30, "number of individuals per generation (how many images per generation).")
	flag.IntVar(&genomS, "genom-size", 300, "how large is DNA of each individ (how many triangles per image).")
	flag.Float64Var(&mutatP, "mutation-power", 0.2, "value from 0.0+ to to 1.0, that defines stregth of mutations.")
	flag.IntVar(&mutatC, "mutation-count", 5, "value from 1 to {genom-size} that defines how many genes are mutated per individ mutation.")
	flag.IntVar(&maxGen, "max-generations", 100000, "maximum number of generations.")

	return genimg.Config{
		XOverCount:     xoverC,
		PopulSize:      populS,
		GenomSize:      genomS,
		MutationPower:  mutatP,
		GenesToMutate:  mutatC,
		MaxGenerations: maxGen,
	}
}

package genimg

import (
	"math/rand"
	"sort"
	"sync"
)

// Config contains all the parameters for the genetic algorithm
type Config struct {
	XOverCount     int
	PopulSize      int
	GenomSize      int
	MutationPower  float64
	GenesToMutate  int
	MaxGenerations int
}

// Oracle knows how to choose best individs
type Oracle func(i *Individ) int64

func elite(popul []Individ) Individ {
	idx := int(rand.Float32() * rand.Float32() * float32(len(popul)-1))
	return popul[idx]
}

func breed(p1 Individ, p2 Individ) Individ {
	childDna := make([]Gene, len(p1.DNA))
	lo := rand.Intn(len(p1.DNA) - 1)
	hi := lo + rand.Intn(len(p1.DNA)-lo)
	for i := 0; i < len(childDna); i++ {
		if i >= lo && i < hi {
			childDna[i] = p2.DNA[i]
		} else {
			childDna[i] = p1.DNA[i]
		}
	}
	return Individ{childDna, 0}
}

func crossover(popul []Individ, genesToMutate int, power float64, oracle Oracle, ch chan Individ, wg *sync.WaitGroup) {
	parent1 := elite(popul)
	parent2 := elite(popul)
	child := breed(parent1, parent2)
	child.mutate(genesToMutate, power)
	child.Fit = oracle(&child)
	ch <- child
	wg.Done()
}

func crossoverAsync(popul []Individ, genesToMutate int, power float64, xOverCount int, oracle Oracle) []Individ {
	ch := make(chan Individ, xOverCount)
	wg := sync.WaitGroup{}
	wg.Add(xOverCount)
	for i := 0; i < xOverCount; i++ {
		go crossover(popul, genesToMutate, power, oracle, ch, &wg)
	}
	wg.Wait()
	close(ch)
	children := make([]Individ, 0)
	for child := range ch {
		children = append(children, child)
	}

	return children
}

// Observe is a function that can be used to monitor progress of the algorithm
type Observe func(generation int, popul []Individ)

// FindFittest returns the best Individ
func FindFittest(conf Config, oracle Oracle, observe Observe) Individ {
	popul := make([]Individ, conf.PopulSize)
	generation := 0
	for i := 0; i < len(popul); i++ {
		gene := make([]Gene, conf.GenomSize)
		for j := 0; j < len(gene); j++ {
			gene[j] = NewRandomGene()
		}
		popul[i] = Individ{gene, 0}
		popul[i].Fit = oracle(&popul[i])
	}

	observe(generation, popul)

	for popul[0].Fit > 0 && generation < conf.MaxGenerations {
		sort.Sort(ByFitness(popul))
		if len(popul) > conf.PopulSize {
			popul = popul[:conf.PopulSize]
		}

		newGeneration := crossoverAsync(popul, conf.GenesToMutate, conf.MutationPower, conf.XOverCount, oracle)
		popul = append(popul, newGeneration...)
		generation++
		observe(generation, popul)
	}

	return popul[0]
}

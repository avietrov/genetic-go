package genstr

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Config that sets various parameters of the experiment
type Config struct {
	Source         string
	PopulationSize int
}

// ExperimentResult represents final result and how many generations it took to achieve it
type ExperimentResult struct {
	Reuslt     string
	Iterations int
}

// Individ represents one child in a generation
type Individ struct {
	gene string
	fit  int
}

// ByFitness is used to sort population by their score
type ByFitness []Individ

func (a ByFitness) Len() int           { return len(a) }
func (a ByFitness) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFitness) Less(i, j int) bool { return a[i].fit < a[j].fit }

func rndString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func fitness(src string, candidate string) int {
	// for now assuming the same length
	fit := 0
	for i := 0; i < len(src); i++ {
		si := int(src[i])
		ci := int(candidate[i])
		fit += (si - ci) * (si - ci)
	}
	return fit
}

func mutate(p string) string {
	idx := rand.Intn(len(p))
	mutation := rand.Intn(3) - 1

	out := []rune(p)
	n := rune(int(out[idx]) + mutation)
	if n > rune('A') && n < rune('z') {
		out[idx] = n
	}

	return string(out)
}

func crossover(p1 string, p2 string) string {
	idx := rand.Intn(len(p2) - 1)
	size := 1 + rand.Intn(len(p2)-idx)

	p1r := []rune(p1)
	p2r := []rune(p2)
	for i := idx; i < idx+size; i++ {
		p1r[i] = p2r[i]
	}

	return string(p1r)
}

func elite(popul []Individ) Individ {
	idx := int(rand.Float32()*rand.Float32()*float32(len(popul)) - 1)
	return popul[idx]
}

func contains(popul []Individ, p Individ) bool {
	for _, a := range popul {
		if a == p {
			return true
		}
	}
	return false
}

// RunExperiment performs actual experiment
func RunExperiment(r *Config) ExperimentResult {
	rand.Seed(time.Now().UTC().UnixNano())

	popul := make([]Individ, r.PopulationSize)

	for i := 0; i < len(popul); i++ {
		str := rndString(len(r.Source))
		popul[i] = Individ{str, fitness(r.Source, str)}
	}

	i := 0
	for popul[0].fit > 0 {
		sort.Sort(ByFitness(popul))
		fmt.Println(popul[:5])

		parent1 := elite(popul)
		parent2 := elite(popul)
		child := crossover(parent1.gene, parent2.gene)
		child = mutate(child)

		childIndivid := Individ{child, fitness(r.Source, child)}

		if childIndivid.fit < popul[len(popul)-1].fit {
			if !contains(popul, childIndivid) {
				popul[len(popul)-1] = childIndivid
			}
		}

		i++
	}

	return ExperimentResult{popul[0].gene, i}
}

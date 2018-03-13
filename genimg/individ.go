package genimg

import "math/rand"

// Individ represents one representative in the population
type Individ struct {
	DNA []Gene
	Fit int64
}

func (i *Individ) mutate(genesToMutate int, power float64) {
	indexes := rand.Perm(len(i.DNA))
	for j := 0; j < genesToMutate; j++ {
		idx := indexes[j]
		i.DNA[idx].Mutate(power)
	}
}

// ByFitness is used to sort population by their score
type ByFitness []Individ

func (a ByFitness) Len() int           { return len(a) }
func (a ByFitness) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFitness) Less(i, j int) bool { return a[i].Fit < a[j].Fit }

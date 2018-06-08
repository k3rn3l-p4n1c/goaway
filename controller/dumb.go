package controller

import (
	"fmt"
	"github.com/MaxHalford/gago"
	m "math"
	"math/rand"
)

// A Vector contains float64s.
type Vector []float64

// Evaluate a Vector with the Drop-Wave function which takes two variables as
// input and reaches a minimum of -1 in (0, 0). The function is rather pure so
// there isn't any error handling to do.
func (X Vector) Evaluate() float64 {
	var (
		numerator   = 1 + m.Cos(12*m.Sqrt(m.Pow(X[0], 2)+m.Pow(X[1], 2)))
		denominator = 0.5*(m.Pow(X[0], 2)+m.Pow(X[1], 2)) + 2
	)
	return -numerator / denominator
}

// Mutate a Vector by resampling each element from a normal distribution with
// probability 0.8.
func (X Vector) Mutate(rng *rand.Rand) {
	gago.MutNormalFloat64(X, 0.8, rng)
}

// Crossover a Vector with another Vector by applying uniform crossover.
func (X Vector) Crossover(Y gago.Genome, rng *rand.Rand) {
	gago.CrossUniformFloat64(X, Y.(Vector), rng)
}

// Clone a Vector to produce a new one that points to a different slice.
func (X Vector) Clone() gago.Genome {
	var Y = make(Vector, len(X))
	copy(Y, X)
	return Y
}

// VectorFactory returns a random vector by generating 2 values uniformally
// distributed between -10 and 10.
func VectorFactory(rng *rand.Rand) gago.Genome {
	return Vector(gago.InitUnifFloat64(2, -10, 10, rng))
}

func Run() {
	var ga = gago.Generational(VectorFactory)
	ga.Initialize()

	fmt.Printf("Best fitness at generation 0: %f\n", ga.HallOfFame[0].Fitness)
	for i := 1; i < 10; i++ {
		err := ga.Evolve()
		if err != nil {
			fmt.Println("Handle error!")
		}
		fmt.Printf("Best fitness at generation %d: %f\n", i, ga.HallOfFame[0].Fitness)
	}
}

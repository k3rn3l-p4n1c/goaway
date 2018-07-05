package scheduler

import (
	"math/rand"
	"github.com/MaxHalford/gago"
	"fmt"
)

func QOS(m Model) float64 {
	return 1.0
}

func Capacity(m Model) bool {
	return true
}

func ModelFactory(rng *rand.Rand) gago.Genome {
	cluster := GenerateRandomCluster()
	return Model{
		cluster:     cluster,
		objectives:  []func(m Model) float64{QOS},
		constraints: []func(m Model) bool{Capacity},
	}
}

func Run() {
	var ga = gago.Generational(ModelFactory)
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

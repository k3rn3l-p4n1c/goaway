package scheduler

import (
	"math/rand"
	"github.com/MaxHalford/gago"
	"fmt"
	"container/list"
	"time"
)

func GenerateRandomCluster() (cluster Cluster) {
	const dataCenterCount = 2
	const serverCount = 6
	const deploymentCount = 20
	var dataCenters []*DataCenter
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)

	cluster.coupling = make([][]float64, deploymentCount)

	for i := 0; i < dataCenterCount; i++ {
		dc := &DataCenter{i, []*Server{}}
		dataCenters = append(dataCenters, dc)
		cluster.dataCenter = append(cluster.dataCenter, dc)

		for i := 0; i < serverCount/dataCenterCount; i++ {
			cluster.addServer(&Server{
				i,
				dc,
				*list.New(),
				20,
				8})
		}
	}
	for i := 0; i < deploymentCount; i++ {
		deployment := Deployment{
			&cluster,
			i,
			1,
			nil,
		}
		newPod := deployment.generateNewPod()
		newPod.memoryUsage = uint(r.Intn(6)) + 1
		cluster.deployments.PushFront(&deployment)
		cluster.pods.PushFront(newPod)
		deployment.podsHead = cluster.pods.Front()

		cluster.coupling[i] = make([]float64, deploymentCount)
		for j := 0; j < deploymentCount; j++ {
			if r.Intn(5) == 0 || i == j {
				cluster.coupling[i][j] = 0
			} else {
				cluster.coupling[i][j] = r.Float64()
			}
		}
	}

	return
}

func ModelFactory(rng *rand.Rand) gago.Genome {
	cluster := GenerateRandomCluster()
	return Model{
		cluster:     &cluster,
		objectives:  []func(c *Cluster) float64{rpc},
		constraints: []func(c *Cluster) bool{capacity},
	}
}

func Run() {
	var ga = gago.Generational(ModelFactory)
	ga.Initialize()

	fmt.Printf("Best fitness at generation 0: %f\n", ga.HallOfFame[0].Fitness)
	for i := 1; i < 20; i++ {
		err := ga.Evolve()
		if err != nil {
			fmt.Println("Handle error!")
		}
		fmt.Printf("Best fitness at generation %d: %f\n", i, ga.HallOfFame[0].Fitness)
	}
	c := ga.HallOfFame[0].Genome.(Model).cluster
	for _, dc := range c.dataCenter {
		fmt.Printf("Datacenter %d:\n", dc.id)
		for _, server := range dc.servers {
			fmt.Printf("\tServer %d: %d\n", server.id, server.pods.Len())
			var sum uint
			for elem := server.pods.Front(); elem != nil && elem.Value != nil; elem = elem.Next() {
				pod := elem.Value.(*Pod)
				sum += pod.memoryUsage
				fmt.Printf("\t\tPod %d:\n", pod.deployment.id)
			}
			fmt.Printf("\tUtilization: %d of %d\n\n", sum, server.memoryCap)
		}
	}
}

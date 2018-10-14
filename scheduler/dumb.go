package scheduler

import (
	"fmt"
	"github.com/MaxHalford/eaopt"
	"github.com/beevik/guid"
	"math/rand"
	"net"
	"time"
)

var databaseServer = 0
var databaseDeployment = 0

func GenerateRandomCluster() Cluster {
	const dataCenterCount = 2
	const serverCount = 6
	const deploymentCount = 10
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)

	cluster := Cluster{
		dataCenters: make([]DataCenter, dataCenterCount),
		coupling:    make([][]float64, deploymentCount),
		placement:   make([]int, deploymentCount),
		servers:     make([]Server, serverCount),
		deployments: make([]Deployment, deploymentCount),
	}
	for i := range cluster.placement {
		cluster.placement[i] = -1
	}
	cluster.placement[databaseDeployment] = databaseServer

	for i := 0; i < dataCenterCount; i++ {
		dc := DataCenter{i, fmt.Sprintf("dc-%d", i), []int{}}
		cluster.dataCenters[i] = dc
	}
	for i := 0; i < serverCount; i++ {
		dcId := r.Intn(len(cluster.dataCenters))
		cluster.servers[i] = Server{
			i,
			dcId,
			20,
			8,
			net.ParseIP(fmt.Sprintf("192.168.0.%d", i)),
			fmt.Sprintf("server-%d", i),
			Healthy,
		}
		cluster.dataCenters[dcId].serverIds = append(cluster.dataCenters[dcId].serverIds, i)
	}

	for i := 0; i < deploymentCount; i++ {
		cluster.deployments[i] = Deployment{
			&cluster,
			i,
		}
		cluster.coupling[i] = make([]float64, deploymentCount)
		for j := 0; j < deploymentCount; j++ {
			if r.Intn(5) == 0 || i == j {
				cluster.coupling[i][j] = 0
			} else {
				cluster.coupling[i][j] = r.Float64()
			}
		}
	}

	return cluster
}
func GenerateRandomStack(c *Cluster) Stack {
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)

	var stack = Stack{
		cluster:     c,
		pods:        map[string]*Pod{},
		replicaSets: make([]*ReplicaSet, len(c.deployments)),
	}

	for i, d := range c.deployments {
		replica := 1
		stack.replicaSets[d.id] = &ReplicaSet{
			d.id,
			replica,
			[]string{},
		}

		memoryUsage := uint64(r.Intn(6) + 3)

		for j := 0; j < replica; j++ {
			var serverId int

			if c.placement[i] != -1 {
				serverId = c.placement[i]
			} else {
				if server := stack.getFirstCapableServer(memoryUsage); server != nil {
					serverId = server.id
				} else {
					serverId = stack.getRandomServer().id
				}
			}
			pod := &Pod{
				guid.NewString(),
				&stack,
				d.id,
				serverId,
				memoryUsage,
			}
			stack.replicaSets[d.id].podIds = append(stack.replicaSets[d.id].podIds, pod.uuid) // add to replica set
			stack.pods[pod.uuid] = pod                                                        // add to cluster
		}
	}
	return stack
}

func Run(cluster *Cluster, stack *Stack) (res string) {
	ga, err := eaopt.NewDefaultGAConfig().NewGA()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the number of generations to run for
	ga.NGenerations = 20

	// Add a custom print function to track progress
	ga.Callback = func(ga *eaopt.GA) {
		res += fmt.Sprintf("Best fitness at generation %d: %f\n", ga.NGenerations, ga.HallOfFame[0].Fitness)
	}

	// Find an minimum
	err = ga.Minimize(func(_ *rand.Rand) eaopt.Genome {
		return Model{
			cluster: cluster,
			stack:   *stack,

			objectives:  []func(s *Stack, c *Cluster) float64{utilization, rpc},
			constraints: []func(s *Stack, c *Cluster) bool{capacity, placement},
		}

	})
	if err != nil {
		fmt.Println(err)
		return
	}

	model := ga.HallOfFame[0].Genome.(Model)
	ga.HallOfFame[0].Evaluate()
	res += fmt.Sprintf("Best fitness at last generation: %f\n", ga.HallOfFame[0].Fitness)
	if v, e := ga.HallOfFame[0].Genome.Evaluate(); e == nil {
		res += fmt.Sprintf("Best fitness at last generation: %f\n", v)
	}

	res += "\n\n" + formatCluster(model)
	return
}

func formatCluster(m Model) (res string) {
	for _, dc := range m.cluster.dataCenters {
		res += fmt.Sprintf("Datacenter #%d:\n", dc.id)
		for _, serverId := range dc.serverIds {
			server := m.cluster.servers[serverId]
			res += fmt.Sprintf("\tServer %d:\n", serverId)
			var sum uint64
			for _, pod := range m.stack.pods {
				if pod.serverId == server.id {
					sum += pod.memoryUsage
					res += fmt.Sprintf("\t\tPod %d:\n", pod.deploymentId)
				}
			}
			res += fmt.Sprintf("\tUtilization: %d of %d\n\n", sum, server.memoryCap)
		}
	}
	return
}

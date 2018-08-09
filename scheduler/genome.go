package scheduler

import (
	"math/rand"
	"github.com/MaxHalford/eaopt"
	"reflect"
	"math"
)

const bigValue = -100000.0

type Model struct {
	cluster     *Cluster
	stack       Stack
	objectives  []func(s *Stack, c *Cluster) float64
	constraints []func(s *Stack, c *Cluster) bool
}

func (m Model) Evaluate() (float64, error) {
	value := bigValue
	for _, constraint := range m.constraints {
		if !constraint(&m.stack, m.cluster) {
			return +1.0, nil
		}
	}

	for _, objective := range m.objectives {
		value *= objective(&m.stack, m.cluster)
	}
	return value, nil
}

func (m Model) Mutate(rng *rand.Rand) {
	for dIndex := 0; dIndex < len(m.cluster.deployments); dIndex++ {
		r := rng.Float32()
		if r > 0.8 {
			m.stack.scaleUp(&m.cluster.deployments[dIndex])
		} else if r < 2.0 {
			m.stack.scaleDown(&m.cluster.deployments[dIndex])
		}
	}
	for _, pod := range m.stack.pods {
		r := rng.Float32()
		if r > 0.8 {
			server := m.cluster.servers[rng.Intn(len(m.cluster.servers))]
			pod.MigrateTo(server.id)
		}
	}
}

func (m Model) Crossover(mate eaopt.Genome, rng *rand.Rand) {
	m2 := mate.(Model)
	for i := 0; i < len(m.cluster.deployments); i++ {
		var p = rng.Float64()
		d1, d2 := m.cluster.deployments[i], m2.cluster.deployments[i]
		r1, r2 := m.stack.replicaSets[i], m.stack.replicaSets[i]
		m.stack.scale(&d1, int((p*float64(r1.replica))+((1-p)*float64(r2.replica))))
		m.stack.scale(&d2, int(((1-p)*float64(r1.replica))+(p*float64(r2.replica))))
	}

	keys1 := reflect.ValueOf(m.stack.pods).MapKeys()
	keys2 := reflect.ValueOf(m2.stack.pods).MapKeys()

	for i := 0; i < int(math.Min(float64(len(keys1)), float64(len(keys2)))); i++ {
		if rng.Float64() > 0.8 {
			p1, p2 := m.stack.pods[keys1[i].String()], m2.stack.pods[keys2[i].String()]
			p1.MigrateTo(p2.serverId)
			p2.MigrateTo(p1.serverId)
		}
	}
}

func (m Model) Clone() eaopt.Genome {

	var newStack = Stack{
		cluster:     m.cluster,
		pods:        make(map[string]*Pod, len(m.stack.pods)),
		replicaSets: make([]*ReplicaSet, len(m.stack.replicaSets)),
	}

	for i, replicaSet := range m.stack.replicaSets {
		newStack.replicaSets[i] = &ReplicaSet{
			deploymentId: replicaSet.deploymentId,
			podIds:       replicaSet.podIds,
			replica:      replicaSet.replica,
		}
	}

	for podId, pod := range m.stack.pods {
		newStack.pods[podId] = &Pod{
			stack:        &newStack,
			deploymentId: pod.deploymentId,
			memoryUsage:  pod.memoryUsage,
			serverId:     pod.serverId,
		}
	}

	return Model{
		stack:       newStack,
		cluster:     m.cluster,
		constraints: m.constraints,
		objectives:  m.objectives,
	}

}
func (m Model) GetCluster() Cluster {
	return *m.cluster

}

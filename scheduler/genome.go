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
	m.stack = GenerateRandomStack(m.cluster)
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
		p1, p2 := m.stack.pods[keys1[i].String()], m2.stack.pods[keys2[i].String()]
		if rng.Float64() > 0.8 {
			s1, s2 := p1.serverId, p2.serverId
			p1.MigrateTo(s2)
			p2.MigrateTo(s1)
		}
	}
}

func (m Model) Clone() eaopt.Genome {

	var newStack = Stack{
		cluster:     m.cluster,
		pods:        make(map[string]Pod, len(m.stack.pods)),
		replicaSets: make([]ReplicaSet, len(m.stack.replicaSets)),
	}

	for i, replicaSet := range m.stack.replicaSets {
		newStack.replicaSets[i] = ReplicaSet{
			deploymentId: replicaSet.deploymentId,
			podIds:       replicaSet.podIds,
			replica:      replicaSet.replica,
		}
	}

	for podId, pod := range m.stack.pods {
		newStack.pods[podId] = Pod{
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

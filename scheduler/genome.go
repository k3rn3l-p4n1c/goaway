package scheduler

import (
	"math/rand"
	"github.com/MaxHalford/gago"
	"github.com/jinzhu/copier"
)

const bigValue = -100000

type Model struct {
	cluster     *Cluster
	objectives  []func(c *Cluster) float64
	constraints []func(c *Cluster) bool
}

func (m Model) Evaluate() (value float64) {
	value = bigValue
	for _, constraint := range m.constraints {
		if constraint(m.cluster) {
			println("constrain failed")
			return +1.0
		}
	}

	for _, objective := range m.objectives {
		value *= objective(m.cluster)
	}
	return value
}

func (m Model) Mutate(rng *rand.Rand) {
	for elem := m.cluster.pods.Front(); elem != nil; elem = elem.Next() {
		pod := elem.Value.(*Pod)
		nextServerId := rng.Intn(len(m.cluster.servers) - 1)
		pod.MigrateTo(m.cluster.servers[nextServerId])
	}

	for elem := m.cluster.deployments.Front(); elem != nil; elem = elem.Next() {
		deployment := elem.Value.(*Deployment)
		randomDiff := 3 - int(rng.Float32() * 6)
		deployment.scale(deployment.replica + randomDiff)

	}
}

func (m Model) Crossover(mate gago.Genome, rng *rand.Rand) {
	m2 := mate.(Model)
	for e1, e2 := m.cluster.deployments.Front(), m2.cluster.deployments.Front();
		e1 != nil && e2 != nil; e1, e2 = e1.Next(), e2.Next() {
		var p = rng.Float64()
		d1, d2 := e1.Value.(*Deployment), e2.Value.(*Deployment)
		d1.scale(int((p * float64(d1.replica)) + ((1 - p) * float64(d2.replica))))
		d2.scale(int(((1 - p) * float64(d1.replica)) + (p * float64(d2.replica))))
	}

	for e1, e2 := m.cluster.pods.Front(), m2.cluster.pods.Front();
		e1 != nil && e2 != nil; e1, e2 = e1.Next(), e2.Next() {
		p1, p2 := e1.Value.(*Pod), e2.Value.(*Pod)
		if rng.Float64() > 0.8 {
			s1, s2 := p1.server, p2.server
			p1.MigrateTo(s2)
			p2.MigrateTo(s1)
		}
	}

}

func (m Model) Clone() gago.Genome {
	var newIns = Model{}
	copier.Copy(&newIns, &m)
	return newIns

}
func (m Model) GetCluster() Cluster {
	return *m.cluster

}

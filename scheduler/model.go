package scheduler

import (
	"math/rand"
	"container/list"
	"github.com/MaxHalford/gago"
	"github.com/jinzhu/copier"
)

const bigValue = 10000

type Deployment struct {
	cluster  Cluster
	id       int
	replica  int
	podsHead *list.Element
}

func (d Deployment) scaleDown() {
	d.podsHead = d.podsHead.Next()
	d.cluster.pods.Remove(d.podsHead.Prev())
}

func (d Deployment) scaleUp() {
	newPod := d.generateNewPod()
	d.cluster.pods.InsertBefore(newPod, d.podsHead)
	d.podsHead = d.podsHead.Prev()
}

func (c Cluster) getRandomServer() *Server {
	serverCount := len(c.servers)
	return c.servers[rand.Intn(serverCount)]
}

func (d Deployment) generateNewPod() *Pod {
	server := d.cluster.getRandomServer()

	return &Pod{
		&d,
		server,
	}
}

func (d Deployment) scale(replica int) {
	if replica < 1 {
		return
	}

	if d.replica > replica {
		for i := 0; i < d.replica-replica; i++ {
			d.scaleDown()
		}
	} else if d.replica < replica {
		for i := 0; i < d.replica-replica; i++ {
			d.scaleUp()
		}
	}
}

type DataCenter struct {
	id int
}

type Server struct {
	id int
	dc *DataCenter
}

type Pod struct {
	deployment *Deployment
	server     *Server
}

type Evaluator func(m Cluster) float64

type Cluster struct {
	pods        list.List
	deployments list.List
	servers     []*Server
}

func GenerateRandomCluster() (cluster Cluster) {
	const dataCenterCount = 3
	const serverCount = 10
	const deploymentCount = 20
	var dataCenters []*DataCenter
	for i := 0; i < dataCenterCount; i++ {
		dataCenters = append(dataCenters, &DataCenter{i})
	}
	for i := 0; i < serverCount; i++ {
		cluster.servers = append(cluster.servers, &Server{i, dataCenters[rand.Intn(dataCenterCount)]})
	}
	for i := 0; i < deploymentCount; i++ {
		deployment := Deployment{
			cluster,
			i,
			1,
			nil,
		}
		newPod := deployment.generateNewPod()
		cluster.deployments.PushFront(deployment)
		cluster.pods.PushFront(newPod)
		deployment.podsHead = cluster.pods.Front()
	}

	return
}

type Model struct {
	cluster     Cluster
	objectives  []func(m Model) float64
	constraints []func(m Model) bool
}

func (m Model) Evaluate() (value float64) {
	for _, objective := range m.objectives {
		value += objective(m)
	}
	for _, constraint := range m.constraints {
		if constraint(m) {
			value += bigValue
		}
	}
	return
}

func (m Model) Mutate(rng *rand.Rand) {
	for elem := m.cluster.pods.Front(); elem != nil; elem = elem.Next() {
		if rng.Float32() > 0.8 {
			pod := elem.Value.(*Pod)
			serverId := pod.server.id
			nextServerId := (serverId + 1) % len(m.cluster.servers)
			pod.server = m.cluster.servers[nextServerId]
		}
	}

	for elem := m.cluster.deployments.Front(); elem != nil; elem = elem.Next() {
		randNum := rng.Float32()
		deployment := elem.Value.(Deployment)
		if randNum > 0.9 {
			deployment.scale(deployment.replica + 1)
		}
		if randNum < 0.1 {
			deployment.scale(deployment.replica - 1)
		}
	}
}

func (m Model) Crossover(mate gago.Genome, rng *rand.Rand) {

}

func (m Model) Clone() gago.Genome {
	var newIns = Model{}
	copier.Copy(&newIns, &m)
	return newIns

}

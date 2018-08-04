package scheduler

import (
	"container/list"
	"math/rand"
)

type Evaluator func(m Cluster) float64

type Deployment struct {
	cluster  *Cluster
	id       int
	replica  int
	podsHead *list.Element
}

func (d *Deployment) scaleDown() {
	d.podsHead = d.podsHead.Next()
	d.cluster.pods.Remove(d.podsHead.Prev())
}

func (d *Deployment) scaleUp() {
	newPod := d.generateNewPod()
	d.cluster.pods.InsertBefore(newPod, d.podsHead)
	d.podsHead = d.podsHead.Prev()
}

func (d *Deployment) generateNewPod() *Pod {
	server := d.cluster.getRandomServer()

	pod := &Pod{
		d,
		server,
		0,
	}
	pod.server.pods.PushFront(pod)
	return pod
}

func (d *Deployment) scale(replica int) {
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
	id      int
	servers []*Server
}

type Server struct {
	id         int
	dataCenter *DataCenter
	pods       list.List
	memoryCap  uint
	cpu        uint
}

type Pod struct {
	deployment  *Deployment
	server      *Server
	memoryUsage uint
}

func (p *Pod) MigrateTo(server *Server) {
	for elem := p.server.pods.Front(); elem != nil; elem = elem.Next() {
		if elem.Value.(*Pod) == p {
			p.server.pods.Remove(elem)
			break
		}
	}
	p.server = server
	p.server.pods.PushFront(p)
}

type Cluster struct {
	dataCenter  []*DataCenter
	pods        list.List
	deployments list.List
	servers     []*Server
	coupling    [][]float64
}

func (c *Cluster) getRandomServer() *Server {
	serverCount := len(c.servers)
	return c.servers[rand.Intn(serverCount)]
}
func (c *Cluster) addServer(server *Server) {
	c.servers = append(c.servers, server)
	server.dataCenter.servers = append(server.dataCenter.servers, server)
}

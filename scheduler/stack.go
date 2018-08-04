package scheduler

import "github.com/beevik/guid"

type ReplicaSet struct {
	deploymentId int
	replica      int
	podIds       []string
}

type Pod struct {
	uuid         string
	stack        *Stack
	deploymentId int
	serverId     int
	memoryUsage  uint
}

type Stack struct {
	cluster     *Cluster
	pods        map[string]Pod
	replicaSets []ReplicaSet
}

func (s Stack) scaleDown(d *Deployment) {
	rs := s.replicaSets[d.id]
	if rs.replica <= 1 {
		return
	}
	deletedPodId := rs.podIds[len(rs.podIds)-1]
	rs.replica --
	rs.podIds = rs.podIds[:len(rs.podIds)-1]                  // remove from deployment
	delete(s.pods, deletedPodId)                              // remove from cluster
}

func (s Stack) scaleUp(d *Deployment) {
	newPod := s.generateNewPod(d)
	rs := s.replicaSets[d.id]
	rs.podIds = append(rs.podIds, newPod.uuid) // add to replica set
	s.pods[newPod.uuid] = newPod           // add to stack
}

func (s Stack) scale(d *Deployment, replica int) {
	if replica < 1 {
		return
	}
	rs := s.replicaSets[d.id]

	if rs.replica > replica {
		for i := 0; i < rs.replica-replica; i++ {
			s.scaleDown(d)
		}
	} else if rs.replica < replica {
		for i := 0; i < rs.replica-replica; i++ {
			s.scaleUp(d)
		}
	}
}

func (s Stack) generateNewPod(deployment *Deployment) Pod {
	server := s.getRandomServer()
	pod := Pod{
		guid.New().String(),
		&s,
		deployment.id,
		server.id,
		2,
	}
	return pod
}

func (p Pod) MigrateTo(serverId int) {
	if p.stack.cluster.placement[p.deploymentId] != -1 {
		if serverId != p.stack.cluster.placement[p.deploymentId] {
			return // do not migrate if destination doesn't match placement constraint
		}
	}

	p.serverId = serverId
}

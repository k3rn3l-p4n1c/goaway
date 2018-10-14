package scheduler

import (
	"errors"
	"net"
)

type NodeState int

const (
	Idle NodeState = iota
	Starting
	Healthy
	Unhealthy
)

type Evaluator func(m Cluster) float64

type Deployment struct {
	cluster *Cluster
	id      int
}

type DataCenter struct {
	id        int
	name      string
	serverIds []int
}

type Server struct {
	id           int
	dataCenterId int
	memoryCap    uint64
	cpu          uint
	ip           net.IP
	name         string
	state        NodeState
}

type Cluster struct {
	dataCenters []DataCenter
	deployments []Deployment
	servers     []Server
	coupling    [][]float64
	placement   []int
}

func (cluster *Cluster) findDataCenterByName(name string) (*DataCenter, error) {
	for _, dataCenter := range cluster.dataCenters {
		if dataCenter.name == name {
			return &dataCenter, nil
		}
	}
	return nil, errors.New("data center named `" + name + "` not found")
}

func (cluster *Cluster) findServerByIp(ip net.IP) (*Server, error) {
	for _, server := range cluster.servers {
		if server.ip.Equal(ip) {
			return &server, nil
		}
	}
	return nil, errors.New("server with ip <" + ip.String() + "> not found")
}

func (cluster *Cluster) AddOrUpdateServer(name, dataCenterName string, memoryCap uint64, ip net.IP) error {
	var newId int
	if len(cluster.servers) == 0 {
		newId = 0
	} else {
		newId = cluster.servers[len(cluster.servers)-1].id + 1
	}
	dataCenter, err := cluster.findDataCenterByName(dataCenterName)
	if err != nil {
		return err
	}

	server, _ := cluster.findServerByIp(ip)

	newServer := Server{
		id:           newId,
		dataCenterId: dataCenter.id,
		memoryCap:    memoryCap,
		ip:           ip,
		name:         name,
		state:        Healthy, // todo
	}
	if server == nil {
		cluster.servers = append(cluster.servers, newServer)
		dataCenter.serverIds = append(dataCenter.serverIds, newId)
	} else {
		cluster.servers[server.id] = newServer
	}

	return nil
}

func (cluster *Cluster) AddOrUpdateDataCenter(name string) error {
	var lastId int
	if len(cluster.dataCenters) == 0 {
		lastId = 0
	} else {
		lastId = cluster.dataCenters[len(cluster.dataCenters)-1].id
	}
	dataCenter, _ := cluster.findDataCenterByName(name)

	newDataCenter := DataCenter{
		id:        lastId + 1,
		name:      name,
		serverIds: []int{},
	}
	if dataCenter == nil {
		cluster.dataCenters = append(cluster.dataCenters, newDataCenter)
	} else {
		cluster.dataCenters[dataCenter.id] = newDataCenter
	}

	return nil
}

func (cluster *Cluster) AddOrUpdateDeployment(id int) {
	cluster.placement = append(cluster.placement, -1) // todo
	cluster.deployments = append(cluster.deployments, Deployment{
		cluster: cluster,
		id: id,
	})
}
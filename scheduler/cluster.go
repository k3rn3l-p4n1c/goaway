package scheduler

type Evaluator func(m Cluster) float64

type Deployment struct {
	cluster *Cluster
	id      int
}

type DataCenter struct {
	id        int
	serverIds []int
}

type Server struct {
	id           int
	dataCenterId int
	memoryCap    uint
	cpu          uint
}

type Cluster struct {
	dataCenters []DataCenter
	deployments []Deployment
	servers     []Server
	coupling    [][]float64
	placement   []int
}

func (s *Stack) getServerUsage(serverId int) (sum uint) {
	for _, pod := range s.pods {
		if pod.serverId == serverId {
			sum += pod.memoryUsage
		}
	}
	return sum
}
func (s *Stack) getRandomServer() Server {
	var minUsed Server
	var minUsage uint = 100000
	for _, server := range s.cluster.servers {
		usage := s.getServerUsage(server.id)
		if usage < minUsage {
			minUsage = usage
			minUsed = server
		}
	}
	return minUsed
}

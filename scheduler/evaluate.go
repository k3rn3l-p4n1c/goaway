package scheduler

// return value between 0 and 1. bigger is better

func rpc(s *Stack, c *Cluster) float64 {
	value := 1.0
	for _, p1 := range s.pods {
		for _, p2 := range s.pods {
			weigh := c.coupling[p1.deploymentId][p2.deploymentId]
			server1, server2 := c.servers[p1.serverId], c.servers[p2.serverId]
			switch {
			case server1.id == server2.id:
				break
			case server1.dataCenterId == server2.dataCenterId:
				value *= 1 - (1-weigh)/1000
				break
			default:
				value *= 1 - (1-weigh)/100
				break
			}
		}
	}
	return value
}

func utilization(s *Stack, c *Cluster) float64 {
	value := 1.0

	memoryUsage := make([]uint64, len(c.servers))

	for _, pod := range s.pods {
		memoryUsage[pod.serverId] += pod.memoryUsage
	}

	for i, server := range c.servers {
		if memoryUsage[i] != 0 {
			value *= float64(memoryUsage[i]) / float64(server.memoryCap)
		} else {
			value = (3 + value) / 4
		}
	}
	return value
}

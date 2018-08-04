package scheduler

// return value between 0 and 1. bigger is better

func rpc(c *Cluster) float64 {
	value := 1.0
	for elem := c.pods.Front(); elem != nil; elem = elem.Next() {
		p1 := elem.Value.(*Pod)
		for elem := c.pods.Front(); elem != nil; elem = elem.Next() {
			p2 := elem.Value.(*Pod)
			weigh := c.coupling[p1.deployment.id][p2.deployment.id]
			switch {
			case p1.server == p2.server:
				break
			case p1.server.dataCenter == p2.server.dataCenter:
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

func utilization(c *Cluster) float64 {
	value := 1.0

	memoryUsage := make([]uint, len(c.servers))

	for elem := c.pods.Front(); elem != nil; elem = elem.Next() {
		pod := elem.Value.(*Pod)
		memoryUsage[pod.server.id] += pod.memoryUsage
	}

	for i, server := range c.servers {
		if memoryUsage[i] != 0 {
			value *= float64(memoryUsage[i]) / float64(server.memoryCap)
		}
	}
	return value
}

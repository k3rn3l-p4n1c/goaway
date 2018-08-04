package scheduler


func capacity(c *Cluster) bool {
	memoryUsage := make([]uint, len(c.servers))

	for elem := c.pods.Front(); elem != nil; elem = elem.Next() {
		pod := elem.Value.(*Pod)
		memoryUsage[pod.server.id] += pod.memoryUsage
	}

	for i, server := range c.servers {
		if memoryUsage[i] > server.memoryCap {
			return false
		}
	}
	return true
}

func placement(c *Cluster) bool {
	for elem := c.pods.Front(); elem != nil; elem = elem.Next() {
		pod := elem.Value.(*Pod)
		if c.placement[pod.deployment.id] != nil {
			if pod.server.id != *c.placement[pod.deployment.id]{
				return false
			}
		}
	}
	return true
}
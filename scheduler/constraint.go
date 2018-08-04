package scheduler


func capacity(c *Cluster) bool {
	memoryUsage := make([]uint, len(c.servers))

	for elem := c.pods.Front(); elem != nil; elem = elem.Next() {
		pod := elem.Value.(*Pod)
		memoryUsage[pod.server.id] += pod.memoryUsage
	}

	for i, server := range c.servers {
		if memoryUsage[i] > server.memoryCap {
			return true
		}
	}
	return false
}

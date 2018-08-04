package scheduler


func capacity(s *Stack, c *Cluster) bool {
	memoryUsage := make([]uint, len(c.servers))

	for _, pod := range s.pods{
		memoryUsage[pod.serverId] += pod.memoryUsage

	}

	for i, server := range c.servers {
		if memoryUsage[i] > server.memoryCap {
			return false
		}
	}
	return true
}

func placement(s *Stack, c *Cluster) bool {
	for _, pod := range s.pods{
		if c.placement[pod.deploymentId] != -1 {
			if pod.serverId != c.placement[pod.deploymentId]{
				return false
			}
		}
	}
	return true
}
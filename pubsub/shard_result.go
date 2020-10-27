package pubsub

type shardResult struct {
	topicsUp, topicsDown []string
}

type shardResultCollector struct {
	m map[int]shardResult
}

func (c shardResultCollector) collect(s int, r shardResult) {
	if collected, exists := c.m[s]; exists {
		collected.topicsDown = append(collected.topicsDown, r.topicsDown...)
		collected.topicsUp = append(collected.topicsUp, r.topicsUp...)
		return
	}

	c.m[s] = r
}

func (s shardResultCollector) Results() map[int]shardResult {
	return s.m
}

func newShardResultCollector() shardResultCollector {
	return shardResultCollector{map[int]shardResult{}}
}
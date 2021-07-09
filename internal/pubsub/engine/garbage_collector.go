package engine

import "time"

type garbageCollector struct {
	invalidationTime int64
	keys             map[string]int64
}

func (g *garbageCollector) add(key string, ts int64) {
	t, ok := g.keys[key]
	if ok && t >= ts {
		return
	}
	g.keys[key] = ts
}

func (g *garbageCollector) delete(key string, ts int64) {
	t, ok := g.keys[key]
	if !ok || (ok && t >= ts) {
		return
	}
	delete(g.keys, key)
}

func (g *garbageCollector) collect(f func(string)) {
	t := time.Now().UnixNano()
	for key, keyTime := range g.keys {
		if t-keyTime >= g.invalidationTime {
			f(key)
			delete(g.keys, key)
		}
	}
}

func newGarbageCollector(invalidationTime int64) *garbageCollector {
	return &garbageCollector{
		invalidationTime: invalidationTime,
		keys:             map[string]int64{},
	}
}

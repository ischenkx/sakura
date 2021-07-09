package engine

type presenceInfo struct {
	Active bool
	TimeStamp int64
}

type presenceManager struct {
	clients map[string]presenceInfo
}

func (p *presenceManager) lastTouch(id string) int64 {
	return p.clients[id].TimeStamp
}

func (p *presenceManager) delete(id string) {
	delete(p.clients, id)
}

func (p *presenceManager) isActive(id string) bool {
	info, ok := p.clients[id]
	return ok && info.Active
}

func (p *presenceManager) update(id string, isActive bool, ts int64) {
	info, ok := p.clients[id]
	if ok && info.TimeStamp > ts {
		return
	}
	p.clients[id] = presenceInfo{
		Active:    isActive,
		TimeStamp: ts,
	}
}

func newPresenceManager() *presenceManager {
	return &presenceManager{clients: map[string]presenceInfo{}}
}

package user

type presenceInfo struct {
	TimeStamp int64
	Active bool
}

type presenceManager struct {
	users map[string]map[string]presenceInfo
}

func (p *presenceManager) isActive(id string) bool {
	infos, ok := p.users[id]
	if !ok {
		return false
	}
	for _, info := range infos {
		if info.Active {
			return true
		}
	}
	return false
}

func (p *presenceManager) delete(id, client string) {
	delete(p.users[id], client)
}

func (p *presenceManager) update(id, client string, isActive bool, ts int64) {
	u, ok := p.users[id]
	if !ok {
		u = map[string]presenceInfo{}
		p.users[id] = u
	}

	info, ok := u[client]
	if info.TimeStamp > ts && ok {
		return
	}

	info.TimeStamp = ts
	info.Active = isActive

	u[client] = info
}
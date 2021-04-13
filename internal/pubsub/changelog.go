package pubsub

type TopicAlteration struct {
	TopicID string
	ClientsAdded, UsersAdded, ClientsDeleted, UsersDeleted []string
}

func newTopicAlteration(t string) TopicAlteration {
	return TopicAlteration{
		TopicID:        t,
	}
}

type PerTopicAlterationsAggregator struct {
	Alterations []TopicAlteration
	currentTopic string
}

func (p *PerTopicAlterationsAggregator) checkLast(t string) bool {
	if len(p.Alterations) == 0 {
		return false
	}
	return p.Alterations[len(p.Alterations)-1].TopicID == t
}

func (p *PerTopicAlterationsAggregator) findOrCreate(t string) *TopicAlteration {
	if p.checkLast(t) {
		return &p.Alterations[len(p.Alterations) - 1]
	}
	for i := 0; i < len(p.Alterations); i++ {
		if p.Alterations[i].TopicID == t {
			p.Alterations[i], p.Alterations[len(p.Alterations)-1] = p.Alterations[len(p.Alterations)-1], p.Alterations[i]
			return &p.Alterations[i]
		}
	}
	p.Alterations = append(p.Alterations, newTopicAlteration(t))
	return &p.Alterations[len(p.Alterations) - 1]
}

func (p *PerTopicAlterationsAggregator) addUser(t, user string) {
	alt := p.findOrCreate(t)
	alt.UsersAdded = append(alt.UsersAdded, user)
}

func (p *PerTopicAlterationsAggregator) deleteUser(t, user string) {
	alt := p.findOrCreate(t)
	alt.UsersDeleted = append(alt.UsersDeleted, user)
}

func (p *PerTopicAlterationsAggregator) addClient(t, client string) {
	alt := p.findOrCreate(t)
	alt.ClientsAdded = append(alt.ClientsAdded, client)
}

func (p *PerTopicAlterationsAggregator) deleteClient(t, client string) {
	alt := p.findOrCreate(t)
	alt.ClientsDeleted = append(alt.ClientsDeleted, client)
}

type ChangeLog struct {
	TopicsCreated      []string
	TopicsDeleted      []string
	UsersCreated       []string
	UsersDeleted       []string
	ClientsCreated     []string
	ClientsDeleted     []string
	ClientsInactivated []string
	Topics 			   PerTopicAlterationsAggregator
	timestamp          int64
}

func (c *ChangeLog) addCreatedTopic(id string) {
	c.TopicsCreated = append(c.TopicsCreated, id)
}

func (c *ChangeLog) addDeletedTopic(id string) {
	c.TopicsDeleted = append(c.TopicsDeleted, id)
}

func (c *ChangeLog) addCreatedUser(id string) {
	c.UsersCreated = append(c.UsersCreated, id)
}

func (c *ChangeLog) addDeletedUser(id string) {
	c.UsersDeleted = append(c.UsersDeleted, id)
}

func (c *ChangeLog) addCreatedClient(id string) {
	c.ClientsCreated = append(c.ClientsCreated, id)
}

func (c *ChangeLog) addDeletedClient(id string) {
	c.ClientsDeleted = append(c.ClientsDeleted, id)
}

func (c *ChangeLog) addInactivatedClient(id string) {
	c.ClientsInactivated = append(c.ClientsInactivated, id)
}

func newChangeLog(ts int64) ChangeLog {
	return ChangeLog{
		timestamp: ts,
	}
}

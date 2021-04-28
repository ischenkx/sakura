package pubsub

import "github.com/RomanIschenko/notify/pubsub"

type topicAlteration struct {
	topicID                                                string
	clientsAdded, usersAdded, clientsDeleted, usersDeleted []string
}

func (alt *topicAlteration) TopicID() string {
	return alt.topicID
}

func (alt *topicAlteration) ClientsAdded() []string{return alt.clientsAdded}
func (alt *topicAlteration) UsersAdded() []string{return alt.usersAdded}
func (alt *topicAlteration) ClientsDeleted() []string{return alt.clientsDeleted}
func (alt *topicAlteration) UsersDeleted() []string{return alt.usersDeleted}

func newTopicAlteration(t string) *topicAlteration {
	return &topicAlteration{
		topicID: t,
	}
}

type topicAggregator struct {
	alterations  []pubsub.TopicChangeLog
	currentTopic string
}

func (p *topicAggregator) checkLast(t string) bool {
	if len(p.alterations) == 0 {
		return false
	}
	return p.alterations[len(p.alterations)-1].TopicID() == t
}

func (p *topicAggregator) findOrCreate(t string) *topicAlteration {
	if p.checkLast(t) {
		d := p.alterations[len(p.alterations) - 1].(*topicAlteration)
		return d
	}
	for i := 0; i < len(p.alterations); i++ {
		if p.alterations[i].TopicID() == t {
			p.alterations[i], p.alterations[len(p.alterations)-1] = p.alterations[len(p.alterations)-1], p.alterations[i]
			return p.alterations[i].(*topicAlteration)
		}
	}
	p.alterations = append(p.alterations, newTopicAlteration(t))
	return p.alterations[len(p.alterations) - 1].(*topicAlteration)
}

func (p *topicAggregator) addUser(t, user string) {
	alt := p.findOrCreate(t)
	alt.usersAdded = append(alt.usersAdded, user)
}

func (p *topicAggregator) deleteUser(t, user string) {
	alt := p.findOrCreate(t)
	alt.usersDeleted = append(alt.usersDeleted, user)
}

func (p *topicAggregator) addClient(t, client string) {
	alt := p.findOrCreate(t)
	alt.clientsAdded = append(alt.clientsAdded, client)
}

func (p *topicAggregator) deleteClient(t, client string) {
	alt := p.findOrCreate(t)
	alt.clientsDeleted = append(alt.clientsDeleted, client)
}

type ChangeLog struct {
	topicCreated       []string
	topicsDeleted      []string
	usersCreated       []string
	usersDeleted       []string
	clientsCreated     []string
	clientsDeleted     []string
	clientsInactivated []string
	topics             topicAggregator
	timestamp          int64
}

func (c *ChangeLog) addCreatedTopic(id string) {
	c.topicCreated = append(c.topicCreated, id)
}

func (c *ChangeLog) addDeletedTopic(id string) {
	c.topicsDeleted = append(c.topicsDeleted, id)
}

func (c *ChangeLog) addCreatedUser(id string) {
	c.usersCreated = append(c.usersCreated, id)
}

func (c *ChangeLog) addDeletedUser(id string) {
	c.usersDeleted = append(c.usersDeleted, id)
}

func (c *ChangeLog) addCreatedClient(id string) {
	c.clientsCreated = append(c.clientsCreated, id)
}

func (c *ChangeLog) addDeletedClient(id string) {
	c.clientsDeleted = append(c.clientsDeleted, id)
}

func (c *ChangeLog) addInactivatedClient(id string) {
	c.clientsInactivated = append(c.clientsInactivated, id)
}

func (c *ChangeLog) TopicsCreated() []string {return c.topicCreated}
func (c *ChangeLog) TopicsDeleted()      []string {return c.topicsDeleted}
func (c *ChangeLog) UsersCreated()       []string {return c.usersCreated}
func (c *ChangeLog) UsersDeleted()       []string {return c.usersDeleted}
func (c *ChangeLog) ClientsCreated()     []string {return c.clientsCreated}
func (c *ChangeLog) ClientsDeleted()     []string {return c.clientsDeleted}
func (c *ChangeLog) ClientsInactivated() []string {return c.clientsInactivated}
func (c *ChangeLog) Topics()             []pubsub.TopicChangeLog {return c.topics.alterations}
func (c *ChangeLog) Timestamp()          int64 {return c.timestamp}

func newChangeLog(ts int64) *ChangeLog {
	return &ChangeLog{
		timestamp: ts,
	}
}

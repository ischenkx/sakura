package pubsub

type ChangeLog struct {
	TopicsCreated []string
	TopicsDeleted []string
	UsersCreated []string
	UsersDeleted []string
	ClientsCreated []string
	ClientsDeleted []string
	ClientsInactivated []string
	timestamp int64
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
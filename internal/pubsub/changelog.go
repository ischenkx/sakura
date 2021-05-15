package pubsub

type ChangeLog struct {
	topicsCreated      []string
	topicsDeleted      []string
	usersCreated       []string
	usersDeleted       []string
	clientsCreated     []string
	clientsDeleted     []string
	clientsInactivated []string
	timestamp int64
}

func (c *ChangeLog) addCreatedTopic(id string) {
	c.topicsCreated = append(c.topicsCreated, id)
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

func (c *ChangeLog) TopicsCreated() []string      { return c.topicsCreated }
func (c *ChangeLog) TopicsDeleted() []string      { return c.topicsDeleted }
func (c *ChangeLog) UsersCreated() []string       { return c.usersCreated }
func (c *ChangeLog) UsersDeleted() []string       { return c.usersDeleted }
func (c *ChangeLog) ClientsCreated() []string     { return c.clientsCreated }
func (c *ChangeLog) ClientsDeleted() []string     { return c.clientsDeleted }
func (c *ChangeLog) ClientsInactivated() []string { return c.clientsInactivated }
func (c *ChangeLog) Timestamp() int64             { return c.timestamp }

func newChangeLog(ts int64) *ChangeLog {
	return &ChangeLog{
		timestamp: ts,
	}
}

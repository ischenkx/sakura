package changelog

type ChangeLog struct {
	ClientsUp   []string
	ClientsDown []string
	UsersUp     []string
	UsersDown   []string
	TopicsUp    []string
	TopicsDown  []string
	TimeStamp   int64
}

func (c *ChangeLog) Merge(c1 ChangeLog) {
	c.ClientsUp = append(c.ClientsUp, c1.ClientsUp...)
	c.ClientsDown = append(c.ClientsDown, c1.ClientsDown...)
	c.UsersUp = append(c.UsersUp, c1.UsersUp...)
	c.UsersDown = append(c.UsersDown, c1.UsersDown...)
	c.TopicsUp = append(c.TopicsUp, c1.TopicsUp...)
	c.TopicsDown = append(c.TopicsDown, c1.TopicsDown...)
}

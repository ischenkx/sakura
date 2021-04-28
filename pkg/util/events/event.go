package events

type Event struct {
	Name                   string
	Data                   interface{}
	Topics, Users, Clients []string
	TimeStamp              int64
	MetaInfo               interface{}
}

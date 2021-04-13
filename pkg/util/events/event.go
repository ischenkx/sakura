package events

type Event struct {
	Name                   string
	Data                   interface{}
	Topics, Users, Clients []string
	NoBuffering            bool
	Seq                    int64
	MetaInfo               interface{}
}

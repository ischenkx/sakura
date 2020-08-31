package notifier

import "github.com/RomanIschenko/notify"

type Auth interface {
	Verify(interface{}) (notify.ClientInfo, bool)
}

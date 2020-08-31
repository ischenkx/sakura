package notifier

import (
	"context"
	"github.com/RomanIschenko/notify"
	"io"
)

type Server interface {
	//first argument - data to be passed to auth or ClientInfo
	//second - transport
	Connect(interface{}, notify.Transport) (*notify.Client, error)
	DisconnectClient(client *notify.Client)
	//u pass client and reader, it sends it to data handler,
	//if it returns non-nil error then client is to be deleted.
	Handle(client *notify.Client, r io.Reader) error
	Run(ctx context.Context)
}

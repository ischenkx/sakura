package balancer

import (
	"encoding/json"
	"errors"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/cluster/broker"
	taskctx "github.com/RomanIschenko/notify/task_context"
	"time"
)

func Register(ctx *taskctx.TaskContext, b broker.Broker, app *notify.App, addrFunc func() string) error {
	if b == nil {
		return errors.New("can not connect the load balancer_pb network (broker is nil)")
	}
	b.Subscribe([]string{GlobalChannel}, time.Now().UnixNano())

	ctx.Defer(func() {
		b.Publish([]string{GlobalChannel}, broker.NewEvent(InstanceDownEvent, []byte(app.ID())))
	})

	lbLogger.Debugf("registered app(%s) in a load balancer_pb", app.ID())
	b.Handle(func(e broker.Event) {
		if e.Name == PingEvent {
			ip := addrFunc()
			id := string(e.Data)
			pong := Pong{
				ID:       id,
				IP:       ip,
			}
			data, err := json.Marshal(pong)
			if err != nil {
				return
			}
			b.Publish([]string{GlobalChannel}, broker.NewEvent(PongEvent, data))
		}
	})
	b.Publish([]string{GlobalChannel}, broker.NewEvent(InstanceUpEvent, []byte(addrFunc())))

	return nil
}
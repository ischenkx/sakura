package gateway

import (
	"fmt"
	"github.com/ischenkx/notify"
	"github.com/ischenkx/notify/examples/framework/app/services"
	"strings"
)

// notify:handler prefix=""
type Handler struct {
	uniqueService *services.UniqueIDProvider
}

type Request struct {
	Name string
}

// notify:on event=message
func (h *Handler) OnMessage(client notify.Client, data string, req Request, strs []string) {
	fmt.Println(strs)
	client.Emit("response", fmt.Sprintf("recv(\"%s\"), uid(\"%s\"), req_name(\"%s\")", data, h.uniqueService.Get(), req.Name))
}

// notify:on event=joinStrings
func (h *Handler) OnJoinStrings(client notify.Client, strs []string, sep string) {
	client.Emit("response", strings.Join(strs, sep))
}

// notify:on event=request
func (h *Handler) OnRequest(client notify.Client, tid string, payload string) {
	client.Emit("response", tid, fmt.Sprintf("from_server(%s)", payload))
}
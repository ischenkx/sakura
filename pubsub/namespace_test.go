package pubsub

import (
	"fmt"
	"github.com/google/uuid"
	_ "net/http/pprof"
	"testing"
)

func TestParseTopicData(t *testing.T) {
	fmt.Println(parseTopicData("chats:global"))
}

func TestTopic(t *testing.T) {
	registry := newNamespaceRegistry()
	registry.Register("2c", NamespaceConfig{
		MaxClients: 2,
	})

	registry.Register("2u", NamespaceConfig{
		MaxUsers: 2,
	})

	registry.Register("5c2u", NamespaceConfig{
		MaxUsers: 2,
		MaxClients: 5,
	})

	client1, _ := newClient(NewClientID("1", uuid.New().String()), 1)
	client2, _ := newClient(NewClientID("2", uuid.New().String()), 1)
	client3, _ := newClient(NewClientID("3", uuid.New().String()), 1)
	client4, _ := newClient(NewClientID("1", uuid.New().String()), 1)
	client5, _ := newClient(NewClientID("1", uuid.New().String()), 1)
	client6, _ := newClient(NewClientID("1", uuid.New().String()), 1)

	c2 := registry.generate("2c:2c")
	u2 := registry.generate("2u:2u")
	fmt.Println(u2.cfg)
	c5u2 := registry.generate("5c2u:5c2u")

	fmt.Println("c2")
	fmt.Println("(fine)", c2.add(client1))
	fmt.Println("(fine)", c2.add(client2))
	fmt.Println("(error)", c2.add(client3))

	fmt.Println("u2")

	fmt.Println("(fine)", u2.add(client1))
	fmt.Println("(fine)", u2.add(client2))
	fmt.Println("(error)", u2.add(client3))
	fmt.Println(len(u2.subs))
	fmt.Println("c5u2")

	fmt.Println("(fine)", c5u2.add(client1))
	fmt.Println("(fine)", c5u2.add(client4))
	fmt.Println("(fine)", c5u2.add(client5))
	fmt.Println("(fine)", c5u2.add(client6))
	fmt.Println("(fine)", c5u2.add(client2))
	fmt.Println("(error)", c5u2.add(client3))

}
package pubsub

import (
	"os"
	"testing"
	"time"
)

var PS = New(Config{
	InvalidationTime: time.Minute,
	CleanInterval:    time.Second,
})

func TestPubSub_Connect(t *testing.T) {
	clientId := "client1"
	userId := "user1"

	client, changelog, reconnected, err := PS.Connect(ConnectOptions{
		ClientID:  clientId,
		UserID:    userId,
		Writer:    os.Stdout,
		TimeStamp: time.Now().UnixNano(),
		Meta:      "testing the connect method",
	})

	// checking the validity of simple connection establishment
	if err != nil {
		t.Errorf("failed to connect: %v", err)
		return
	}
	if reconnected {
		t.Errorf("false reconnection")
		return
	}
	if client == nil {
		t.Errorf("returned client is nil")
		return
	}
	if len(changelog.clientsCreated) == 0 {
		t.Errorf("changelog does not contain any created clients")
		return
	}
	if len(changelog.clientsCreated) > 1 {
		t.Errorf("changelog contains too many created clients (> 1): %v", changelog.clientsCreated)
		return
	}
	if changelog.clientsCreated[0] != clientId {
		t.Errorf("returned created client does not match real created client:\n\t%s - expected\n\t%s - received", clientId, changelog.clientsCreated[0])
		return
	}
	if len(changelog.usersCreated) == 0 {
		t.Errorf("changelog does not contain any created users")
		return
	}
	if len(changelog.usersCreated) > 1 {
		t.Errorf("changelog contains too many created users (> 1): %v", changelog.usersCreated)
		return
	}
	if changelog.usersCreated[0] != userId {
		t.Errorf("returned created user does not match real created user:\n\t%s - expected\n\t%s - received", userId, changelog.usersCreated[0])
		return
	}

	// checking the validity of double connecting
	client, changelog, reconnected, err = PS.Connect(ConnectOptions{
		ClientID:  clientId,
		UserID:    userId,
		Writer:    os.Stdout,
		TimeStamp: time.Now().UnixNano(),
		Meta:      "testing the connect method",
	})

	if err != nil {
		t.Errorf("double connecting must have failed, but it didn't => that's not good")
		return
	}

	t.Log("Checked first time connection and double connection")
}

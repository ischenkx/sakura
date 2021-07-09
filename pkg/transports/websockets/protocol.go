package websockets

import (
	"errors"
	"io"
)

const (
	messageCode byte = iota + 1
	pingCode
	pongCode
	authAckCode
	authReqCode
)

type message struct {
	opCode byte
	data   []byte
}

func readMessage(data []byte) (message, error) {
	if len(data) == 0 {
		return message{}, errors.New("invalid length")
	}
	return message{data: data[1:], opCode: data[0]}, nil
}

func writeMessage(writer io.Writer, m message) (int, error) {
	n, err := writer.Write([]byte{m.opCode})

	if err != nil {
		return 0, err
	}

	n1, err := writer.Write(m.data)
	if err != nil {
		return 0, err
	}

	return n1 + n, nil
}

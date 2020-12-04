package redibroker

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/RomanIschenko/notify/cluster/broker"
	"io"
)

func encodeEvent(e broker.Event, writer io.Writer) error {
	// time + name length + id (uuid - 36 bytes)
	err := binary.Write(writer, binary.LittleEndian, e.Time)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte{byte(len(e.Name))})
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte{byte(len(e.BrokerID))})
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(e.ID))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(e.BrokerID))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(e.Name))
	if err != nil {
		return err
	}
	_, err = writer.Write(e.Data)
	if err != nil {
		return err
	}
	return nil
}

func decodeEvent(data []byte) (broker.Event, error) {
	var event broker.Event
	if len(data) < 46 {
		return event, errors.New("too short event")
	}
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &event.Time); err != nil {
		return event, err
	}
	nameLen := int(data[8])
	brokerIDLen := int(data[9])
	if 9 + brokerIDLen + 36 + nameLen > len(data) {
		return event, errors.New("failed to parse data")
	}

	event.ID = string(data[9:46])
	event.BrokerID = string(data[46:46+brokerIDLen])
	event.Name = string(data[46+brokerIDLen:46+brokerIDLen+nameLen])
	event.Data = data[46+brokerIDLen+nameLen:]
	return event, nil
}


package brokers

import (
"encoding/binary"
"encoding/json"
"errors"
"github.com/RomanIschenko/notify"
"github.com/RomanIschenko/pubsub"
"sync"
)

type Codec struct {
	Decoder func([]byte) (interface{}, error)
	Encoder func(data interface{}) ([]byte, error)
}

//RBP stands for Roman's Broker Protocol
type RBP struct {
	eventCodecs sync.Map
}
//Data     interface{}
//AppID    string
//BrokerID string
//Event    eventsps.EventType
//Time     int64

//BrokerID length - 1 byte
//AppID length - 1 byte
//Event length - 2 bytes
//Data length - 2 bytes
//BrokerID
//AppID
//Event
//Time - int64
//Data
//Minimal length = 1 + 1 + 4 + 4 + 1 + 1 + 1 + 8 + 1 = 26
//Header length = 1 + 1 + 2 + 2 + 8 = 13

func (rbp *RBP) Encode(message notify.BrokerEvent) (data []byte, err error) {
	l := len(message.Event) + len(message.BrokerID) + len(message.AppID) + 19
	var encodedData []byte = nil
	if codec, ok := rbp.eventCodecs.Load(message.Event); ok {
		encodedData, err = codec.(Codec).Encoder(message.Data)
	} else {
		encodedData, err = json.Marshal(message.Data)
	}
	if err != nil {
		return
	}

	l += len(encodedData)

	data = make([]byte, l)
	data[0] = byte(len(message.BrokerID))
	data[1] = byte(len(message.AppID))
	binary.BigEndian.PutUint16(data[2:], uint16(len(message.Event)))
	binary.BigEndian.PutUint16(data[4:], uint16(len(encodedData)))
	offset := 6
	copy(data[offset:], message.BrokerID)
	offset += len(message.BrokerID)
	copy(data[offset:], message.AppID)
	offset += len(message.AppID)
	copy(data[offset:], message.Event)
	offset += len(message.Event)
	binary.BigEndian.PutUint64(data[offset:], uint64(message.Time))
	offset += 8
	copy(data[offset:], encodedData)
	return
}

func (rbp *RBP) Decode(data []byte) (message notify.BrokerEvent, err error) {
	if len(data) < 22 {
		err = errors.New("no data to decode (byte slice is nil or empty)")
		return
	}
	offset := 0
	brokerIDLength := data[offset]
	offset += 1
	appIDLength := data[offset]
	offset += 1
	eventLength := binary.BigEndian.Uint16(data[offset:])
	offset += 2
	dataLength := binary.BigEndian.Uint16(data[offset:])
	offset += 2
	totalLength := 13 + uint16(appIDLength) + uint16(brokerIDLength) + eventLength + dataLength
	if totalLength > uint16(len(data)) {
		err = errors.New("total length is less than actual length")
		return
	}
	message.BrokerID = string(data[offset:offset+int(brokerIDLength)])
	offset += int(brokerIDLength)
	message.AppID = string(data[offset:offset+int(appIDLength)])
	offset += int(appIDLength)
	message.Event = string(data[offset:offset+int(eventLength)])
	offset += int(eventLength)
	message.Time = int64(binary.BigEndian.Uint64(data[offset:]))
	offset += 8
	messageData := data[offset:offset+int(dataLength)]
	switch message.Event {
	case notify.PublishEvent:
		opts := pubsub.PublishOptions{}
		err = json.Unmarshal(messageData, &opts)
		message.Data = opts
	case notify.SubscribeEvent:
		opts := pubsub.SubscribeOptions{}
		err = json.Unmarshal(messageData, &opts)
		message.Data = opts
	case notify.UnsubscribeEvent:
		opts := pubsub.UnsubscribeOptions{}
		err = json.Unmarshal(messageData, &opts)
		message.Data = opts
	default:
		if codec, ok := rbp.eventCodecs.Load(message.Event); ok {
			message.Data, err = codec.(Codec).Decoder(messageData)
		} else {
			var data interface{}
			err = json.Unmarshal(messageData, &data)
			message.Data = data
		}
	}
	return
}

func NewRBP() *RBP {
	return &RBP{sync.Map{}}
}

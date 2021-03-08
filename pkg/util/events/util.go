package events

import (
	"errors"
	"github.com/sirupsen/logrus"
	"reflect"
)

var logger = logrus.New().WithField("source", "emitter")

func parseIncomingData(data []byte) (string, []byte, error) {
	if len(data) < 2 {
		return "", nil, errors.New("failed to parse event #1")
	}
	length := data[0]
	if len(data) < 1+int(length) {
		return "", nil, errors.New("failed to parse event #2")
	}
	name := data[1:1+length]
	event := data[1+length:]
	return string(name), event, nil
}

func compareTypes(t1, t2 reflect.Type) bool {
	return t1.String() == t2.String()
}

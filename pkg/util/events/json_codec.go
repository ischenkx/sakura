package events

import "encoding/json"

type JSONCodec struct {}

func (JSONCodec) Marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (JSONCodec) Unmarshal(from []byte, to interface{}) error {
	return json.Unmarshal(from, to)
}

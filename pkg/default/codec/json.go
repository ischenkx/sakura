package codec

import "encoding/json"

type JSON struct{}

func (JSON) Marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (JSON) Unmarshal(from []byte, to interface{}) error {
	return json.Unmarshal(from, to)
}

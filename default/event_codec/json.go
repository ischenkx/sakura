package evcodec

import "encoding/json"

type JSON struct {}

func (JSON) Marshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}
func (JSON) Unmarshal(d []byte, i interface{}) error {
	return json.Unmarshal(d, i)
}

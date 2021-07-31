package util

import "encoding/json"

func TryJSON(i interface{}) ([]byte, bool) {
	data, err := json.Marshal(i)
	return data, err == nil
}

package evcodec

import (
	"encoding/binary"
	"encoding/json"
	"errors"
)

type JSON struct {}

func (JSON) Marshal(i []interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (JSON) Unmarshal(d []byte, args []interface{}) error {
	argsRead := 0
	for argsRead < len(args) {
		if len(d) < 2 {
			return errors.New("error while decoding #1")
		}
		dataLength := int(binary.LittleEndian.Uint16(d))
		d = d[2:]

		if len(d) < dataLength {
			return errors.New("error while decoding #2")
		}
		argData := d[:dataLength]
		d = d[dataLength:]
		if err := json.Unmarshal(argData, args[argsRead]); err != nil {
			return err
		}
		argsRead++
	}

	return nil
}

package mime

import (
	"bytes"
	"encoding/gob"
)

func BinaryEncode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BinaryDecode(data []byte) (interface{}, error) {
	return data, nil
}

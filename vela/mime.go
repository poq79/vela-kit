package vela

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type EncodeFunc func(interface{}) ([]byte, error)
type DecodeFunc func([]byte) (interface{}, error)

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

func JsonEncode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func JsonDecode(data []byte) (interface{}, error) {
	var v interface{}
	err := json.Unmarshal(data, &v)
	return v, err
}

type MimeByEnv interface {
	Mime(interface{}, EncodeFunc, DecodeFunc)
	MimeDecode(string, []byte) (interface{}, error) //func(mime_name , chunk) (object , error)
	MimeEncode(interface{}) ([]byte, string, error) //func(object) (chunk , mime_name , error)
}

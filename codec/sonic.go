package codec

import "encoding/json"

type Sonic struct{}

func (s Sonic) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (s Sonic) Unmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}

func (s Sonic) Name() string {
	return "sonic"
}

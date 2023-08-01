package codec

type Mime interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(b []byte, v interface{}) error
	// name of this codec
	Name() string
}

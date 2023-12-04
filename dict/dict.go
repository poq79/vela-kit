package dict

type Than func(string) (over bool)

type Iterator interface {
	Next() bool
	Done()
	Text() string
	Close() error
	Reset() error
}

type Dictionary interface {
	Wrap() error
	ForEach(Than) error
	Iterator() Iterator
}

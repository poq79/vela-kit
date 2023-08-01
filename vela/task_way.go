package vela

const (
	CONSOLE Way = iota + 10
	TRANSPORT
	INLINE
	AGAIN
	Scanner
)

type Way uint8

func (way Way) String() string {
	switch way {
	case CONSOLE:
		return "console"
	case TRANSPORT:
		return "tunnel"
	case INLINE:
		return "inline"
	case AGAIN:
		return "again"
	default:
		return "unknown"
	}
}

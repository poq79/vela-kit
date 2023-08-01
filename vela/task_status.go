package vela

const (
	Running TaskStatus = iota + 1
	Doing
	Panic
	Fail
	Register
	Update
	Updating
)

type TaskStatus uint8

func (t TaskStatus) String() string {
	switch t {
	case Running:
		return "running"
	case Doing:
		return "doing"
	case Fail:
		return "fail"

	case Panic:
		return "panic"
	case Register:
		return "reg"

	case Update:
		return "update"

	default:
		return ""
	}
}

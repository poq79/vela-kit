package audit

import (
	"github.com/vela-ssoc/vela-kit/grep"
)

type match func(*Event) bool

func newFilter(key string, pattern string) match {
	filter := grep.New(pattern)
	return func(ev *Event) bool {
		val := ev.Field(key)
		//xEnv.Debugf("%s grep %s", val, pattern)
		return filter(val)
	}
}

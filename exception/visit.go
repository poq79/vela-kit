package exception

import (
	"github.com/vela-ssoc/vela-kit/stdutil"
)

func Fatal(e error) {
	if e == nil {
		return
	}
	out := stdutil.New(stdutil.Daemon())
	defer func() { _ = out.Close() }()
	out.ERR("fatal error %s", e.Error())
}

func Try(e error, protect bool) {
	if e == nil {
		return
	}

	if protect {
		defer func() {
			if cause := recover(); cause == nil {
				return
			} else {
				out := stdutil.New(stdutil.Daemon())
				defer func() { _ = out.Close() }()
				out.ERR("recover error %v , stack %s", cause, StackTrace(0))
			}
		}()
	}

	panic(e)
}

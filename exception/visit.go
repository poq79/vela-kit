package exception

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
)

func Fatal(e error) {
	if e == nil {
		return
	}
	show, file := auxlib.Output()
	if file != nil {
		defer func() { _ = file.Close() }()
	}

	show("fatal error %s", e.Error())
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
				show, file := auxlib.Output()
				show("recover error %v , stack %s", cause, StackTrace(0))
				if file != nil {
					file.Close()
				}
			}
		}()
	}

	panic(e)
}

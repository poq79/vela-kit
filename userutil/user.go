package userutil

import (
	"github.com/vela-ssoc/vela-kit/strutil"
	"os/user"
)

const (
	expiration = 3 * 3600
)

func Username(uid string) (uname string, err error) {
	if u, e := uc.Get(strutil.S2B(uid)); e == nil {
		uname = strutil.B2S(u)
		return
	}

	var u *user.User
	u, err = user.LookupId(uid)
	if err != nil {
		return
	}

	uname = u.Username
	uc.Set(strutil.S2B(uid), strutil.S2B(uname), expiration)
	return
}

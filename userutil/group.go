package userutil

import (
	"github.com/vela-ssoc/vela-kit/strutil"
	"os/user"
)

func Group(gid string) (name string, err error) {
	if v, ok := gc.Get([]byte(gid)); ok == nil {
		name = strutil.B2S(v)
		return
	} else {
		var g *user.Group
		g, err = user.LookupGroupId(gid)
		if err != nil {
			return
		}
		name = g.Name
		gc.Set([]byte(gid), []byte(name), expiration)
	}

	return
}

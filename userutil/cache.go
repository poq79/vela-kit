package userutil

import (
	"github.com/coocood/freecache"
)

var ( // username cache
	uc = freecache.NewCache(256)
	gc = freecache.NewCache(256)
)

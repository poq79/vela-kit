package node

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"os"
	"sync"
)

type node struct {
	once      sync.Once
	upgrading uint32
	id        string
	prefix    string
}

func newNode() *node {
	return &node{
		upgrading: 0,
		id:        "",
		prefix:    "share",
	}
}

func newLuaNode(L *lua.LState) int {

	if !L.CheckCodeVM("startup") {
		L.RaiseError("rock.node.* not allow with %s", L.CodeVM())
		return 0
	}

	tab := L.CheckTable(1)
	tab.Range(func(key string, val lua.LValue) {
		switch key {
		case "resolve":
			resolve = val.String()
		case "id":
			_G.id = val.String()

		case "prefix":
			_G.prefix = val.String()

		default:
			L.RaiseError("node not found %s", key)
		}
	})

	if e := _G.valid(); e != nil {
		xEnv.Errorf("node startup failure error %v", e)
		os.Exit(-1)
		return 0
	}
	return 0
}

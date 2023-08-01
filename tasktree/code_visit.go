package tasktree

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/radix"
	"github.com/vela-ssoc/vela-kit/vela"
	"reflect"
	"strings"
	"time"
)

func (cd *Code) AssertCodeLuaStateTagFunc() {}

func (cd *Code) Update(co *lua.LState, chunk []byte, env vela.Environment, way vela.Way) {
	cd.cancel()
	ctx, cancel := context.WithCancel(env.Context())
	co.SetContext(ctx)
	cd.cancel = cancel

	cd.header.hash = fmt.Sprintf("%x", md5.Sum(chunk))
	cd.header.way = way
	cd.header.env = env
	cd.header.err = nil
	cd.chunk = chunk
	return
}

func (cd *Code) IsReg() bool {
	t := cd.T()
	return t == vela.Register || t == vela.Update
}

func (cd *Code) Disable() bool {
	return cd.header.disable
}

func (cd *Code) Key() string {
	return cd.header.key
}

func (cd *Code) Hash() string {
	return cd.header.hash
}

func (cd *Code) From() string {
	return cd.header.way.String()
}

func (cd *Code) Time() string {
	return cd.header.uptime.Format(time.RFC3339Nano)
}

func (cd *Code) Uptime() time.Time {
	return cd.header.uptime
}

func (cd *Code) Status() string {
	return cd.T().String()
}

func (cd *Code) Link() string {
	if len(cd.header.link) == 0 {
		return ""
	}
	return strings.Join(cd.header.link, ",")
}

func (cd *Code) Exist(proc string) bool {
	return cd.vela(proc) != nil
}

func (cd *Code) Wrap() error {
	return cd.header.err
}

func (cd *Code) CompareVM(L *lua.LState) bool {
	cname := L.CodeVM()
	if cname == "" {
		return false
	}

	return cd.Key() == cname
}

func (cd *Code) NewVelaData(L *lua.LState, name string, typeof string) *lua.VelaData {
	if cd.Disable() {
		L.RaiseError("disable create new vela")
		return nil
	}

	vla := cd.vela(name)
	if vla != nil {
		if reflect.TypeOf(vla.Data).String() != typeof {
			L.RaiseError("invalid type , must %T , got %s", vla.Data, typeof)
			return nil
		}
		cd.addMask(name)
	} else {
		vla = lua.NewVelaData(nil) //有可能是nil 注意的判断
		cd.tree.Insert(radix.Key(name), vla)
		cd.addMask(name)
	}

	return vla
}

func (cd *Code) List() []*vela.Runner {
	idx := 0
	ars := make([]*vela.Runner, cd.tree.Size())
	cd.tree.ForEach(func(node radix.Node) (cont bool) {
		pd := node.Value().(*lua.VelaData)
		ars[idx] = &vela.Runner{
			Name:    string(node.Key()),
			CodeVM:  pd.CodeVM(),
			Type:    pd.Data.Type(),
			Status:  pd.Data.State().String(),
			Private: pd.IsPrivate(),
		}
		idx++
		return true
	})

	return ars
}

func (cd *Code) Get(name string) *lua.VelaData {
	return cd.vela(name)
}

func (cd *Code) newEvent(sub string) *audit.Event {
	return audit.NewEvent("TaskTree.code").Subject(sub).From(cd.Key())
}

package pipe

import (
	"fmt"
	auxlib "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	lvalue "github.com/vela-ssoc/vela-kit/reflectx"
	"github.com/vela-ssoc/vela-kit/vela"
	"io"
)

type Call interface {
	PCall(v ...interface{}) error
}

func (px *Chains) Len() int {
	return len(px.chain)
}

func (px *Chains) LValue(lv lua.LValue) {
	switch lv.Type() {

	case lua.LTUserData:
		px.Object(lv.Peek().(*lua.LUserData).Value)

	case lua.LTVelaData:
		px.Object(lv.Peek().(*lua.VelaData).Data)

	case lua.LTAnyData:
		px.Object(lv.Peek().(*lua.AnyData).Data)

	case lua.LTObject:
		px.Object(lv.Peek())

	case lua.LTGoFuncErr:
		fn := px.LFuncErr(lv.Peek().(lua.GoFuncErr))
		px.append(fn)

	case lua.LTGoFuncStr:
		fn := px.LFuncStr(lv.Peek().(lua.GoFuncStr))
		px.append(fn)

	case lua.LTGoFuncInt:
		fn := px.LFuncInt(lv.Peek().(lua.GoFuncInt))
		px.append(fn)
	case lua.LTGoFunction:
		fn := px.GoFunc(lv.Peek().(lua.GoFunction))
		px.append(fn)

	case lua.LTFunction:
		px.append(px.LFunc(lv.Peek().(*lua.LFunction)))
	default:
		px.invalid("invalid pipe lua type , got %s", lv.Type().String())
	}
}

func (px *Chains) Object(v interface{}) {
	fn := px.Preprocess(v)
	if fn == nil {
		return
	}

	px.append(fn)
}

func (px *Chains) Preprocess(v interface{}) Fn {
	switch item := v.(type) {

	case io.Writer:
		return px.Writer(item)

	case *lua.LFunction:
		return px.LFunc(item)
	case lua.Console:
		return px.Console(item)
	case Call:
		return item.PCall

	case func():
		return func(...interface{}) error {
			item()
			return nil
		}

	case func(interface{}):
		return func(...interface{}) error {
			item(v)
			return nil
		}

	case func() error:
		return func(...interface{}) error {
			item()
			return nil
		}

	case func(interface{}) error:
		return func(v ...interface{}) error {
			if len(v) == 0 {
				return nil
			}
			return item(v[0])
		}

	default:
		px.invalid("invalid pipe object")
	}

	return nil
}

func (px *Chains) GoFunc(fn lua.GoFunction) Fn {
	return func(v ...interface{}) error {
		return fn()
	}
}

func (px *Chains) LFuncErr(fn lua.GoFuncErr) Fn {
	return func(v ...interface{}) error {
		return fn(v...)
	}
}

func (px *Chains) LFuncStr(fn lua.GoFuncStr) Fn {
	return func(v ...interface{}) error {
		fn(v...)
		return nil
	}
}

func (px *Chains) LFuncInt(fn lua.GoFuncInt) Fn {
	return func(v ...interface{}) error {
		fn(v...)
		return nil
	}
}

func (px *Chains) LFunc(fn *lua.LFunction) Fn {
	return func(v ...interface{}) error {
		size := len(v)
		if size == 0 {
			return nil
		}

		var co *lua.LState
		L, ok := v[size-1].(*lua.LState)
		if ok {
			co = px.clone(L)
			v = v[:size-1]
		}
		cp := px.xEnv.P(fn)

		var args []lua.LValue
		for _, item := range v {
			args = append(args, lvalue.ToLValue(item, co))
		}
		defer px.xEnv.Free(co)

		if len(args) == 0 {
			return fmt.Errorf("xreflect to LValue fail %v", v)
		}

		return co.CallByParam(cp, args...)
	}
}

func (px *Chains) write(w io.Writer, v ...interface{}) error {
	size := len(v)
	if size == 0 {
		return nil
	}

	data, err := auxlib.ToStringE(v[0])
	if err != nil {
		return err
	}
	_, err = w.Write(auxlib.S2B(data))
	return err
}

func (px *Chains) Writer(w io.Writer) Fn {
	return func(v ...interface{}) error {
		if w == nil {
			return fmt.Errorf("invalid io writer %p", w)
		}

		return px.write(w, v...)
	}
}

func (px *Chains) SetEnv(env vela.Environment) *Chains {
	if env != nil {
		px.xEnv = env
	}
	return px
}

func (px *Chains) Console(out lua.Console) Fn {
	return func(v ...interface{}) error {
		size := len(v)
		if size == 0 {
			return nil
		}

		data, err := auxlib.ToStringE(v[0])
		if err != nil {
			return err
		}
		out.Println(data)
		return nil
	}
}

//兼容老的数据

func (px *Chains) Do(arg interface{}, co *lua.LState, x func(error)) {
	n := len(px.chain)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		fn := px.chain[i]
		if e := fn(arg, co); e != nil && x != nil {
			x(e)
		}
	}
}

func (px *Chains) Case(v interface{}, id int, cnd string, co *lua.LState) error {
	n := len(px.chain)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		fn := px.chain[i]
		if e := fn(v, id, cnd, co); e != nil {
			return e
		}
	}

	return nil
}

func (px *Chains) Call2(v1 interface{}, v2 interface{}, co *lua.LState) error {
	n := len(px.chain)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		fn := px.chain[i]
		if e := fn(v1, v2, co); e != nil {
			return e
		}
	}

	return nil
}

func (px *Chains) Call(co *lua.LState, v ...interface{}) error {
	n := len(px.chain)
	if n == 0 {
		return nil
	}

	param := append(v, co)
	for i := 0; i < n; i++ {
		fn := px.chain[i]
		if e := fn(param...); e != nil {
			return e
		}
	}

	return nil
}

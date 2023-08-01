package audit

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
)

type config struct {
	name string
	bkt  []string
	rate []inhibitMatch
	pass []match
	file string
	pipe *pipe.Chains
	sdk  lua.Writer
	co   *lua.LState
}

func velaMinConfig() *config {
	return &config{
		name: "vela.audit",
		file: "vela.audit.log",
		pipe: pipe.New(),
		bkt:  []string{"audit_inhibit_record"},
		rate: []inhibitMatch{newInhibitMatch("$inet_$id_$typeof_$from", 5*60)},
	}
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := velaMinConfig()
	cfg.co = xEnv.Clone(L)

	tab.Range(func(key string, val lua.LValue) {
		switch key {

		case "file":
			cfg.file = val.String()

		case "to":
			cfg.sdk = auxlib.CheckWriter(val, L)

		default:
			L.RaiseError("not found %s", key)
		}

	})

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	return cfg
}

func (cfg *config) verify() error {
	return nil
}

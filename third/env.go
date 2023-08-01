package third

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"path/filepath"
)

var xEnv vela.Environment

type VelaThird interface {
	Load(name string) (*vela.ThirdInfo, error)
	Info(name string) *vela.ThirdInfo
}

func Constructor(env vela.Environment, callback func(VelaThird) error) {
	xEnv = env

	//生成对象
	t := &third{
		dir:    filepath.Join(env.ExecDir(), "3rd"),
		cache:  make(map[string]*vela.ThirdInfo, 32),
		bucket: env.Bucket("VELA_THIRD_INFO_DB"),
	}
	t.check()
	t.define() //定义路由
	t.online() //上线后操作

	if err := callback(t); err != nil {
		env.Errorf("third callback exec fail %v", err)
	} //保存全局

	env.Set("third",
		lua.NewExport("lua.third.export", lua.WithIndex(t.indexL), lua.WithFunc(t.loadL)))
}

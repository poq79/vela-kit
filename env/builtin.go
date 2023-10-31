package env

import (
	"context"
	"github.com/vela-ssoc/vela-common-mba/definition"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/denoise"
	"github.com/vela-ssoc/vela-kit/env/sys"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"github.com/vela-ssoc/vela-kit/proxy"
	"github.com/vela-ssoc/vela-kit/rtable"
	"github.com/vela-ssoc/vela-kit/tasktree"
	"github.com/vela-ssoc/vela-kit/third"
	"github.com/vela-ssoc/vela-kit/variable"
	"github.com/vela-ssoc/vela-kit/vela"
	tunnel "github.com/vela-ssoc/vela-tunnel"
	"sync"
)

type substance struct {
	region vela.Region
	task   *tasktree.TaskTree
}

type onConnectEv struct {
	name string
	todo func() error
}

type monitor interface {
	Quiet() bool
	CPU() float64
	AgentCPU() float64
}

type Environment struct {
	ctx  context.Context
	stop context.CancelFunc
	tab  *EnvL    //lua vm environment
	rou  *routine //routine pool cache
	//bdb         *bboltDB               //bbolt database cache
	db          database               //database
	log         vela.Log               //External Log interface
	mbc         []vela.Closer          //Must be closed
	sub         *substance             //substance object cache
	tnl         tunnel.Tunneler        //版本2 tunnel
	adt         *audit.Audit           //审计模块
	vhu         *variable.Hub          //变量状态
	rtm         monitor                //runtime
	shm         shared                 //共享内存
	link        interface{}            //本地调用链接
	mime        *MimeHub               //mime hub object
	third       third.VelaThird        //third 三方存储
	tupMutex    sync.Mutex             //并发锁
	tuple       map[string]interface{} //存储一些关键信息
	onConnect   []onConnectEv
	router      *rtable.TnlRouter //存储注册路由信息
	hide        definition.MinionHide
	onReconnect *pipe.Chains
}

func with(env *Environment) {

	//注入系统变量
	sys.WithEnv(env)
	proxy.WithEnv(env)
	denoise.With(env)

	//注入函数
	env.Set("go", lua.NewFunction(env.thread))

	//注入信号量
	env.Set("notify", lua.NewFunction(env.notifyL))

	//设置节点信息
	env.Set("ID", lua.NewFunction(env.nodeIDL))
	env.Set("inet", lua.NewFunction(env.inetL))
	env.Set("kernel", lua.S2L(env.Kernel()))
	env.Set("inet6", lua.NewFunction(env.inet6L))
	env.Set("mac", lua.NewFunction(env.macL))
	env.Set("addr", lua.NewFunction(env.addrL))
	env.Set("arch", lua.NewFunction(env.archL))
	env.Set("prefix", lua.S2L(env.ExecDir()))

	env.Set("tunnel",
		lua.NewExport("lua.tunnel.export",
			lua.WithIndex(env.tunnelIndexL),
		))

	env.Set("exdata",
		lua.NewExport("vela.exdata.export",
			lua.WithIndex(env.exdataIndexL),
			lua.WithNewIndex(env.setExdataL)))

	env.Set("hide",
		lua.NewExport("vela.hide.export",
			lua.WithIndex(env.hideIndexL)))

	env.Set("broker",
		lua.NewExport("vela.broker.export",
			lua.WithIndex(env.brokerIndexL)))

}

func Create(mode string, name string, protect bool) *Environment {
	ctx, cancel := context.WithCancel(context.TODO())

	env := &Environment{
		sub:   &substance{},
		tuple: make(map[string]interface{}, 32),
		ctx:   ctx,
		stop:  cancel,
	}

	env.newEnvL(mode, name, protect)
	env.newRoutine()
	env.newMimeHub()
	env.invoke()

	with(env)
	return env
}

package agent

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/env"
	"github.com/vela-ssoc/vela-kit/vela"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Usage() {
	fmt.Printf("ssc worker -d            说明: 后台运行\n")
	fmt.Printf("ssc worker               说明: 直接运行\n")
	fmt.Printf("ssc load script.lua      说明: 加载指定配置\n")
	fmt.Printf("ssc install              说明: 安装服务\n")
	fmt.Printf("ssc uninstall            说明: 卸载服务\n")
	return
}

func Load(fn construct) {
	if len(os.Args) < 3 {
		fmt.Printf("not found script")
		return
	}

	xEnv := fn("load")
	err := xEnv.DoTaskFile(os.Args[2], vela.CONSOLE)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	xEnv.Infof("vela success by %s\n", os.Args[2])
	//signal.Notify()
}
func Start(fn construct) {

	//是否开启后台服务
	if len(os.Args) >= 3 && os.Args[2] == "-d" {
		daemon()
		return
	}

	//工作进程
	Worker(fn)
	return
}

func Worker(fn construct) {
	var cmd exec.Cmd
	var shutdown bool
	run(fn, &cmd, &shutdown)

	chn := make(chan os.Signal, 1)
	signal.Notify(chn, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-chn

	shutdown = true
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
}

func constructor(name string, startup func(env vela.Environment)) func(string) *env.Environment {
	return func(mode string) *env.Environment {
		xEnv := env.Create(mode, name, true)
		startup(xEnv)
		return xEnv
	}
}

func By(name string, use func(env vela.Environment)) {

	if len(os.Args) < 2 {
		Usage()
		return
	}

	mode := os.Args[1]
	deploy := constructor(name, use)

	switch mode {
	//case "load":
	//	Load(deploy)
	case "start":
		Start(deploy)
	case "worker":
		Worker(deploy)
	case "install":
		Install(deploy)
	case "uninstall":
		Uninstall(deploy)
	case "service":
		Service(deploy)
	case "upgrade":
		Upgrade()
	default:
		Usage()
		return
	}

}

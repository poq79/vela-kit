package agent

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/env"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	velaExecIdx = "SSC_DAEMON_IDX"
)

type construct func(string) *env.Environment

func SetEnv(cmd *exec.Cmd) {
	idx := os.Getenv(velaExecIdx)
	if idx == "fork" {
		return
	}
	cmd.Env = append(os.Environ(), velaExecIdx+"=fork")
}

func forkExec(fc *exec.Cmd) {
	output, file := auxlib.Output()
	if file != nil {
		defer func() { _ = file.Close() }()
	}

	//启动新的命令
	cmd := exec.Cmd{
		SysProcAttr: NewSysProcAttr(),
	}

	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
			cmd.Process = nil
		}
	}()

	//获取文件
	exe := executable(output)
	if exe == "" {
		return
	}

	//配置环境
	cmd.Path = exe
	cmd.Dir = filepath.Dir(exe)
	cmd.Args = []string{exe, "worker"}

	//控制环境变量
	SetEnv(&cmd)
	//cmd.Env = append(os.Environ(), velaExecIdx+"=fork")

	//控制文件输出
	if file != nil {
		cmd.Stderr = file
		cmd.Stdout = file
	}

	//开始服务
	if e := cmd.Start(); e != nil {
		output(`"msg":"fork ssc error %s" , "exe":"%s" , "dir":"%s"`, e.Error(), cmd.Path, cmd.Dir)
		return
	}

	output(`"msg":"fork ssc pid %d" , "exe":"%s" , "dir":"%s"`, cmd.Process.Pid, cmd.Path, cmd.Dir)
	*fc = cmd

	//等待命令
	if e := cmd.Wait(); e != nil {
		output(`"msg":"fork ssc pid %d exit error %s" , "exe":"%s" , "dir":"%s"`,
			e.Error(), cmd.Process.Pid, cmd.Path, cmd.Dir)
		return
	}

	output(`"msg":"fork ssc pid %d exit" , "exe":"%s" , "dir":"%s"`, cmd.Process.Pid, cmd.Path, cmd.Dir)

	return
}

// run 启动长连接服务器
func run(fn construct, cmd *exec.Cmd, shutdown *bool) {
	idx := os.Getenv(velaExecIdx)
	if idx == "fork" {
		xEnv := fn("worker") //工作模式:worker 环境变量:rock
		xEnv.Spawn(0, xEnv.Worker)
		xEnv.Notify()
		return
	}

	go func(c *exec.Cmd) {
		for {
			if *shutdown {
				return
			}
			forkExec(c)
			<-time.After(time.Second * 30)
		}
	}(cmd)
}

func daemon() {
	exe := Exe()
	cmd := exec.Cmd{
		Path:        exe,
		Dir:         filepath.Dir(exe),
		SysProcAttr: NewSysProcAttr(),
		Args:        []string{exe, "worker"},
	}

	err := cmd.Start()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("start ssc daemon %d\n", cmd.Process.Pid)
	os.Exit(0)
}

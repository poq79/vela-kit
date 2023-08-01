package node

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	ps "github.com/shirou/gopsutil/process"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"os"
)

var shell = ""

func fileMd5(filename string) (string, error) {
	pFile, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("open %v error %v", filename, err)
	}
	defer pFile.Close()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	hub := md5.New()
	auxlib.Copy(ctx, hub, pFile)
	return hex.EncodeToString(hub.Sum(nil)), nil
}

func NotUpgrade(exe string) bool {
	pro, err := ps.NewProcess(int32(os.Getppid()))
	if err != nil {
		xEnv.Errorf("not found parent process %v", err)
		return true
	}

	pex, err := pro.Exe()
	if err != nil {
		xEnv.Errorf("not found parent process executable %v", err)
		return true
	}

	h1, err := fileMd5(pex)
	if err != nil {
		xEnv.Errorf("not found parent process executable hash %v", err)
		return true
	}

	h2, err := fileMd5(exe)
	if err != nil {
		xEnv.Errorf("not found process executable hash %v", err)
		return true
	}

	if h1 == h2 {
		return true
	}

	return false
}

func (nd *node) hot(save, abs string, out func(string, ...interface{})) error {
	if err := nd.daemon(save); err != nil {
		out("升级主进程服务失败exe:%s 原因:%v", abs, err)
	}

	os.Exit(0)
	return nil
}

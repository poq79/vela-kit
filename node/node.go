package node

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-common-mba/netutil"
	"github.com/vela-ssoc/vela-kit/fileutil"
	"github.com/vela-ssoc/vela-kit/stdutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

func (nd *node) valid() error {
	if err := socket(&resolve); err != nil {
		return err
	}

	return nil
}

func (nd *node) Path() (exe string, work string, backup string) {
	out := stdutil.New(stdutil.Console())
	defer out.Close()
	var err error

	exe, err = xEnv.Exe()
	if err != nil {
		out.ERR("executable got fail %v", err)
		return
	}

	work, _ = filepath.Split(exe)
	backup = filepath.Join(work, "backup")
	return
}

func (nd *node) Name() string {
	return nd.id
}

func (nd *node) startup(ctx *fasthttp.RequestCtx) error {
	body := ctx.Request.Body()

	var kv map[string]string
	err := json.Unmarshal(body, &kv)
	if err != nil {
		return err
	}

	nd.prefix = kv["prefix"]
	nd.id = kv["id"]
	resolve = kv["resolve"]

	return nil
}

func (nd *node) query(u upgrade) string {
	return fmt.Sprintf("/api/v1/broker/upgrade/download?%s", u.Query().Encode())
}

func (nd *node) download(u string, save string) (string, error) {
	attr, err := xEnv.Attachment(u)
	if err == nil {
		hash, e := attr.File(save)
		attr.Close()
		return hash, e
	}

	tk := time.NewTicker(30 * time.Second)
	defer tk.Stop()

	ignore := func(err error) bool {
		if err == nil {
			return false
		}
		if e, ok := err.(*netutil.HTTPError); ok && e.Code == http.StatusTooManyRequests {
			return true
		}

		xEnv.Errorf("download agent binary fail %v", err)
		return false
	}

	failure := 0
	for range tk.C {
		failure++
		if failure > 360 {
			return "", fmt.Errorf("升级包下载超过360次失败")
		}

		attr, err = xEnv.Attachment(u)
		if ignore(err) {
			continue
		}

		if err == nil {
			hash, e := attr.File(save)
			attr.Close()
			return hash, e
		}
	}

	return "", fmt.Errorf("下载升级包失败")
}

func (nd *node) daemon(exe string) error {
	if NotUpgrade(exe) {
		return fmt.Errorf("daemon not upgrade %s", exe)
	}

	cmd := exec.Cmd{
		Path: exe,
		Dir:  filepath.Dir(exe),
		Args: []string{exe, "upgrade"},
	}

	return cmd.Start()
}

func (nd *node) upgrade(ctx *fasthttp.RequestCtx) error {
	out := stdutil.New(stdutil.Console())
	defer out.Close()

	if atomic.AddUint32(&nd.upgrading, 1) > 1 {
		out.Info("多次指令接收,正在升级......")
		return nil
	}

	var up upgrade
	body := ctx.Request.Body()

	err := json.Unmarshal(body, &up)
	if err != nil {
		out.Info("upgrade unmarshal fail %v", err)
		return err
	}

	out.Info("upgrade ssc %#v", up)

	err = xEnv.Spawn(0, func() {
		out2 := stdutil.New(stdutil.Console())
		defer out2.Close()

		defer atomic.StoreUint32(&nd.upgrading, 0)

		abs, er := xEnv.Exe()
		if er != nil {
			out2.ERR("executable got fail %v", err)
			return
		}

		// 获取当前的工作目录
		workdir, name := filepath.Split(abs)
		ext := filepath.Ext(name)
		if len(ext) > 0 {
			name = strings.SplitN(name, ext, 2)[0]
		}

		backDir := filepath.Join(workdir, "backup")
		backName := filepath.Join(backDir, fmt.Sprintf("%s-%s%s", name, xEnv.Hide().Edition, ext))

		// 只备份本次的二进制包, 历史备份二进制包不留存, 简单粗暴: 删除备份目录, 将本次二进制放到备份目录
		er = fileutil.CreateIfNotExists(backDir, true)
		if er != nil {
			out2.ERR("失败备份当前二进制文件: %s ---> %s fail %v", abs, backName, er)
			return
		}

		out2.Info("开始备份当前二进制文件: %s ---> %s", abs, backName)
		n, er := fileutil.CopyFile(abs, backName)
		if er != nil {
			out2.ERR("开始备份当前二进制文件: %s ---> %s %v", abs, backName, er)
			return
		}
		out2.Info("备份当前二进制成功: %s ---> %s size: %d", abs, backName, n)
		// 下载最新版本

		save := filepath.Join(workdir, fmt.Sprintf("ssc-%d%s", time.Now().Unix(), ext))

		_, er = nd.download(nd.query(up), save)
		if er != nil {
			out2.ERR("下载文件失败 %v %v", up, er)
			return
		}
		out2.Info("下载文件结束 %v", up)

		info, er := os.Stat(save)
		if er != nil {
			out2.ERR("可执行程序未成功保存%v", er)
			return
		}

		if info.Size() < 4096 {
			out2.ERR("可执行程序size < 4096")
			return
		}

		er = nd.hot(save, abs, out)
		if er != nil {
			out2.ERR("可执行程序执行失败%v", er)
		}

	})
	out.ERR("upgrade spawn start fail %v", err)
	return nil
}

func (nd *node) command(ctx *fasthttp.RequestCtx) error {

	body := ctx.Request.Body()
	if len(body) == 0 {
		return fmt.Errorf("command body empty")
	}

	var cmd command
	err := json.Unmarshal(body, &cmd)
	if err != nil {
		return err
	}

	switch cmd.Cmd {
	case "offline":
		os.Exit(0)
	}

	return nil
}

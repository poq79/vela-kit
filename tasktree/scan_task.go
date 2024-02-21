package tasktree

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/radix"
	"github.com/vela-ssoc/vela-kit/stdutil"
	"github.com/vela-ssoc/vela-kit/vela"
	"time"
)

/*

local sbom = vela.sbom{

}
local track = vela.track("java")

track.pipe(function(section)
	sbom.file(section.name)
end)
*/

type ScanTask struct {
	ctx  context.Context
	kill context.CancelFunc
	code *Code
	co   *lua.LState
}

func (s *ScanTask) report() {

}

func (s *ScanTask) StopScanTask() {
	if s.kill != nil {
		return
	}

	s.kill()
}

func (s *ScanTask) call(env vela.Environment) error {
	if s.code == nil {
		return fmt.Errorf("not found code")
	}

	fn, err := s.co.Load(bytes.NewReader(s.code.chunk), s.code.Key())
	if err != nil {
		return s.code.Wrap()
	}

	go func() {
		s.co.CallByParam(env.P(fn))
		s.code.free()
		s.report()
	}()

	return nil
}

func newScanTask(env vela.Environment, id int64, key string, chunk []byte,
	metadata map[string]interface{}, timeout int) *ScanTask {
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout <= 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	}

	co := env.Coroutine()
	co.SetContext(ctx)
	code := &Code{
		header: &header{
			id:      id,
			hash:    fmt.Sprintf("%x", md5.Sum(chunk)),
			status:  uint32(vela.Doing),
			way:     vela.Scanner,
			key:     key,
			env:     env,
			dialect: false,
		},
		tree:     radix.NewSafeTree(),
		cancel:   cancel,
		chunk:    chunk,
		metadata: metadata,
	}

	co.Exdata = code

	console := stdutil.New(stdutil.Console())
	defer console.Close()
	console.Debug("新增扫描服务 %s 通过 %s 成功", key, env.Name())

	return &ScanTask{
		co:   co,
		ctx:  ctx,
		code: code,
		kill: cancel,
	}
}
